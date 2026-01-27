package config

import (
	"errors"
	"flag"
	"os"
)

// Config holds the application configuration
type Config struct {
	ListFile    string
	SingleURL   string
	OutputFile  string
	ThreadCount int
	Timeout     int
}

// Load parses CLI flags and returns a Config struct
func Load() (*Config, error) {
	list := flag.String("l", "", "Input file containing bucket URLs or names")
	single := flag.String("u", "", "Single bucket URL to scan")
	output := flag.String("o", DefaultOutputFile, "Output file for found secrets")
	threads := flag.Int("t", DefaultThreadCount, "Number of concurrent threads")
	timeout := flag.Int("to", DefaultTimeoutSec, "Request timeout in seconds")

	flag.Parse()

	cfg := &Config{
		ListFile:    *list,
		SingleURL:   *single,
		OutputFile:  *output,
		ThreadCount: *threads,
		Timeout:     *timeout,
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// validate checks the configuration values
func (c *Config) validate() error {
	if c.ListFile == "" && c.SingleURL == "" {
		return errors.New("either -l (list file) or -u (single URL) must be provided")
	}
	if c.ListFile != "" {
		if _, err := os.Stat(c.ListFile); os.IsNotExist(err) {
			return errors.New("input file does not exist: " + c.ListFile)
		}
	}
	if c.ThreadCount <= 0 {
		return errors.New("thread count must be positive")
	}
	return nil
}
