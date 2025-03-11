package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

// Configuration
const (
	DefaultAPIURL = "https://2.intelx.io/"
	DefaultAPIKey = "" // You should set your API key here or use environment variables
	ServerPort    = "8080"
)

// IntelXClient handles communication with the IntelX API
type IntelXClient struct {
	APIURL string
	APIKey string
}

// SearchRequest represents a search request to the IntelX API
type SearchRequest struct {
	Term       string   `json:"term"`
	MaxResults int      `json:"maxresults"`
	Media      int      `json:"media"`
	Sort       int      `json:"sort"`
	Terminate  []string `json:"terminate"`
	DateFrom   string   `json:"datefrom,omitempty"` // Using omitempty to ensure it's not included if empty
	DateTo     string   `json:"dateto,omitempty"`   // Using omitempty to ensure it's not included if empty
}

// SearchResponse represents the response from a search request
type SearchResponse struct {
	ID     string `json:"id"`
	Status int    `json:"status"`
}

// SearchResult represents the result of a search
type SearchResult struct {
	ID      string   `json:"id"`
	Records []Record `json:"records"`
	Status  int      `json:"status"`
}

// Record represents a single record in the search results
type Record struct {
	SystemID   string `json:"systemid"`
	StorageID  string `json:"storageid"`
	Name       string `json:"name"`
	Date       string `json:"date"`
	Media      int    `json:"media"`
	Bucket     string `json:"bucket"`
	XScore     int    `json:"xscore"`
	Similarity int    `json:"similarity"`
}

// TemplateData holds data to be passed to HTML templates
type TemplateData struct {
	Title       string
	Error       string
	SearchTerm  string
	Results     *SearchResult
	Preview     string
	RecordID    string
	StorageID   string
	Bucket      string
	APIKey      string
	APIURL      string
	CurrentYear int
	DateFrom    string
	DateTo      string
}

// NewIntelXClient creates a new IntelX API client
func NewIntelXClient() *IntelXClient {
	apiURL := os.Getenv("INTELX_API_URL")
	if apiURL == "" {
		apiURL = DefaultAPIURL
	}

	// Ensure the API URL ends with a slash
	if !strings.HasSuffix(apiURL, "/") {
		apiURL = apiURL + "/"
	}

	apiKey := os.Getenv("INTELX_API_KEY")
	if apiKey == "" {
		apiKey = DefaultAPIKey
	}

	return &IntelXClient{
		APIURL: apiURL,
		APIKey: apiKey,
	}
}

// Search initiates a search with the IntelX API
func (c *IntelXClient) Search(ctx context.Context, term string, maxResults int, dateFrom, dateTo string) (*SearchResponse, error) {
	searchReq := SearchRequest{
		Term:       term,
		MaxResults: maxResults,
		Media:      0,
		Sort:       2,
		Terminate:  []string{},
	}

	// Add date range if provided
	if dateFrom != "" {
		searchReq.DateFrom = dateFrom
		log.Printf("Adding dateFrom filter: %s", dateFrom)
	}

	if dateTo != "" {
		searchReq.DateTo = dateTo
		log.Printf("Adding dateTo filter: %s", dateTo)
	}

	reqBody, err := json.Marshal(searchReq)
	if err != nil {
		return nil, fmt.Errorf("error marshaling search request: %w", err)
	}

	log.Printf("Sending search request: %s", string(reqBody))

	req, err := http.NewRequestWithContext(ctx, "POST", c.APIURL+"intelligent/search", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating search request: %w", err)
	}

	req.Header.Set("x-key", c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing search request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s, status code: %d", string(body), resp.StatusCode)
	}

	var searchResp SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&searchResp); err != nil {
		return nil, fmt.Errorf("error decoding search response: %w", err)
	}

	return &searchResp, nil
}

// GetResults retrieves search results from the IntelX API
func (c *IntelXClient) GetResults(ctx context.Context, searchID string, limit int) (*SearchResult, error) {
	url := fmt.Sprintf("%sintelligent/search/result?id=%s&limit=%d&statistics=1&previewlines=8", c.APIURL, searchID, limit)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating get results request: %w", err)
	}

	req.Header.Set("x-key", c.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing get results request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: %s, status code: %d", string(body), resp.StatusCode)
	}

	var result SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding search results: %w", err)
	}

	return &result, nil
}

// GetFilePreview retrieves a file preview from the IntelX API
func (c *IntelXClient) GetFilePreview(ctx context.Context, storageID, bucket string) (string, error) {
	url := fmt.Sprintf("%sfile/preview?sid=%s&f=0&l=8&c=1&m=1&b=%s&k=%s", c.APIURL, storageID, bucket, c.APIKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating file preview request: %w", err)
	}

	req.Header.Set("x-key", c.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error executing file preview request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error: %s, status code: %d", string(body), resp.StatusCode)
	}

	preview, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading file preview: %w", err)
	}

	return string(preview), nil
}

// GetFileContent retrieves the full file content from the IntelX API
func (c *IntelXClient) GetFileContent(ctx context.Context, systemID string, bucket string) ([]byte, error) {
	// Based on Python FILE_READ(systemid, 0, bucket, filename)
	// Note: The API expects systemid to be passed differently than in the preview endpoint
	url := fmt.Sprintf("%sfile/read?type=0&systemid=%s&bucket=%s", c.APIURL, systemID, bucket)

	log.Printf("Download URL: %s", url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating file content request: %w", err)
	}

	req.Header.Set("x-key", c.APIKey)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing file content request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Error response: %s", string(body))
		return nil, fmt.Errorf("API error: %s, status code: %d", string(body), resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading file content: %w", err)
	}

	return content, nil
}

// GetFileContentAlternative provides an alternative method for downloading files
func (c *IntelXClient) GetFileContentAlternative(ctx context.Context, storageID string, bucket string) ([]byte, error) {
	// Try with storageID instead of systemID
	url := fmt.Sprintf("%sfile/read?storageid=%s&bucket=%s", c.APIURL, storageID, bucket)

	log.Printf("Alternative download URL: %s", url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating file content request: %w", err)
	}

	req.Header.Set("x-key", c.APIKey)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing file content request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Error response: %s", string(body))
		return nil, fmt.Errorf("API error: %s, status code: %d", string(body), resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading file content: %w", err)
	}

	return content, nil
}

// GetFileContentThirdMethod tries a third method for downloading files
func (c *IntelXClient) GetFileContentThirdMethod(ctx context.Context, storageID, formatType, bucket string) ([]byte, error) {
	// Try a third format based on the HTML example
	url := fmt.Sprintf("%sfile/view?f=%s&storageid=%s&k=%s&m=1&b=%s", c.APIURL, formatType, storageID, c.APIKey, bucket)

	log.Printf("Third method download URL: %s", url)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating file content request: %w", err)
	}

	req.Header.Set("x-key", c.APIKey)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing file content request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		log.Printf("Error response: %s", string(body))
		return nil, fmt.Errorf("API error: %s, status code: %d", string(body), resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading file content: %w", err)
	}

	return content, nil
}

// renderTemplate renders an HTML template with the given data
func renderTemplate(w http.ResponseWriter, tmpl string, data *TemplateData) {
	// Define template functions
	funcMap := template.FuncMap{
		"or": func(a, b string) string {
			if a != "" {
				return a
			}
			return b
		},
	}

	// Parse templates with function map
	t, err := template.New("layout.html").Funcs(funcMap).ParseFiles("templates/layout.html", "templates/"+tmpl+".html")
	if err != nil {
		log.Printf("Error parsing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Add current year to template data
	data.CurrentYear = time.Now().Year()

	err = t.ExecuteTemplate(w, "layout", data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// loggingMiddleware logs HTTP requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}

func main() {
	// Create a new IntelX client
	client := NewIntelXClient()

	// Check if API key is set
	if client.APIKey == "" {
		log.Println("Warning: No API key set. Please set your API key in the environment variable INTELX_API_KEY or in the DefaultAPIKey constant.")
	}

	// Create a new ServeMux
	mux := http.NewServeMux()

	// Serve static files
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Home page handler
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "index", &TemplateData{
			Title:  "IntelX Web Interface",
			APIURL: client.APIURL,
			APIKey: client.APIKey,
		})
	})

	// Search handler
	mux.HandleFunc("POST /search", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		searchTerm := r.FormValue("search")
		if searchTerm == "" {
			renderTemplate(w, "index", &TemplateData{
				Title: "IntelX Web Interface",
				Error: "Please enter a search term",
			})
			return
		}

		maxResults := 10
		if maxResultsStr := r.FormValue("maxresults"); maxResultsStr != "" {
			if val, err := strconv.Atoi(maxResultsStr); err == nil && val > 0 {
				maxResults = val
			}
		}

		// Check if date range filtering is enabled
		dateRangeEnabled := r.FormValue("date_range_enabled") == "on"

		// Initialize date variables
		var dateFrom, dateTo string

		// Only process date parameters if date range is enabled
		if dateRangeEnabled {
			// Get date range parameters if the toggle is enabled
			startDate := r.FormValue("start_date")
			endDate := r.FormValue("end_date")

			log.Printf("Date range filtering is enabled")
			log.Printf("Raw input dates - startDate: %s, endDate: %s", startDate, endDate)

			// Convert HTML datetime-local format to expected API format (YYYY-MM-DD HH:MM:SS)
			if startDate != "" {
				// Parse the HTML datetime-local format (YYYY-MM-DDTHH:MM)
				t, err := time.Parse("2006-01-02T15:04", startDate)
				if err == nil {
					// Format to API expected format
					dateFrom = t.Format("2006-01-02 15:04:05")
					log.Printf("Start date parsed successfully: %s", dateFrom)
				} else {
					log.Printf("Error parsing start date: %v", err)
				}
			}

			if endDate != "" {
				// Parse the HTML datetime-local format (YYYY-MM-DDTHH:MM)
				t, err := time.Parse("2006-01-02T15:04", endDate)
				if err == nil {
					// Format to API expected format
					dateTo = t.Format("2006-01-02 15:04:05")
					log.Printf("End date parsed successfully: %s", dateTo)
				} else {
					log.Printf("Error parsing end date: %v", err)
				}
			}
		} else {
			log.Printf("Date range filtering is disabled - ignoring any date values in the form")
		}

		log.Printf("Searching with term: %s, maxResults: %d, dateRangeEnabled: %v, dateFrom: %s, dateTo: %s",
			searchTerm, maxResults, dateRangeEnabled, dateFrom, dateTo)

		ctx := r.Context()
		searchResp, err := client.Search(ctx, searchTerm, maxResults, dateFrom, dateTo)
		if err != nil {
			log.Printf("Search error: %v", err)
			renderTemplate(w, "index", &TemplateData{
				Title:      "IntelX Web Interface",
				Error:      "Error performing search: " + err.Error(),
				SearchTerm: searchTerm,
				DateFrom:   r.FormValue("start_date"), // Return the original input date format
				DateTo:     r.FormValue("end_date"),   // Return the original input date format
			})
			return
		}

		// Redirect to results page with date range parameters only if date range is enabled
		redirectURL := fmt.Sprintf("/results?id=%s&term=%s", searchResp.ID, url.QueryEscape(searchTerm))
		if dateRangeEnabled && (dateFrom != "" || dateTo != "") {
			if dateFrom != "" {
				redirectURL += fmt.Sprintf("&from=%s", url.QueryEscape(r.FormValue("start_date")))
			}
			if dateTo != "" {
				redirectURL += fmt.Sprintf("&to=%s", url.QueryEscape(r.FormValue("end_date")))
			}
		}

		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	})

	// Results handler
	mux.HandleFunc("GET /results", func(w http.ResponseWriter, r *http.Request) {
		searchID := r.URL.Query().Get("id")
		searchTerm := r.URL.Query().Get("term")
		fromDate := r.URL.Query().Get("from")
		toDate := r.URL.Query().Get("to")

		if searchID == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		limit := 10
		if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
			if val, err := strconv.Atoi(limitStr); err == nil && val > 0 {
				limit = val
			}
		}

		ctx := r.Context()
		results, err := client.GetResults(ctx, searchID, limit)
		if err != nil {
			log.Printf("Get results error: %v", err)
			renderTemplate(w, "index", &TemplateData{
				Title:      "IntelX Web Interface",
				Error:      "Error retrieving results: " + err.Error(),
				SearchTerm: searchTerm,
				DateFrom:   fromDate,
				DateTo:     toDate,
			})
			return
		}

		// Log the date parameters being passed to the template
		log.Printf("Rendering results template with dateFrom: %s, dateTo: %s", fromDate, toDate)

		renderTemplate(w, "results", &TemplateData{
			Title:      "Search Results - IntelX Web Interface",
			SearchTerm: searchTerm,
			Results:    results,
			APIURL:     client.APIURL,
			APIKey:     client.APIKey,
			DateFrom:   fromDate,
			DateTo:     toDate,
		})
	})

	// Preview handler
	mux.HandleFunc("GET /preview", func(w http.ResponseWriter, r *http.Request) {
		storageID := r.URL.Query().Get("sid")
		recordID := r.URL.Query().Get("rid")
		bucket := r.URL.Query().Get("bucket")

		if storageID == "" || bucket == "" {
			http.Error(w, "Missing required parameters", http.StatusBadRequest)
			return
		}

		ctx := r.Context()
		preview, err := client.GetFilePreview(ctx, storageID, bucket)
		if err != nil {
			log.Printf("Get preview error: %v", err)
			http.Error(w, "Error retrieving preview: "+err.Error(), http.StatusInternalServerError)
			return
		}

		renderTemplate(w, "preview", &TemplateData{
			Title:     "File Preview - IntelX Web Interface",
			Preview:   preview,
			RecordID:  recordID,
			StorageID: storageID,
			Bucket:    bucket,
		})
	})

	// Download handler
	mux.HandleFunc("GET /download", func(w http.ResponseWriter, r *http.Request) {
		// Get query parameters
		systemID := r.URL.Query().Get("rid")  // system ID (from the Python example)
		storageID := r.URL.Query().Get("sid") // storage ID (sometimes used instead)
		bucket := r.URL.Query().Get("bucket") // bucket parameter is required
		filename := r.URL.Query().Get("name") // optional filename

		// Check required parameters
		if systemID == "" && storageID == "" {
			http.Error(w, "Missing system ID (rid) or storage ID (sid)", http.StatusBadRequest)
			return
		}

		if bucket == "" {
			bucket = "pastes" // Default bucket if not specified
		}

		if filename == "" {
			filename = "intelx-file"
		}

		// For debugging
		if systemID != "" {
			log.Printf("Attempting download with SystemID=%s, Bucket=%s, Filename=%s", systemID, bucket, filename)
		} else {
			log.Printf("Attempting download with StorageID=%s, Bucket=%s, Filename=%s", storageID, bucket, filename)
		}

		ctx := r.Context()
		var content []byte
		var err error

		// Try several approaches in sequence

		// Approach 1: Use systemID
		if systemID != "" {
			log.Printf("Trying primary download method with systemID")
			content, err = client.GetFileContent(ctx, systemID, bucket)
			if err == nil {
				log.Printf("Primary download method successful!")
			} else {
				log.Printf("Primary download method failed: %v", err)
			}
		}

		// Approach 2: Use storageID if available and Approach 1 failed
		if (err != nil || content == nil) && storageID != "" {
			log.Printf("Trying alternative download method with storageID")
			content, err = client.GetFileContentAlternative(ctx, storageID, bucket)
			if err == nil {
				log.Printf("Alternative download method successful!")
			} else {
				log.Printf("Alternative download method also failed: %v", err)
			}
		}

		// Approach 3: Try systemID with storageID approach if all else failed
		if err != nil && systemID != "" && storageID == "" {
			log.Printf("Trying to use systemID as storageID")
			content, err = client.GetFileContentAlternative(ctx, systemID, bucket)
			if err == nil {
				log.Printf("Using systemID as storageID was successful!")
			} else {
				log.Printf("Using systemID as storageID also failed: %v", err)
			}
		}

		// Try various format types with the third method
		formatTypes := []string{"0", "2", "3", "5", "7"}
		if err != nil && storageID != "" {
			for _, formatType := range formatTypes {
				log.Printf("Trying third method with format type %s", formatType)
				content, err = client.GetFileContentThirdMethod(ctx, storageID, formatType, bucket)
				if err == nil {
					log.Printf("Third method successful with format type %s!", formatType)
					break
				} else {
					log.Printf("Third method failed with format type %s: %v", formatType, err)
				}
			}
		}

		// If all download attempts failed, redirect to IntelX
		if err != nil {
			log.Printf("All download attempts failed, redirecting to IntelX website")
			if systemID != "" {
				http.Redirect(w, r, fmt.Sprintf("https://intelx.io/?did=%s", systemID), http.StatusFound)
			} else {
				http.Redirect(w, r, fmt.Sprintf("https://intelx.io/?did=%s", storageID), http.StatusFound)
			}
			return
		}

		log.Printf("Successfully downloaded file. Size: %d bytes", len(content))

		// Set appropriate content type based on filename
		contentType := "application/octet-stream"
		if strings.HasSuffix(strings.ToLower(filename), ".txt") {
			contentType = "text/plain"
		} else if strings.HasSuffix(strings.ToLower(filename), ".json") {
			contentType = "application/json"
		} else if strings.HasSuffix(strings.ToLower(filename), ".html") || strings.HasSuffix(strings.ToLower(filename), ".htm") {
			contentType = "text/html"
		}

		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		w.Header().Set("Content-Type", contentType)
		w.Header().Set("Content-Length", strconv.Itoa(len(content)))
		w.Write(content)
	})

	// API key update handler
	mux.HandleFunc("POST /update-api-key", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		apiKey := r.FormValue("api_key")
		apiURL := r.FormValue("api_url")

		if apiKey != "" {
			client.APIKey = apiKey
		}

		if apiURL != "" {
			// Ensure the API URL ends with a slash
			if !strings.HasSuffix(apiURL, "/") {
				apiURL = apiURL + "/"
			}
			client.APIURL = apiURL
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
	})

	// Apply middleware
	handler := loggingMiddleware(mux)

	// Start the server
	log.Printf("Starting server on port %s...", ServerPort)
	log.Fatal(http.ListenAndServe(":"+ServerPort, handler))
}
