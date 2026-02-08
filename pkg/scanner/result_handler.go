package scanner

import (
	"bufio"
	"fmt"
	"log/slog"
	"os"
	"sync"
)

// ResultHandler yönetir sonuçların dosyaya yazılmasını
// ve goroutine senkronizasyonunu sağlar
type ResultHandler struct {
	outputPath string
	done       chan struct{}
	wg         sync.WaitGroup
}

// NewResultHandler yeni bir ResultHandler oluşturur
func NewResultHandler(outputPath string) *ResultHandler {
	return &ResultHandler{
		outputPath: outputPath,
		done:       make(chan struct{}),
	}
}

// Start sonuçları dinlemeye başlar ve dosyaya yazar
// Bu metod bir goroutine içinde çalıştırılmalıdır
func (h *ResultHandler) Start(results <-chan Result) {
	h.wg.Add(1)
	go func() {
		defer h.wg.Done()
		defer close(h.done)
		h.processResults(results)
	}()
}

// processResults sonuçları işler ve dosyaya yazar
func (h *ResultHandler) processResults(results <-chan Result) {
	outputFile, err := os.Create(h.outputPath)
	if err != nil {
		slog.Error("Output dosyası oluşturulamadı", "path", h.outputPath, "error", err)
		return
	}
	defer outputFile.Close()

	writer := bufio.NewWriter(outputFile)
	defer func() {
		if flushErr := writer.Flush(); flushErr != nil {
			slog.Error("Buffer flush hatası", "error", flushErr)
		}
	}()

	if err := h.writeHeader(writer); err != nil {
		slog.Error("Header yazma hatası", "error", err)
		return
	}

	for result := range results {
		if err := h.writeResult(writer, result); err != nil {
			slog.Error("Sonuç yazma hatası", "url", result.URL, "error", err)
			continue
		}
	}
}

// writeHeader dosyaya başlık yazar
func (h *ResultHandler) writeHeader(writer *bufio.Writer) error {
	_, err := writer.WriteString("--- S3 SCAN RESULTS ---\n")
	return err
}

// writeResult tek bir sonucu dosyaya ve stdout'a yazar
func (h *ResultHandler) writeResult(writer *bufio.Writer, result Result) error {
	msg := fmt.Sprintf("✅ [FOUND] %s (Size: %d)", result.URL, result.Size)
	fmt.Println(msg)
	
	_, err := writer.WriteString(msg + "\n")
	return err
}

// Wait handler'ın tüm sonuçları işlemesini bekler
func (h *ResultHandler) Wait() {
	h.wg.Wait()
}

// Done senkronizasyon için done channel'ını döndürür
func (h *ResultHandler) Done() <-chan struct{} {
	return h.done
}
