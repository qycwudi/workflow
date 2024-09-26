package utils

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type HTTPLogWriter struct {
	Url        string
	User       string
	Pass       string
	Buffer     *bytes.Buffer
	Interval   time.Duration
	EntryCount int
	mu         sync.Mutex
	done       chan struct{}
}

func NewHTTPLogWriter(url, user, pass string, interval time.Duration) *HTTPLogWriter {
	writer := &HTTPLogWriter{
		Url:        url,
		User:       user,
		Pass:       pass,
		Buffer:     new(bytes.Buffer),
		Interval:   interval,
		EntryCount: 0,
		done:       make(chan struct{}),
	}
	go writer.startTimer()
	return writer
}

func (w *HTTPLogWriter) startTimer() {
	ticker := time.NewTicker(w.Interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			w.Flush()
		case <-w.done:
			return
		}
	}
}

func (w *HTTPLogWriter) Write(p []byte) (n int, err error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.Buffer.Len() > 0 {
		w.Buffer.WriteString(",")
	}

	n, err = w.Buffer.Write(p)
	if err != nil {
		return 0, err
	}

	w.EntryCount++
	if w.EntryCount >= 1000 {
		w.Flush()
	}

	return n, nil
}

func (w *HTTPLogWriter) sendLogs() {
	if w.Buffer.Len() == 0 {
		return
	}

	logs := fmt.Sprintf("[%s]", w.Buffer.String())
	req, err := http.NewRequest("POST", w.Url, bytes.NewBufferString(logs))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(w.User, w.Pass)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send logs: %v", err)
		return
	}
	defer resp.Body.Close()

	w.Buffer.Reset()
	w.EntryCount = 0 // Reset entry count after sending
}

func (w *HTTPLogWriter) Flush() {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.sendLogs()
}

func (w *HTTPLogWriter) Stop() {
	close(w.done)
}
