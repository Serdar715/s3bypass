package output

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
)

// Result sonuç verisini temsil eder
type Result struct {
	URL  string
	Size int64
}

// OutputWriter sonuçların yazılması için interface
type OutputWriter interface {
	WriteHeader() error
	Write(url string, size int64) error
	Close() error
}

// FileOutputWriter dosyaya yazan OutputWriter implementasyonu
type FileOutputWriter struct {
	file   *os.File
	writer *bufio.Writer
}

// NewFileOutputWriter yeni bir FileOutputWriter oluşturur
func NewFileOutputWriter(path string) (*FileOutputWriter, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("output dosyası oluşturulamadı: %w", err)
	}
	
	return &FileOutputWriter{
		file:   file,
		writer: bufio.NewWriter(file),
	}, nil
}

// WriteHeader dosyaya başlık yazar
func (w *FileOutputWriter) WriteHeader() error {
	_, err := w.writer.WriteString("--- S3 SCAN RESULTS ---\n")
	if err != nil {
		return fmt.Errorf("header yazma hatası: %w", err)
	}
	return nil
}

// Write tek bir sonucu yazar
func (w *FileOutputWriter) Write(url string, size int64) error {
	msg := fmt.Sprintf("✅ [FOUND] %s (Size: %d)\n", url, size)
	_, err := w.writer.WriteString(msg)
	if err != nil {
		return fmt.Errorf("sonuç yazma hatası: %w", err)
	}
	return nil
}

// Close buffer'ı flush eder ve dosyayı kapatır
func (w *FileOutputWriter) Close() error {
	if err := w.writer.Flush(); err != nil {
		slog.Error("Buffer flush hatası", "error", err)
	}
	return w.file.Close()
}

// ConsoleOutputWriter stdout'a yazan OutputWriter implementasyonu
type ConsoleOutputWriter struct{}

// NewConsoleOutputWriter yeni bir ConsoleOutputWriter oluşturur
func NewConsoleOutputWriter() *ConsoleOutputWriter {
	return &ConsoleOutputWriter{}
}

// WriteHeader stdout'a başlık yazar
func (w *ConsoleOutputWriter) WriteHeader() error {
	fmt.Println("--- S3 SCAN RESULTS ---")
	return nil
}

// Write stdout'a sonuç yazar
func (w *ConsoleOutputWriter) Write(url string, size int64) error {
	fmt.Printf("✅ [FOUND] %s (Size: %d)\n", url, size)
	return nil
}

// Close ConsoleOutputWriter için no-op
func (w *ConsoleOutputWriter) Close() error {
	return nil
}

// MultiOutputWriter birden fazla writer'a yazar
type MultiOutputWriter struct {
	writers []OutputWriter
}

// NewMultiOutputWriter yeni bir MultiOutputWriter oluşturur
func NewMultiOutputWriter(writers ...OutputWriter) *MultiOutputWriter {
	return &MultiOutputWriter{writers: writers}
}

// WriteHeader tüm writer'lara header yazar
func (w *MultiOutputWriter) WriteHeader() error {
	for _, writer := range w.writers {
		if err := writer.WriteHeader(); err != nil {
			slog.Error("Header yazma hatası", "error", err)
		}
	}
	return nil
}

// Write tüm writer'lara sonuç yazar
func (w *MultiOutputWriter) Write(url string, size int64) error {
	for _, writer := range w.writers {
		if err := writer.Write(url, size); err != nil {
			slog.Error("Sonuç yazma hatası", "error", err)
		}
	}
	return nil
}

// Close tüm writer'ları kapatır
func (w *MultiOutputWriter) Close() error {
	for _, writer := range w.writers {
		if err := writer.Close(); err != nil {
			slog.Error("Writer kapatma hatası", "error", err)
		}
	}
	return nil
}
