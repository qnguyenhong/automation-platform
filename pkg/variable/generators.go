package variable

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	firstNames = []string{"James", "Mary", "John", "Patricia", "Robert", "Jennifer", "Michael", "Linda", "William", "Elizabeth", "David", "Barbara", "Richard", "Susan", "Joseph", "Jessica", "Thomas", "Sarah", "Christopher", "Karen"}
	lastNames  = []string{"Smith", "Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis", "Rodriguez", "Martinez", "Hernandez", "Lopez", "Gonzalez", "Wilson", "Anderson", "Thomas", "Taylor", "Moore", "Jackson", "Martin"}
	domains    = []string{"example.com", "test.com", "mail.com", "demo.org", "sample.io"}
)

func randomUUID() string {
	return uuid.New().String()
}

func randomEmail() string {
	prefix := randomString(8)
	domain := domains[randomIntN(len(domains))]
	return strings.ToLower(prefix + "@" + domain)
}

func randomString(length int) string {
	if length <= 0 {
		length = 8
	}
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[randomIntN(len(chars))]
	}
	return string(b)
}

func randomInt(min, max int) int {
	if min >= max {
		return min
	}
	return min + randomIntN(max-min+1)
}

func randomIntN(n int) int {
	if n <= 0 {
		return 0
	}
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		return 0
	}
	return int(nBig.Int64())
}

func randomFloat(min, max float64, decimals int) float64 {
	if min >= max {
		return min
	}
	nBig, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return min
	}
	ratio := float64(nBig.Int64()) / 1000000.0
	val := min + ratio*(max-min)
	factor := 1.0
	for i := 0; i < decimals; i++ {
		factor *= 10
	}
	return float64(int(val*factor)) / factor
}

func randomBool() bool {
	return randomIntN(2) == 1
}

func randomName() string {
	first := firstNames[randomIntN(len(firstNames))]
	last := lastNames[randomIntN(len(lastNames))]
	return first + " " + last
}

func randomFirstName() string {
	return firstNames[randomIntN(len(firstNames))]
}

func randomLastName() string {
	return lastNames[randomIntN(len(lastNames))]
}

func randomPhone() string {
	area := randomInt(200, 999)
	p1 := randomInt(200, 999)
	p2 := randomInt(1000, 9999)
	return fmt.Sprintf("+1-%d-%d-%d", area, p1, p2)
}

func randomIPv4() string {
	return fmt.Sprintf("%d.%d.%d.%d", randomInt(1, 254), randomInt(0, 254), randomInt(0, 254), randomInt(1, 254))
}

func randomIPv6() string {
	parts := make([]string, 8)
	for i := range parts {
		parts[i] = fmt.Sprintf("%04x", randomInt(0, 65535))
	}
	return strings.Join(parts, ":")
}

func randomEnum(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return values[randomIntN(len(values))]
}

func timestampSeconds() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

func timestampMS() string {
	return fmt.Sprintf("%d", time.Now().UnixMilli())
}

func isoDate() string {
	return time.Now().UTC().Format(time.RFC3339)
}

func dateWithFormat(format string) string {
	if format == "" {
		format = "2006-01-02"
	}
	return time.Now().UTC().Format(format)
}

func randomHex(length int) string {
	if length <= 0 {
		length = 16
	}
	b := make([]byte, (length+1)/2)
	rand.Read(b)
	return hex.EncodeToString(b)[:length]
}

func randomAlpha(length int) string {
	if length <= 0 {
		length = 8
	}
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = chars[randomIntN(len(chars))]
	}
	return string(b)
}

func randomDigits(length int) string {
	if length <= 0 {
		length = 6
	}
	b := make([]byte, length)
	for i := range b {
		b[i] = byte('0' + randomIntN(10))
	}
	return string(b)
}

func randomSentence(wordCount int) string {
	if wordCount <= 0 {
		wordCount = 6
	}
	words := []string{"the", "quick", "brown", "fox", "jumps", "over", "lazy", "dog", "a", "an", "is", "was", "has", "had", "will", "would", "could", "should", "may", "might", "can", "shall", "do", "does", "did", "be", "been", "being", "have", "having", "go", "going", "gone", "went", "come", "coming", "came", "take", "taking", "took", "give", "giving", "gave", "make", "making", "made", "know", "knowing", "knew", "think", "thinking", "thought"}
	result := make([]string, wordCount)
	for i := range result {
		result[i] = words[randomIntN(len(words))]
	}
	return strings.Join(result, " ")
}

func randomURL() string {
	domain := domains[randomIntN(len(domains))]
	path := randomAlpha(8)
	return fmt.Sprintf("https://%s/%s", domain, strings.ToLower(path))
}
