package filter

import (
	"net/http"
	"s3bypass/pkg/config"
)

// Engine handles filtering of scan results
type Engine struct {
	Codes map[int]struct{}
	Sizes map[int]struct{}
	Words map[int]struct{}
	Lines map[int]struct{}
}

// New creates a new Filter Engine by parsing config using FilterBuilder
func New(cfg *config.Config) *Engine {
	builder := NewFilterBuilder().
		WithCodes(cfg.FilterCode).
		WithSizes(cfg.FilterSize).
		WithWords(cfg.FilterWord).
		WithLines(cfg.FilterLine)
	
	return builder.Build()
}

// ShouldSkip returns true if the response matches any filter
func (e *Engine) ShouldSkip(resp *http.Response) bool {
	// Filter by Status Code
	if len(e.Codes) > 0 {
		if _, ok := e.Codes[resp.StatusCode]; ok {
			return true
		}
	}

	// Filter by Content Size
	if len(e.Sizes) > 0 {
		if _, ok := e.Sizes[int(resp.ContentLength)]; ok {
			return true
		}
	}

	// Filter by Words/Lines (Note: HEAD output is 0)
	// If user specifically filters 0, we skip.
	if len(e.Words) > 0 {
		if _, ok := e.Words[0]; ok {
			return true
		}
	}
	
	if len(e.Lines) > 0 {
		if _, ok := e.Lines[0]; ok {
			return true
		}
	}

	return false
}
