package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (s *server) StreamStart() error {
	url := "https://download.samplelib.com/mp4/sample-5s.mp4"
	data, err := DownloadBytes(url)

	if err != nil {
		fmt.Println("error downloading:", err)
		return err
	}

	log.Println("length of data:", len(data))
	key := "video-key"
	streamer := &MockVideoStreamer{
		Store: map[string][]byte{key: data},
	}

	http.HandleFunc("/video", VideoStreamHandler(streamer, key, len(data)))
	return http.ListenAndServe(":8080", nil)
}

type VideoStreamer interface {
	Seek(key string, start, end int) ([]byte, error)
}

type MockVideoStreamer struct {
	Store map[string][]byte
}

func (m *MockVideoStreamer) Seek(key string, start, end int) ([]byte, error) {
	videoData, exists := m.Store[key]

	if !exists {
		return nil, fmt.Errorf("video not found")
	}

	videoLen := len(videoData)

	if start < 0 || start >= videoLen {
		return nil, fmt.Errorf("start byte %d out of range", start)
	}

	if end < start || end >= videoLen {
		end = videoLen - 1
	}

	return videoData[start : end+1], nil
}

func VideoStreamHandler(streamer VideoStreamer, videoKey string, totalSize int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rangeHeader := r.Header.Get("Range")

		var start, end int

		if rangeHeader == "" {
			start = 0
			end = 1024*1024 - 1
		} else {
			rangeParts := strings.TrimPrefix(rangeHeader, "bytes=")
			rangeValues := strings.Split(rangeParts, "-")

			var err error
			start, err = strconv.Atoi(rangeValues[0])
			if err != nil {
				http.Error(w, "Invalid start byte", http.StatusBadRequest)
				return
			}

			if len(rangeValues) > 1 && rangeValues[1] != "" {
				end, err = strconv.Atoi(rangeValues[1])
				if err != nil {
					http.Error(w, "Invalid end byte", http.StatusBadRequest)
					return
				}
			} else {
				end = start + 1024*1024 - 1
			}
		}

		if end >= totalSize {
			end = totalSize - 1
		}

		videoData, err := streamer.Seek(videoKey, start, end)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error retrieving vide: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, totalSize))
		w.Header().Set("Accept-Ranges", "bytes")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(videoData)))
		w.Header().Set("Content-Type", "video/mp4")
		w.WriteHeader(http.StatusPartialContent)

		_, err = w.Write(videoData)
		if err != nil {
			http.Error(w, "error streaming video", http.StatusInternalServerError)
		}
	}
}

func DownloadBytes(url string) ([]byte, error) {
	// Send a GET request to the URL
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
