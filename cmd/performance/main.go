package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Configuration
const (
	APIURL       = "http://localhost:8080/v1/products"
	OutputCSV    = "././docs/assets/cache_performance.csv"
	StatsCSV     = "././docs/assets/cache_statistics.csv"
	PageSize     = 5
	MaxPageCount = 10
)

type Result struct {
	PageNumber      int
	PageSize        int
	ColdCacheTimeMS int
	WarmCacheTimeMS int
}

type Pagination struct {
	Next string `json:"next"`
}

type APIResponse struct {
	Pagination Pagination `json:"pagination"`
}

func fetchPageTimes() ([]Result, []int, []int) {
	var (
		results        []Result
		coldCacheTimes []int
		warmCacheTimes []int
		nextCursor     string
		pageNumber     int
	)

	for {
		pageNumber++

		url := fmt.Sprintf("%s?pageSize=%d&next=%s", APIURL, PageSize, nextCursor)

		// Cold cache request
		startCold := time.Now()
		respCold, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			break
		}
		elapsedCold := int(time.Since(startCold).Milliseconds())

		if respCold.StatusCode != http.StatusOK {
			fmt.Printf("Error: Failed to fetch data (HTTP %d)\n", respCold.StatusCode)
			break
		}

		// Warm cache request
		startWarm := time.Now()
		respWarm, err := http.Get(url)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			break
		}
		elapsedWarm := int(time.Since(startWarm).Milliseconds())

		// Log response times
		fmt.Printf("Page %d: Cold cache %dms, Warm cache %dms\n", pageNumber, elapsedCold, elapsedWarm)
		results = append(results, Result{
			PageNumber:      pageNumber,
			PageSize:        PageSize,
			ColdCacheTimeMS: elapsedCold,
			WarmCacheTimeMS: elapsedWarm,
		})
		coldCacheTimes = append(coldCacheTimes, elapsedCold)
		warmCacheTimes = append(warmCacheTimes, elapsedWarm)

		// Parse next cursor
		var apiResp APIResponse
		if err := json.NewDecoder(respCold.Body).Decode(&apiResp); err != nil {
			fmt.Printf("Error decoding response: %v\n", err)
			break
		}
		nextCursor = apiResp.Pagination.Next

		respCold.Body.Close()
		respWarm.Body.Close()

		// Stop if there's no next page
		if nextCursor == "" || pageNumber >= MaxPageCount {
			break
		}
	}

	return results, coldCacheTimes, warmCacheTimes
}

func saveResultsToCSV(results []Result, coldCacheTimes, warmCacheTimes []int) {
	// Write performance results to CSV
	file, err := os.Create(OutputCSV)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"page_number", "page_size", "cold_cache_time_ms", "warm_cache_time_ms"}
	writer.Write(headers)

	for _, result := range results {
		writer.Write([]string{
			strconv.Itoa(result.PageNumber),
			strconv.Itoa(result.PageSize),
			strconv.Itoa(result.ColdCacheTimeMS),
			strconv.Itoa(result.WarmCacheTimeMS),
		})
	}
	fmt.Printf("Cache performance results saved to %s\n", OutputCSV)

	// Calculate statistics
	stats := map[string]map[string]float64{
		"cold_cache": calculateStats(coldCacheTimes),
		"warm_cache": calculateStats(warmCacheTimes),
	}

	// Save statistics to CSV
	statsFile, err := os.Create(StatsCSV)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer statsFile.Close()

	statsWriter := csv.NewWriter(statsFile)
	defer statsWriter.Flush()

	statsWriter.Write([]string{"Cache Type", "Mean (ms)", "Median (ms)", "Max (ms)", "Min (ms)"})
	statsWriter.Write([]string{"Cold Cache", fmt.Sprintf("%.2f", stats["cold_cache"]["mean"]), fmt.Sprintf("%.2f", stats["cold_cache"]["median"]), fmt.Sprintf("%.2f", stats["cold_cache"]["max"]), fmt.Sprintf("%.2f", stats["cold_cache"]["min"])})
	statsWriter.Write([]string{"Warm Cache", fmt.Sprintf("%.2f", stats["warm_cache"]["mean"]), fmt.Sprintf("%.2f", stats["warm_cache"]["median"]), fmt.Sprintf("%.2f", stats["warm_cache"]["max"]), fmt.Sprintf("%.2f", stats["warm_cache"]["min"])})

	fmt.Printf("Cache performance statistics saved to %s\n", StatsCSV)
}

func calculateStats(times []int) map[string]float64 {
	total := 0
	for _, t := range times {
		total += t
	}
	mean := float64(total) / float64(len(times))

	// Calculate median
	median := 0.0
	if len(times)%2 == 0 {
		median = float64(times[len(times)/2-1]+times[len(times)/2]) / 2.0
	} else {
		median = float64(times[len(times)/2])
	}

	// Find min and max
	min, max := times[0], times[0]
	for _, t := range times {
		if t < min {
			min = t
		}
		if t > max {
			max = t
		}
	}

	return map[string]float64{
		"mean":   mean,
		"median": median,
		"max":    float64(max),
		"min":    float64(min),
	}
}

func main() {
	fmt.Println("Starting cache performance check...")
	results, coldCacheTimes, warmCacheTimes := fetchPageTimes()
	saveResultsToCSV(results, coldCacheTimes, warmCacheTimes)
	fmt.Println("Cache performance check completed.")
}
