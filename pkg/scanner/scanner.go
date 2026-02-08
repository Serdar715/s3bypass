package scanner

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"s3bypass/pkg/config"
	"s3bypass/pkg/filter"
	"s3bypass/pkg/limiter"
	"s3bypass/pkg/utils"
	"sync"
	"time"
)

// Scanner handles the scanning operations
type Scanner struct {
	cfg      *config.Config
	client   *http.Client
	buckets  []string
	payloads []string
	
	filter   *filter.Engine
	limiter  *limiter.RateLimiter
	strategy RequestStrategy
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

	return &Scanner{
		cfg:      cfg,
		client:   client,
		buckets:  buckets,
		payloads: scanPayloads,
		filter:   filter.New(cfg),
		limiter:  limiter.New(cfg.Delay),
		strategy: CreateRequestStrategy(cfg.FullCheck),
	}
}

// Start initiates the scanning process using a worker pool
func (s *Scanner) Start() {
	jobs := make(chan Job, s.cfg.ThreadCount*config.ChannelBufferMulti)
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

	// Result handler - yeni yapı ile
	resultHandler := NewResultHandler(s.cfg.OutputFile)
	resultHandler.Start(results)

	// Job generator
	s.generateJobs(jobs)
	close(jobs)

	wg.Wait()
	close(results)
	
	// ResultHandler'ın tüm sonuçları yazmasını bekle
	resultHandler.Wait()
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
