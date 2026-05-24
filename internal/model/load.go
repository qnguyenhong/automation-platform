package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type LoadMetrics struct {
	TotalRequests int64            `json:"total_requests"`
	SuccessCount  int64            `json:"success_count"`
	ErrorCount    int64            `json:"error_count"`
	ErrorRate     float64          `json:"error_rate"`
	ThroughputRPS float64          `json:"throughput_rps"`
	MinLatencyMS  int64            `json:"min_latency_ms"`
	MaxLatencyMS  int64            `json:"max_latency_ms"`
	AvgLatencyMS  float64          `json:"avg_latency_ms"`
	P50LatencyMS  int64            `json:"p50_latency_ms"`
	P90LatencyMS  int64            `json:"p90_latency_ms"`
	P95LatencyMS  int64            `json:"p95_latency_ms"`
	P99LatencyMS  int64            `json:"p99_latency_ms"`
	TotalBytes    int64            `json:"total_bytes"`
	TimeSeriesRPS []TimePoint      `json:"time_series_rps"`
	TimeSeriesP95 []TimePoint      `json:"time_series_p95"`
	StatusCodes   map[int]int      `json:"status_codes"`
}

type TimePoint struct {
	Second int     `json:"second"`
	Value  float64 `json:"value"`
}

type LoadMetricsDb struct {
	ID            uuid.UUID       `json:"id"`
	ResultID      uuid.UUID       `json:"result_id"`
	RunID         uuid.UUID       `json:"run_id"`
	TotalRequests int64           `json:"total_requests"`
	SuccessCount  int64           `json:"success_count"`
	ErrorCount    int64           `json:"error_count"`
	ErrorRate     float64         `json:"error_rate"`
	ThroughputRPS float64         `json:"throughput_rps"`
	MinLatencyMS  int64           `json:"min_latency_ms"`
	MaxLatencyMS  int64           `json:"max_latency_ms"`
	AvgLatencyMS  float64         `json:"avg_latency_ms"`
	P50LatencyMS  int64           `json:"p50_latency_ms"`
	P90LatencyMS  int64           `json:"p90_latency_ms"`
	P95LatencyMS  int64           `json:"p95_latency_ms"`
	P99LatencyMS  int64           `json:"p99_latency_ms"`
	TotalBytes    int64           `json:"total_bytes"`
	TimeSeries    json.RawMessage `json:"time_series"`
	StatusCodes   json.RawMessage `json:"status_codes"`
	CreatedAt     time.Time       `json:"created_at"`
}
