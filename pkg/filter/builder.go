package filter

import (
	"fmt"
	"log/slog"
	"s3bypass/pkg/utils"
	"strings"
)

// FilterBuilder fluent API ile filter engine oluşturur
// ve parse hatalarını toplar
type FilterBuilder struct {
	codes  map[int]struct{}
	sizes  map[int]struct{}
	words  map[int]struct{}
	lines  map[int]struct{}
	errors []string
}

// NewFilterBuilder yeni bir FilterBuilder oluşturur
func NewFilterBuilder() *FilterBuilder {
	return &FilterBuilder{
		codes:  make(map[int]struct{}),
		sizes:  make(map[int]struct{}),
		words:  make(map[int]struct{}),
		lines:  make(map[int]struct{}),
		errors: make([]string, 0),
	}
}

// WithCodes status code filtreleri ekler
func (b *FilterBuilder) WithCodes(input string) *FilterBuilder {
	if input == "" {
		return b
	}
	
	parsed, err := utils.ParseIntList(input)
	if err != nil {
		b.errors = append(b.errors, fmt.Sprintf("Status code parse hatası (-fc '%s'): %v", input, err))
		return b
	}
	
	b.codes = parsed
	return b
}

// WithSizes response size filtreleri ekler
func (b *FilterBuilder) WithSizes(input string) *FilterBuilder {
	if input == "" {
		return b
	}
	
	parsed, err := utils.ParseIntList(input)
	if err != nil {
		b.errors = append(b.errors, fmt.Sprintf("Size parse hatası (-fs '%s'): %v", input, err))
		return b
	}
	
	b.sizes = parsed
	return b
}

// WithWords word count filtreleri ekler
func (b *FilterBuilder) WithWords(input string) *FilterBuilder {
	if input == "" {
		return b
	}
	
	parsed, err := utils.ParseIntList(input)
	if err != nil {
		b.errors = append(b.errors, fmt.Sprintf("Word count parse hatası (-fw '%s'): %v", input, err))
		return b
	}
	
	b.words = parsed
	return b
}

// WithLines line count filtreleri ekler
func (b *FilterBuilder) WithLines(input string) *FilterBuilder {
	if input == "" {
		return b
	}
	
	parsed, err := utils.ParseIntList(input)
	if err != nil {
		b.errors = append(b.errors, fmt.Sprintf("Line count parse hatası (-fl '%s'): %v", input, err))
		return b
	}
	
	b.lines = parsed
	return b
}

// Build filter engine'i oluşturur ve hataları raporlar
func (b *FilterBuilder) Build() *Engine {
	// Hatalar varsa logla
	if len(b.errors) > 0 {
		slog.Warn("Filter parse hataları tespit edildi", 
			"count", len(b.errors),
			"details", strings.Join(b.errors, "; "))
	}
	
	return &Engine{
		Codes: b.codes,
		Sizes: b.sizes,
		Words: b.words,
		Lines: b.lines,
	}
}

// HasErrors parse hatası olup olmadığını kontrol eder
func (b *FilterBuilder) HasErrors() bool {
	return len(b.errors) > 0
}

// Errors tüm parse hatalarını döndürür
func (b *FilterBuilder) Errors() []string {
	return b.errors
}
