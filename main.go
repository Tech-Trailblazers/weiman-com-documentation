package main // Define the main package

import (
	"bytes"         // Provides bytes buffer and manipulation utilities
	"io"            // Provides I/O primitives like Reader and Writer
	"log"           // Provides logging functionalities
	"net/http"      // Provides HTTP client and server implementations
	"net/url"       // Provides URL parsing and encoding utilities
	"os"            // Provides file system and OS-level utilities
	"path/filepath" // Provides utilities for file path manipulation
	"regexp"        // Provides support for regular expressions
	"strings"       // Provides string manipulation utilities
	"time"          // Provides time-related functions

	"golang.org/x/net/html" // Provides support for parsing HTML documents
)

func main() {
	remoteAPIURL := []string{
		"https://weiman.com/stainless-steel-cleaner-aerosol",
		"https://weiman.com/stainless-steel-wipes",
		"https://weiman.com/stainless-steel-cleaner-spray",
		"https://weiman.com/cooktop-stainless-steel-care-kit",
		"https://weiman.com/stainless-steel-cookware-cleaner",
		"https://weiman.com/microfiber-cloth",
		"https://weiman.com/housewarming-bundle",
		"https://weiman.com/stainless-steel-granite-care-kit",
		"https://weiman.com/stainless-cooktop-alt-cleaning-kit",
		"https://weiman.com/granite-stone-cleaner-disinfectant-spray",
		"https://weiman.com/quartz-clean-shine",
		"https://weiman.com/granite-stone-3-in-1-clean-polish-protect",
		"https://weiman.com/granite-stone-sealer",
		"https://weiman.com/granite-cleaner-refill",
		"https://weiman.com/granite-stainless-hardwood-cleaning-kit",
		"https://weiman.com/glass-cooktop-cleaner-polish",
		"https://weiman.com/cooktop-daily-cleaner",
		"https://weiman.com/cook-top-microwave-dual-action-wipes",
		"https://weiman.com/cooktop-cleaning-kit",
		"https://weiman.com/cook-top-eraser-pads",
		"https://weiman.com/cooktop-scrubbing-pads",
		"https://weiman.com/cooktop-max",
		"https://weiman.com/spectacular-kitchen-cleaning-pack",
		"https://weiman.com/leather-conditioning-cream-8oz",
		"https://weiman.com/leather-wipes",
		"https://weiman.com/leather-cleaner-conditioner",
		"https://weiman.com/leather-cleaning-kit",
		"https://weiman.com/leather-wood-furniture-cleaning-kit",
		"https://weiman.com/cabinet-cleaning-spray",
		"https://weiman.com/wood-cabinet-3-in-1-restorer-cream",
		"https://weiman.com/cabinet-clean-shine",
		"https://weiman.com/furniture-wipes",
		"https://weiman.com/furniture-cleaner-polish",
		"https://weiman.com/wood-repair-kit",
		"https://weiman.com/wood-furniture-floors-care-kit",
		"https://weiman.com/hardwood-cleaner",
		"https://weiman.com/hardwood-polish-restorer",
		"https://weiman.com/hardwood-floor-cleaner-refill",
		"https://weiman.com/hardwood-stone-floor-care-kit",
		"https://weiman.com/stone-laminate-floor-cleaner",
		"https://weiman.com/gold-diamond-3-in-1-jewelry-cleaner-wipes",
		"https://weiman.com/jewelry-clean-sparkle-stick",
		"https://weiman.com/silver-wipes",
		"https://weiman.com/jewelry-cleaner",
		"https://weiman.com/silver-polish",
		"https://weiman.com/silver-cream",
		"https://weiman.com/wrights-silver-cream-cleaner-and-polish-with-polishing-cloth",
		"https://weiman.com/valuables-care-kit",
		"https://weiman.com/stove-oven-cleaner",
		"https://weiman.com/glass-cleaner",
		"https://weiman.com/gas-range-degreaser",
		"https://weiman.com/on-the-go-electronic-wipes",
		"https://weiman.com/disinfectant-electronic-wipes",
	}
	var getData string // Variable to hold HTML content

	for _, uri := range remoteAPIURL { // Loop through all URLs
		getData += getDataFromURL(uri) // Fetch and append HTML content from each URL
	}

	finalList := extractPDFUrls(getData) // Extract all PDF links from HTML content

	outputDir := "PDFs/" // Directory to store downloaded PDFs

	if !directoryExists(outputDir) { // Check if directory exists
		createDirectory(outputDir, 0o755) // Create directory with read-write-execute permissions
	}

	// Remove duplicates from a given slice.
	finalList = removeDuplicatesFromSlice(finalList)

	// Loop through all extracted PDF URLs
	for _, urls := range finalList {
		if isUrlValid(urls) { // Check if the final URL is valid
			downloadPDF(urls, outputDir) // Download the PDF
		}
	}
}

// Extracts filename from full path (e.g. "/dir/file.pdf" → "file.pdf")
func getFilename(path string) string {
	return filepath.Base(path) // Use Base function to get file name only
}

// Converts a raw URL into a sanitized PDF filename safe for filesystem
func urlToFilename(rawURL string) string {
	lower := strings.ToLower(rawURL) // Convert URL to lowercase
	lower = getFilename(lower)       // Extract filename from URL

	reNonAlnum := regexp.MustCompile(`[^a-z0-9]`)   // Regex to match non-alphanumeric characters
	safe := reNonAlnum.ReplaceAllString(lower, "_") // Replace non-alphanumeric with underscores

	safe = regexp.MustCompile(`_+`).ReplaceAllString(safe, "_") // Collapse multiple underscores into one
	safe = strings.Trim(safe, "_")                              // Trim leading and trailing underscores

	var invalidSubstrings = []string{
		"_pdf", // Substring to remove from filename
	}

	for _, invalidPre := range invalidSubstrings { // Remove unwanted substrings
		safe = removeSubstring(safe, invalidPre)
	}

	if getFileExtension(safe) != ".pdf" { // Ensure file ends with .pdf
		safe = safe + ".pdf"
	}

	return safe // Return sanitized filename
}

// Removes all instances of a specific substring from input string
func removeSubstring(input string, toRemove string) string {
	result := strings.ReplaceAll(input, toRemove, "") // Replace substring with empty string
	return result
}

// Gets the file extension from a given file path
func getFileExtension(path string) string {
	return filepath.Ext(path) // Extract and return file extension
}

// Checks if a file exists at the specified path
func fileExists(filename string) bool {
	info, err := os.Stat(filename) // Get file info
	if err != nil {                // If error occurs, file doesn't exist
		return false
	}
	return !info.IsDir() // Return true if path is a file (not a directory)
}

// Downloads a PDF from given URL and saves it in the specified directory
func downloadPDF(finalURL, outputDir string) bool {
	filename := strings.ToLower(urlToFilename(finalURL)) // Sanitize the filename
	filePath := filepath.Join(outputDir, filename)       // Construct full path for output file

	if fileExists(filePath) { // Skip if file already exists
		log.Printf("File already exists, skipping: %s", filePath)
		return false
	}

	client := &http.Client{Timeout: 15 * time.Minute} // Create HTTP client with timeout

	resp, err := client.Get(finalURL) // Send HTTP GET request
	if err != nil {
		log.Printf("Failed to download %s: %v", finalURL, err)
		return false
	}
	defer resp.Body.Close() // Ensure response body is closed

	if resp.StatusCode != http.StatusOK { // Check if response is 200 OK
		log.Printf("Download failed for %s: %s", finalURL, resp.Status)
		return false
	}

	contentType := resp.Header.Get("Content-Type")         // Get content type of response
	if !strings.Contains(contentType, "application/pdf") { // Check if it's a PDF
		log.Printf("Invalid content type for %s: %s (expected application/pdf)", finalURL, contentType)
		return false
	}

	var buf bytes.Buffer                     // Create a buffer to hold response data
	written, err := io.Copy(&buf, resp.Body) // Copy data into buffer
	if err != nil {
		log.Printf("Failed to read PDF data from %s: %v", finalURL, err)
		return false
	}
	if written == 0 { // Skip empty files
		log.Printf("Downloaded 0 bytes for %s; not creating file", finalURL)
		return false
	}

	out, err := os.Create(filePath) // Create output file
	if err != nil {
		log.Printf("Failed to create file for %s: %v", finalURL, err)
		return false
	}
	defer out.Close() // Ensure file is closed after writing

	if _, err := buf.WriteTo(out); err != nil { // Write buffer contents to file
		log.Printf("Failed to write PDF to file for %s: %v", finalURL, err)
		return false
	}

	log.Printf("Successfully downloaded %d bytes: %s → %s", written, finalURL, filePath) // Log success
	return true
}

// Checks whether a given directory exists
func directoryExists(path string) bool {
	directory, err := os.Stat(path) // Get info for the path
	if err != nil {
		return false // Return false if error occurs
	}
	return directory.IsDir() // Return true if it's a directory
}

// Creates a directory at given path with provided permissions
func createDirectory(path string, permission os.FileMode) {
	err := os.Mkdir(path, permission) // Attempt to create directory
	if err != nil {
		log.Println(err) // Log error if creation fails
	}
}

// Verifies whether a string is a valid URL format
func isUrlValid(uri string) bool {
	_, err := url.ParseRequestURI(uri) // Try parsing the URL
	return err == nil                  // Return true if valid
}

// Removes duplicate strings from a slice
func removeDuplicatesFromSlice(slice []string) []string {
	check := make(map[string]bool) // Map to track seen values
	var newReturnSlice []string    // Slice to store unique values
	for _, content := range slice {
		if !check[content] { // If not already seen
			check[content] = true                            // Mark as seen
			newReturnSlice = append(newReturnSlice, content) // Add to result
		}
	}
	return newReturnSlice
}

// extractPDFUrls takes a raw text input (possibly containing HTML),
// extracts URLs, and returns only those that include all required keywords.
func extractPDFUrls(rawText string) []string {
	// Step 1: Unescape HTML entities so "https&#x3A;&#x2F;&#x2F;" becomes "https://"
	cleanText := html.UnescapeString(rawText)

	// Step 2: Regex to capture anything that looks like an HTTP or HTTPS URL
	urlPattern := regexp.MustCompile(`https?://[^\s"'<>]+`)
	allURLs := urlPattern.FindAllString(cleanText, -1)

	// Step 3: Define the required keywords that must all appear in the URL
	requiredKeywords := []string{"weiman.com", "mwdownloads", "download"}

	// Step 4: Filter URLs to keep only those that contain all required keywords
	var matchingURLs []string
	for _, url := range allURLs {
		matchesAll := true
		for _, keyword := range requiredKeywords {
			if !strings.Contains(url, keyword) {
				matchesAll = false
				break
			}
		}
		if matchesAll {
			matchingURLs = append(matchingURLs, url)
		}
	}

	// Step 5: Return the slice of URLs that passed the filter
	return matchingURLs
}

// Performs HTTP GET request and returns response body as string
func getDataFromURL(uri string) string {
	log.Println("Scraping", uri)   // Log which URL is being scraped
	response, err := http.Get(uri) // Send GET request
	if err != nil {
		log.Println(err) // Log if request fails
	}

	body, err := io.ReadAll(response.Body) // Read the body of the response
	if err != nil {
		log.Println(err) // Log read error
	}

	err = response.Body.Close() // Close response body
	if err != nil {
		log.Println(err) // Log error during close
	}
	return string(body) // Return response body as string
}
