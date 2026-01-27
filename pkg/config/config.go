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
	Delay       int
	Verbose     bool
	Wordlist    string
	FilterCode  string
	FilterSize  string
	FilterWord  string
	FilterLine  string
}

// Load parses CLI flags and returns a Config struct
func Load() (*Config, error) {
	list := flag.String("l", "", "Input file containing bucket URLs or names")
	single := flag.String("u", "", "Single bucket URL to scan")
	output := flag.String("o", DefaultOutputFile, "Output file for found secrets")
	threads := flag.Int("t", DefaultThreadCount, "Number of concurrent threads")
	timeout := flag.Int("to", DefaultTimeoutSec, "Request timeout in seconds")
	delay := flag.Int("delay", DefaultDelayMs, "Delay between requests in milliseconds")
	verbose := flag.Bool("v", false, "Enable verbose logging (debug mode)")
	wordlist := flag.String("w", "", "Path to custom wordlist file (optional)")
	fc := flag.String("fc", "", "Filter HTTP status codes (e.g. 404,403)")
	fs := flag.String("fs", "", "Filter HTTP response sizes (e.g. 0,1024)")
	fw := flag.String("fw", "", "Filter by amount of words (e.g. 0)")
	fl := flag.String("fl", "", "Filter by amount of lines (e.g. 0)")

	flag.Parse()

	cfg := &Config{
		ListFile:    *list,
		SingleURL:   *single,
		OutputFile:  *output,
		ThreadCount: *threads,
		Timeout:     *timeout,
		Delay:       *delay,
		Verbose:     *verbose,
		Wordlist:    *wordlist,
		FilterCode:  *fc,
		FilterSize:  *fs,
		FilterWord:  *fw,
		FilterLine:  *fl,
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
