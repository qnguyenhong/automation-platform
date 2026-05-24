package executors

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/qnguyenhong/automation-platform/internal/model"
	"github.com/qnguyenhong/automation-platform/internal/worker"
	"github.com/qnguyenhong/automation-platform/pkg/config"
	"github.com/qnguyenhong/automation-platform/pkg/variable"
)

type LoadExecutor struct {
	registry *worker.Registry
	cfg      *config.WorkerConfig
}

func NewLoadExecutor(registry *worker.Registry, cfg *config.WorkerConfig) *LoadExecutor {
	return &LoadExecutor{
		registry: registry,
		cfg:      cfg,
	}
}

func (e *LoadExecutor) Name() string {
	return "load"
}

func (e *LoadExecutor) CanHandle(testType string) bool {
	return testType == "load"
}

type requestResult struct {
	timestamp time.Time
	status    int
	duration  time.Duration
	bytes     int64
	err       bool
}

type Limiter struct {
	rate       float64
	capacity   float64
	tokens     float64
	lastRefill time.Time
	mu         sync.Mutex
}

func NewLimiter(rps float64) *Limiter {
	return &Limiter{
		rate:       rps,
		capacity:   rps,
		tokens:     rps,
		lastRefill: time.Now(),
	}
}

func (l *Limiter) Limit(ctx context.Context) error {
	for {
		l.mu.Lock()
		now := time.Now()
		elapsed := now.Sub(l.lastRefill).Seconds()
		l.lastRefill = now
		l.tokens += elapsed * l.rate
		if l.tokens > l.capacity {
			l.tokens = l.capacity
		}

		if l.tokens >= 1.0 {
			l.tokens -= 1.0
			l.mu.Unlock()
			return nil
		}
		l.mu.Unlock()

		sleepTime := time.Duration((1.0 - l.tokens) / l.rate * float64(time.Second))
		if sleepTime < 1*time.Millisecond {
			sleepTime = 1 * time.Millisecond
		}
		if sleepTime > 100*time.Millisecond {
			sleepTime = 100 * time.Millisecond
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(sleepTime):
		}
	}
}

func (e *LoadExecutor) Execute(ctx context.Context, job worker.Job) (*worker.Result, error) {
	// Resolve variables in config
	engine := variable.NewEngine(job.CapturedVars, job.EnvVars)
	resolvedConfig, err := engine.ResolveMap(job.Config)
	if err != nil {
		return &worker.Result{
			Status:       "error",
			ErrorMessage: fmt.Sprintf("variable resolution failed: %v", err),
			DurationMS:   0,
		}, nil
	}

	method := getString(resolvedConfig, "method", "GET")
	url := getString(resolvedConfig, "url", "")
	headers := getMap(resolvedConfig, "headers")
	body := resolvedConfig["body"]

	if url == "" {
		return &worker.Result{
			Status:       "error",
			ErrorMessage: "url is required",
			DurationMS:   0,
		}, nil
	}

	loadConfigMap := getMap(resolvedConfig, "load")
	vus := getInt(loadConfigMap, "vus", 1)
	rateLimitRps := getInt(loadConfigMap, "rate_limit_rps", 0)

	totalDuration := getDuration(loadConfigMap, "duration", 10*time.Second)
	rampUp := getDuration(loadConfigMap, "ramp_up", 0)
	rampDown := getDuration(loadConfigMap, "ramp_down", 0)

	assertions := getAssertions(resolvedConfig)

	var bodyBytes []byte
	if body != nil {
		bodyBytes, _ = json.Marshal(body)
	}

	transport := &http.Transport{
		MaxIdleConns:        1000,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:     30 * time.Second,
	}
	client := &http.Client{
		Transport: transport,
		Timeout:   10 * time.Second,
	}

	var limiter *Limiter
	if rateLimitRps > 0 {
		limiter = NewLimiter(float64(rateLimitRps))
	}

	var mu sync.Mutex
	var totalRequests, successCount, errorCount, totalBytes int64
	var allDurations []int64
	statusCodesMap := make(map[int]int)

	var currentSecRequests int64
	var currentSecDurations []int64

	addSample := func(res requestResult) {
		mu.Lock()
		defer mu.Unlock()
		totalRequests++
		if res.err || res.status >= 400 {
			errorCount++
		} else {
			successCount++
		}
		totalBytes += res.bytes
		statusCodesMap[res.status]++
		allDurations = append(allDurations, res.duration.Milliseconds())

		currentSecRequests++
		currentSecDurations = append(currentSecDurations, res.duration.Milliseconds())
	}

	execCtx, execCancel := context.WithTimeout(ctx, totalDuration+2*time.Second)
	defer execCancel()

	startTime := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < vus; i++ {
		wg.Add(1)
		go func(vuID int) {
			defer wg.Done()

			for {
				select {
				case <-execCtx.Done():
					return
				default:
				}

				elapsed := time.Since(startTime)
				if elapsed >= totalDuration {
					return
				}

				var targetVus float64
				if elapsed < rampUp {
					targetVus = float64(vus) * (elapsed.Seconds() / rampUp.Seconds())
				} else if elapsed > totalDuration-rampDown {
					remaining := totalDuration - elapsed
					if remaining < 0 {
						remaining = 0
					}
					targetVus = float64(vus) * (remaining.Seconds() / rampDown.Seconds())
				} else {
					targetVus = float64(vus)
				}

				if float64(vuID) >= targetVus {
					time.Sleep(50 * time.Millisecond)
					continue
				}

				if limiter != nil {
					if err := limiter.Limit(execCtx); err != nil {
						return
					}
				}

				var bodyReader io.Reader
				if bodyBytes != nil {
					bodyReader = bytes.NewReader(bodyBytes)
				}

				req, err := http.NewRequestWithContext(execCtx, method, url, bodyReader)
				if err != nil {
					addSample(requestResult{timestamp: time.Now(), status: 0, duration: 0, bytes: 0, err: true})
					time.Sleep(10 * time.Millisecond)
					continue
				}

				for k, v := range headers {
					if s, ok := v.(string); ok {
						req.Header.Set(k, s)
					}
				}
				if bodyBytes != nil && req.Header.Get("Content-Type") == "" {
					req.Header.Set("Content-Type", "application/json")
				}

				reqStart := time.Now()
				resp, err := client.Do(req)
				reqDuration := time.Since(reqStart)

				if err != nil {
					addSample(requestResult{timestamp: time.Now(), status: 0, duration: reqDuration, bytes: 0, err: true})
				} else {
					bodyData, _ := io.ReadAll(resp.Body)
					resp.Body.Close()
					addSample(requestResult{
						timestamp: time.Now(),
						status:    resp.StatusCode,
						duration:  reqDuration,
						bytes:     int64(len(bodyData)),
						err:       false,
					})
				}
			}
		}(i)
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var timeSeriesRPS []model.TimePoint
	var timeSeriesP95 []model.TimePoint

	collectorDone := make(chan struct{})
	go func() {
		defer close(collectorDone)
		secondCounter := 0
		for {
			select {
			case <-execCtx.Done():
				return
			case <-ticker.C:
				mu.Lock()
				rps := float64(currentSecRequests)
				p95 := computeP95(currentSecDurations)
				currentSecRequests = 0
				currentSecDurations = nil
				mu.Unlock()

				timeSeriesRPS = append(timeSeriesRPS, model.TimePoint{Second: secondCounter, Value: rps})
				timeSeriesP95 = append(timeSeriesP95, model.TimePoint{Second: secondCounter, Value: p95})

				e.streamMetrics(job, secondCounter, rps, p95)
				secondCounter++
			}
		}
	}()

	wg.Wait()
	execCancel() // stops collector ticker
	<-collectorDone

	// Final calculations
	totalTestTime := time.Since(startTime).Seconds()
	if totalTestTime <= 0 {
		totalTestTime = 1.0
	}

	var minLat, maxLat, p50, p90, p95, p99 int64
	var avgLat float64

	n := len(allDurations)
	if n > 0 {
		sort.Slice(allDurations, func(i, j int) bool { return allDurations[i] < allDurations[j] })
		minLat = allDurations[0]
		maxLat = allDurations[n-1]
		
		var sum int64
		for _, d := range allDurations {
			sum += d
		}
		avgLat = float64(sum) / float64(n)

		p50 = allDurations[int(float64(n)*0.50)]
		p90 = allDurations[int(float64(n)*0.90)]
		p95 = allDurations[int(float64(n)*0.95)]
		p99 = allDurations[int(float64(n)*0.99)]
	}

	var errorRate float64
	if totalRequests > 0 {
		errorRate = float64(errorCount) / float64(totalRequests)
	}

	throughputRPS := float64(totalRequests) / totalTestTime

	finalMetrics := model.LoadMetrics{
		TotalRequests: totalRequests,
		SuccessCount:  successCount,
		ErrorCount:    errorCount,
		ErrorRate:     errorRate,
		ThroughputRPS: throughputRPS,
		MinLatencyMS:  minLat,
		MaxLatencyMS:  maxLat,
		AvgLatencyMS:  avgLat,
		P50LatencyMS:  p50,
		P90LatencyMS:  p90,
		P95LatencyMS:  p95,
		P99LatencyMS:  p99,
		TotalBytes:    totalBytes,
		TimeSeriesRPS: timeSeriesRPS,
		TimeSeriesP95: timeSeriesP95,
		StatusCodes:   statusCodesMap,
	}

	allPassed := true
	var assertionResults []worker.Assertion
	for _, assertion := range assertions {
		workerAssertion := runLoadAssertion(assertion, &finalMetrics)
		assertionResults = append(assertionResults, workerAssertion)
		if !workerAssertion.Passed {
			allPassed = false
		}
	}

	status := "passed"
	var errMsg string
	if !allPassed {
		status = "failed"
		var failures []string
		for _, a := range assertionResults {
			if !a.Passed {
				failures = append(failures, a.Message)
			}
		}
		errMsg = strings.Join(failures, "; ")
	}

	metricsBytes, _ := json.Marshal(finalMetrics)
	var metricsMap map[string]any
	_ = json.Unmarshal(metricsBytes, &metricsMap)

	return &worker.Result{
		Status:       status,
		ErrorMessage: errMsg,
		Assertions:   assertionResults,
		RequestData:  map[string]any{"method": method, "url": url},
		ResponseData: map[string]any{"total_requests": finalMetrics.TotalRequests},
		Metrics:      metricsMap,
		DurationMS:   time.Since(startTime).Milliseconds(),
	}, nil
}

func (e *LoadExecutor) streamMetrics(job worker.Job, second int, rps float64, p95 float64) {
	token := e.registry.GetToken()
	if token == "" {
		return
	}

	payload := map[string]any{
		"stream":  "metrics",
		"content": fmt.Sprintf(`{"second":%d,"rps":%f,"p95":%f}`, second, rps, p95),
	}
	body, _ := json.Marshal(payload)

	url := fmt.Sprintf("%s/api/v1/jobs/%s/log", e.cfg.ServerURL, job.JobID)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Worker-Token", token)

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Do(req)
	if err == nil {
		resp.Body.Close()
	}
}

func computeP95(durations []int64) float64 {
	n := len(durations)
	if n == 0 {
		return 0
	}
	durCopy := make([]int64, n)
	copy(durCopy, durations)
	sort.Slice(durCopy, func(i, j int) bool { return durCopy[i] < durCopy[j] })
	index := int(float64(n) * 0.95)
	if index >= n {
		index = n - 1
	}
	return float64(durCopy[index])
}

func runLoadAssertion(assertion assertionDef, metrics *model.LoadMetrics) worker.Assertion {
	a := worker.Assertion{
		Type:     assertion.Type,
		Expected: assertion.Expected,
	}

	switch assertion.Type {
	case "p95_response_time":
		a.Actual = metrics.P95LatencyMS
		if expected, ok := toInt64(assertion.Expected); ok {
			a.Passed = metrics.P95LatencyMS <= expected
		}
		if !a.Passed {
			a.Message = fmt.Sprintf("expected p95 response time <= %vms, got %dms", assertion.Expected, metrics.P95LatencyMS)
		}

	case "p99_response_time":
		a.Actual = metrics.P99LatencyMS
		if expected, ok := toInt64(assertion.Expected); ok {
			a.Passed = metrics.P99LatencyMS <= expected
		}
		if !a.Passed {
			a.Message = fmt.Sprintf("expected p99 response time <= %vms, got %dms", assertion.Expected, metrics.P99LatencyMS)
		}

	case "error_rate":
		a.Actual = metrics.ErrorRate
		if expected, ok := toFloat64(assertion.Expected); ok {
			a.Passed = metrics.ErrorRate <= expected
		}
		if !a.Passed {
			a.Message = fmt.Sprintf("expected error rate <= %v, got %f", assertion.Expected, metrics.ErrorRate)
		}

	case "throughput_min":
		a.Actual = metrics.ThroughputRPS
		if expected, ok := toFloat64(assertion.Expected); ok {
			a.Passed = metrics.ThroughputRPS >= expected
		}
		if !a.Passed {
			a.Message = fmt.Sprintf("expected throughput >= %v RPS, got %f RPS", assertion.Expected, metrics.ThroughputRPS)
		}

	case "avg_response_time":
		a.Actual = metrics.AvgLatencyMS
		if expected, ok := toFloat64(assertion.Expected); ok {
			a.Passed = metrics.AvgLatencyMS <= expected
		}
		if !a.Passed {
			a.Message = fmt.Sprintf("expected average response time <= %vms, got %fms", assertion.Expected, metrics.AvgLatencyMS)
		}

	default:
		a.Passed = true
		a.Message = fmt.Sprintf("unknown load assertion type: %s", assertion.Type)
	}

	return a
}

func parseDuration(val string, defaultVal time.Duration) time.Duration {
	if val == "" {
		return defaultVal
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return defaultVal
	}
	return d
}

func getDuration(m map[string]any, key string, defaultVal time.Duration) time.Duration {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case string:
			return parseDuration(val, defaultVal)
		case float64:
			return time.Duration(val) * time.Millisecond
		}
	}
	return defaultVal
}

func getInt(m map[string]any, key string, defaultVal int) int {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case float64:
			return int(val)
		case int:
			return val
		case int64:
			return int(val)
		}
	}
	return defaultVal
}

func toFloat64(v any) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int64:
		return float64(n), true
	case json.Number:
		f, err := n.Float64()
		return f, err == nil
	}
	return 0, false
}
