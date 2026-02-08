package scanner

import (
	"bufio"
	"io"
	"log/slog"
	"net/http"
	"s3bypass/pkg/utils"
	"strings"
)

// ResponseData HTTP response'dan elde edilen verileri tutar
type ResponseData struct {
	StatusCode    int
	ContentLength int64
	WordCount     int
	LineCount     int
}

// RequestStrategy HTTP request stratejisi interface'i
type RequestStrategy interface {
	Execute(client *http.Client, url string) (*ResponseData, error)
}

// HeadRequestStrategy HEAD request kullanır (hızlı, body yok)
type HeadRequestStrategy struct{}

// NewHeadRequestStrategy yeni bir HeadRequestStrategy oluşturur
func NewHeadRequestStrategy() *HeadRequestStrategy {
	return &HeadRequestStrategy{}
}

// Execute HEAD request yapar
func (s *HeadRequestStrategy) Execute(client *http.Client, url string) (*ResponseData, error) {
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("User-Agent", utils.GetRandomUserAgent())
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	return &ResponseData{
		StatusCode:    resp.StatusCode,
		ContentLength: resp.ContentLength,
		WordCount:     0, // HEAD request body döndürmez
		LineCount:     0,
	}, nil
}

// GetRequestStrategy GET request kullanır (yavaş ama word/line sayar)
type GetRequestStrategy struct{}

// NewGetRequestStrategy yeni bir GetRequestStrategy oluşturur
func NewGetRequestStrategy() *GetRequestStrategy {
	return &GetRequestStrategy{}
}

// Execute GET request yapar ve body'yi analiz eder
func (s *GetRequestStrategy) Execute(client *http.Client, url string) (*ResponseData, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("User-Agent", utils.GetRandomUserAgent())
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	wordCount, lineCount := s.analyzeBody(resp.Body)
	
	return &ResponseData{
		StatusCode:    resp.StatusCode,
		ContentLength: resp.ContentLength,
		WordCount:     wordCount,
		LineCount:     lineCount,
	}, nil
}

// analyzeBody response body'yi okur ve word/line sayar
func (s *GetRequestStrategy) analyzeBody(body io.Reader) (wordCount int, lineCount int) {
	const maxBodySize = 1024 * 1024 // 1MB limit
	
	limitedReader := io.LimitReader(body, maxBodySize)
	scanner := bufio.NewScanner(limitedReader)
	
	for scanner.Scan() {
		lineCount++
		line := scanner.Text()
		words := strings.Fields(line)
		wordCount += len(words)
	}
	
	if err := scanner.Err(); err != nil {
		slog.Debug("Body okuma hatası", "error", err)
	}
	
	return wordCount, lineCount
}

// CreateRequestStrategy config'e göre uygun stratejiyi oluşturur
func CreateRequestStrategy(fullCheck bool) RequestStrategy {
	if fullCheck {
		return NewGetRequestStrategy()
	}
	return NewHeadRequestStrategy()
}
