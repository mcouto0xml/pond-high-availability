package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// TelemetryNewDataRequest mirrors the API DTO
type TelemetryNewDataRequest struct {
	IotName     string  `json:"iot_name"`
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
	Presence    bool    `json:"presence"`
	Vibration   float64 `json:"vibration"`
	Luminosity  float64 `json:"luminosity"`
	TankLevel   float64 `json:"tank_level"`
}

var iotNames = []string{"sensor-01", "sensor-02", "sensor-03"}

func randomPayload() TelemetryNewDataRequest {
	return TelemetryNewDataRequest{
		IotName:     iotNames[rand.Intn(len(iotNames))],
		Temperature: roundFloat(rand.Float64()*60-10, 2),  // -10°C a 50°C
		Humidity:    roundFloat(rand.Float64()*100, 2),    // 0% a 100%
		Presence:    rand.Intn(2) == 1,
		Vibration:   roundFloat(rand.Float64()*10, 4),     // 0 a 10
		Luminosity:  roundFloat(rand.Float64()*1000, 2),   // 0 a 1000 lux
		TankLevel:   roundFloat(rand.Float64()*100, 2),    // 0% a 100%
	}
}

func roundFloat(val float64, precision int) float64 {
	factor := 1.0
	for i := 0; i < precision; i++ {
		factor *= 10
	}
	return float64(int(val*factor)) / factor
}

// Config holds the load test configuration
type Config struct {
	URL         string
	Concurrency int
	TotalReqs   int
	Timeout     time.Duration
}

// Result holds the result of a single request
type Result struct {
	StatusCode int
	Duration   time.Duration
	Err        error
}

// Stats holds aggregated statistics
type Stats struct {
	mu            sync.Mutex
	totalReqs     int64
	successReqs   int64
	failedReqs    int64
	totalDuration time.Duration
	minDuration   time.Duration
	maxDuration   time.Duration
	statusCodes   map[int]int64
	errors        []string
}

func newStats() *Stats {
	return &Stats{
		statusCodes: make(map[int]int64),
		minDuration: time.Duration(1<<63 - 1),
	}
}

func (s *Stats) record(r Result) {
	s.mu.Lock()
	defer s.mu.Unlock()

	atomic.AddInt64(&s.totalReqs, 1)
	s.totalDuration += r.Duration

	if r.Duration < s.minDuration {
		s.minDuration = r.Duration
	}
	if r.Duration > s.maxDuration {
		s.maxDuration = r.Duration
	}

	if r.Err != nil {
		atomic.AddInt64(&s.failedReqs, 1)
		if len(s.errors) < 10 {
			s.errors = append(s.errors, r.Err.Error())
		}
		return
	}

	s.statusCodes[r.StatusCode]++
	if r.StatusCode >= 200 && r.StatusCode < 300 {
		atomic.AddInt64(&s.successReqs, 1)
	} else {
		atomic.AddInt64(&s.failedReqs, 1)
	}
}

func (s *Stats) print(elapsed time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	total := s.totalReqs
	if total == 0 {
		fmt.Println("No requests were made.")
		return
	}

	avgDuration := s.totalDuration / time.Duration(total)
	rps := float64(total) / elapsed.Seconds()

	fmt.Println()
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║         LOAD TEST RESULTS                ║")
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Printf("\n📊 General\n")
	fmt.Printf("   Total Requests  : %d\n", total)
	fmt.Printf("   Elapsed Time    : %s\n", elapsed.Round(time.Millisecond))
	fmt.Printf("   Req/s (RPS)     : %.2f\n", rps)

	fmt.Printf("\n✅ Results\n")
	fmt.Printf("   Success (2xx)   : %d (%.1f%%)\n", s.successReqs, float64(s.successReqs)/float64(total)*100)
	fmt.Printf("   Failed          : %d (%.1f%%)\n", s.failedReqs, float64(s.failedReqs)/float64(total)*100)

	fmt.Printf("\n⏱️  Latency\n")
	fmt.Printf("   Min             : %s\n", s.minDuration.Round(time.Microsecond))
	fmt.Printf("   Avg             : %s\n", avgDuration.Round(time.Microsecond))
	fmt.Printf("   Max             : %s\n", s.maxDuration.Round(time.Microsecond))

	if len(s.statusCodes) > 0 {
		fmt.Printf("\n📋 Status Codes\n")
		for code, count := range s.statusCodes {
			fmt.Printf("   HTTP %d         : %d\n", code, count)
		}
	}

	if len(s.errors) > 0 {
		fmt.Printf("\n❌ Sample Errors (up to 10)\n")
		for _, e := range s.errors {
			fmt.Printf("   - %s\n", e)
		}
	}

	fmt.Println()
}

func doRequest(client *http.Client, cfg Config) Result {
	start := time.Now()

	payload, err := json.Marshal(randomPayload())
	if err != nil {
		return Result{Err: fmt.Errorf("marshal failed: %w", err), Duration: time.Since(start)}
	}

	req, err := http.NewRequest("POST", cfg.URL, bytes.NewBuffer(payload))
	if err != nil {
		return Result{Err: fmt.Errorf("request creation failed: %w", err), Duration: time.Since(start)}
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return Result{Err: fmt.Errorf("request failed: %w", err), Duration: time.Since(start)}
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)

	return Result{
		StatusCode: resp.StatusCode,
		Duration:   time.Since(start),
	}
}

func main() {
	// Flags
	url := flag.String("url", "http://localhost:8080/telemetry", "Target URL")
	concurrency := flag.Int("c", 10, "Number of concurrent goroutines")
	totalReqs := flag.Int("n", 1000, "Total number of requests")
	timeout := flag.Duration("timeout", 10*time.Second, "HTTP request timeout")
	flag.Parse()

	cfg := Config{
		URL:         *url,
		Concurrency: *concurrency,
		TotalReqs:   *totalReqs,
		Timeout:     *timeout,
	}

	client := &http.Client{
		Timeout: cfg.Timeout,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: cfg.Concurrency,
			DisableKeepAlives:   false,
		},
	}

	fmt.Printf("🚀 Starting load test\n")
	fmt.Printf("   URL          : %s\n", cfg.URL)
	fmt.Printf("   Method       : POST\n")
	fmt.Printf("   Concurrency  : %d goroutines\n", cfg.Concurrency)
	fmt.Printf("   Total Reqs   : %d\n", cfg.TotalReqs)
	fmt.Printf("   Timeout      : %s\n", cfg.Timeout)
	fmt.Printf("   Sensors      : %v\n", iotNames)
	fmt.Println()

	stats := newStats()
	jobs := make(chan struct{}, cfg.TotalReqs)
	var wg sync.WaitGroup

	// Enqueue all jobs
	for i := 0; i < cfg.TotalReqs; i++ {
		jobs <- struct{}{}
	}
	close(jobs)

	start := time.Now()

	// Spawn workers
	for i := 0; i < cfg.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range jobs {
				result := doRequest(client, cfg)
				stats.record(result)
			}
		}()
	}

	// Progress reporter
	done := make(chan struct{})
	go func() {
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				completed := atomic.LoadInt64(&stats.totalReqs)
				fmt.Printf("\r   Progress: %d/%d requests (%.1f%%)", completed, cfg.TotalReqs, float64(completed)/float64(cfg.TotalReqs)*100)
			case <-done:
				return
			}
		}
	}()

	wg.Wait()
	close(done)

	elapsed := time.Since(start)
	stats.print(elapsed)
}