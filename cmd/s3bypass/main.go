package main

import (
	"fmt"
	"os"
	"regexp"
	"s3bypass/pkg/config"
	"s3bypass/pkg/scanner"
	"s3bypass/pkg/utils"
	"time"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}

	// 2. Read Inputs
	var lines []string
	if cfg.ListFile != "" {
		fileLines, err := utils.ReadLines(cfg.ListFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input file: %v\n", err)
			os.Exit(1)
		}
		lines = append(lines, fileLines...)
	}

	if cfg.SingleURL != "" {
		lines = append(lines, cfg.SingleURL)
	}

	// 3. Parse Bucket Names
	buckets := parseBuckets(lines)
	if len(buckets) == 0 {
		fmt.Println("⚠️ No valid buckets found in input.")
		os.Exit(0)
	}

	// 4. Initialize Scanner
	scan := scanner.New(cfg, buckets)

	// 5. Start Scanning
	fmt.Printf("🔥 Starting S3 Hunter v2.0\n")
	fmt.Printf("📦 Buckets: %d | 🧵 Threads: %d\n", len(buckets), cfg.ThreadCount)
	fmt.Println("--------------------------------------------------")
	
	start := time.Now()
	scan.Start()
	
	fmt.Printf("\n🏁 Scan completed in %s. Check '%s'.\n", time.Since(start), cfg.OutputFile)
}

func parseBuckets(lines []string) []string {
	// Regex to extract bucket name from standard S3 URLs
	reURL := regexp.MustCompile(config.S3UrlRegex)
	
	bucketMap := make(map[string]bool)
	var buckets []string

	for _, line := range lines {
		// Clean common prefixes if user copy-pasted output
		match := reURL.FindStringSubmatch(line)
		var name string
		
		if len(match) > 1 {
			name = match[1]
		} else {
			// Clean prefix
			cleaned := line
			if len(cleaned) > len(config.ProtectedPrefix) && cleaned[:len(config.ProtectedPrefix)] == config.ProtectedPrefix {
				cleaned = cleaned[len(config.ProtectedPrefix):]
			}
			
			// Basic validation
			valid := true
			for _, r := range cleaned {
				if r == ' ' || r == '\t' {
					valid = false
					break
				}
			}
			
			if valid && len(cleaned) > config.MinBucketNameLen {
				name = cleaned
			}
		}

		if name != "" && !bucketMap[name] {
			bucketMap[name] = true
			buckets = append(buckets, name)
		}
	}
	return buckets
}
