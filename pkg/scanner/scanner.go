package scanner

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"s3bypass/pkg/config"
	"s3bypass/pkg/utils"
	"sync"
	"time"
)

// Scanner handles the scanning operations
type Scanner struct {
	cfg     *config.Config
	client  *http.Client
	buckets []string
	payloads []string
	
	// Filters
	filterCodes map[int]struct{}
	filterSizes map[int]struct{}
	filterWords map[int]struct{}
	filterLines map[int]struct{}
}

// New creates a new Scanner instance
func New(cfg *config.Config, buckets []string) *Scanner {
	tr := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:        cfg.ThreadCount,
		IdleConnTimeout:     config.DefaultIdleTimeout * time.Second,
		DisableCompression:  true,
		MaxIdleConnsPerHost: cfg.ThreadCount,
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Duration(cfg.Timeout) * time.Second,
	}

	// Load payloads
	var scanPayloads []string
	if cfg.Wordlist != "" {
		lines, err := utils.ReadLines(cfg.Wordlist)
		if err != nil {
			fmt.Printf("⚠️ Failed to load wordlist: %v. Using default payloads.\n", err)
			scanPayloads = Payloads
		} else {
			scanPayloads = lines
		}
	} else {
		scanPayloads = Payloads
	}

	// Parse Filters
	fCodes, _ := utils.ParseIntList(cfg.FilterCode)
	fSizes, _ := utils.ParseIntList(cfg.FilterSize)
	fWords, _ := utils.ParseIntList(cfg.FilterWord)
	fLines, _ := utils.ParseIntList(cfg.FilterLine)

	return &Scanner{
		cfg:      cfg,
		client:   client,
		buckets:  buckets,
		payloads: scanPayloads,
		filterCodes: fCodes,
		filterSizes: fSizes,
		filterWords: fWords,
		filterLines: fLines,
	}
}

// Start initiates the scanning process using a worker pool
func (s *Scanner) Start() {
	jobs := make(chan Job, s.cfg.ThreadCount*config.ChannelBufferMulti) // buffer for jobs
	results := make(chan Result, s.cfg.ThreadCount*config.ChannelBufferMulti)

	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < s.cfg.ThreadCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.worker(jobs, results)
		}()
	}

	// Result handler
	go s.handleResults(results)

	// Job generator
	s.generateJobs(jobs)
	close(jobs)

	wg.Wait()
	close(results)
}

func (s *Scanner) generateJobs(jobs chan<- Job) {
	for _, bucket := range s.buckets {
		for _, prefix := range Prefixes {
			for _, payload := range s.payloads {
				jobs <- Job{
					Bucket:  bucket,
					Prefix:  prefix,
					Payload: payload,
				}
			}
		}
	}
}

func (s *Scanner) handleResults(results <-chan Result) {
	outputFile, err := os.Create(s.cfg.OutputFile)
	if err != nil {
		fmt.Printf("❌ Failed to create output file: %v\n", err)
		return
	}
	defer outputFile.Close()
	outputFile.WriteString("--- S3 SCAN RESULTS ---\n")

	for result := range results {
		msg := fmt.Sprintf("✅ [FOUND] %s (Size: %d)", result.URL, result.Size)
		fmt.Println(msg)
		outputFile.WriteString(msg + "\n")
	}
}
