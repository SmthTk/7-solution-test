package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

const (
	PORT = "5051"
)

type BeefSummary struct {
	Beef map[string]int `json:"beef"`
}

// Cache for avoid repeated API calls
var (
	beefDataCache       string
	beefDataCacheMux    sync.RWMutex
	summaryCache        *BeefSummary
	summaryCacheMux     sync.RWMutex
	cacheInitialized    bool
	cacheInitializedMux sync.RWMutex
)

func main() {
	// Initialize the cache
	go initializeCache()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Define routes
	router.GET("/beef/summary", beefSummaryHandler)

	// Start the server
	log.Println(fmt.Sprintf("Server running on port %s...", PORT))
	log.Fatal(router.Run(":" + PORT))
}

func initializeCache() {
	fetchDataAndUpdateCache()
	cacheInitializedMux.Lock()
	cacheInitialized = true
	cacheInitializedMux.Unlock()
}

// Check if cache is initialized
func isCacheInitialized() bool {
	cacheInitializedMux.RLock()
	defer cacheInitializedMux.RUnlock()
	return cacheInitialized
}

// Handler for the beef summary endpoint using Gin
func beefSummaryHandler(c *gin.Context) {
	if !isCacheInitialized() {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "Service initializing, please try again shortly",
		})
		return
	}

	// Return cached summary
	summaryCacheMux.RLock()
	cachedSummary := summaryCache
	summaryCacheMux.RUnlock()

	c.JSON(http.StatusOK, cachedSummary)
}

func fetchDataAndUpdateCache() {
	// Fetch data from the API
	resp, err := http.Get("https://baconipsum.com/api/?type=meat-and-filler&paras=99&format=text")
	if err != nil {
		log.Printf("Error fetching data from Bacon Ipsum API: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response body: %v", err)
		return
	}

	// Update the beef data cache
	beefDataCacheMux.Lock()
	beefDataCache = string(body)
	beefDataCacheMux.Unlock()

	// Process the data and update the summary cache
	updateSummaryCache()
}

// Process the beef data and update the summary cache
func updateSummaryCache() {
	beefDataCacheMux.RLock()
	data := beefDataCache
	beefDataCacheMux.RUnlock()

	// ลบอักขระ
	re := regexp.MustCompile(`[,.;!?]`)
	data = re.ReplaceAllString(data, " ")

	// lower case
	data = strings.ToLower(data)

	// Split space
	words := strings.Fields(data)

	// Count occurrences of each beef type
	beefCounts := make(map[string]int)
	for _, word := range words {
		beefCounts[word]++
	}

	// Create the summary
	summary := &BeefSummary{
		Beef: beefCounts,
	}

	// Update the summary cache
	summaryCacheMux.Lock()
	summaryCache = summary
	summaryCacheMux.Unlock()
}
