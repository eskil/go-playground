package downloader

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

type DownloadFile struct {
	URL          string   `json:"url"`
	Output       string   `json:"output"`
	Size         int64    `json:"size"`
	Preallocated bool     `json:"preallocated"`
	File         *os.File `json:"-"`
}

func NewDownloadFile(url, output string) DownloadFile {
	df := DownloadFile{
		URL:          url,
		Output:       output,
		Size:         0,
		Preallocated: false,
		File:         nil,
	}
	df.FetchSize()
	df.Preallocate()
	fmt.Println(df.Summary())
	return df
}

func (df *DownloadFile) Close() error {
	if df.File != nil {
		return df.File.Close()
	}
	return nil
}

// Method on the struct
func (df DownloadFile) Summary() string {
	return fmt.Sprintf("Downloading %s to %s (%d bytes)", df.URL, df.Output, df.Size)
}

func (df *DownloadFile) FetchSize() error {
	size, err := get_file_size(df.URL)
	if err != nil {
		return fmt.Errorf("error fetching file size: %w", err)
	}
	df.Size = size
	return nil
}

func (df *DownloadFile) Preallocate() error {
	if df.Size <= 0 {
		return fmt.Errorf("file size is not set or invalid")
	}

	var err error
	df.File, err = os.Create(df.Output)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}

	if err := df.File.Truncate(df.Size); err != nil {
		return fmt.Errorf("error truncating file to size %d: %w", df.Size, err)
	}
	df.Preallocated = true

	return nil
}

func (df *DownloadFile) StartDownload(chunks int64) error {
	var wg sync.WaitGroup
	ch := make(chan DownloadTaskResult)
	done_ch := make(chan struct{})

	pos := int64(0)
	chunksz := df.Size / chunks
	for i := int64(0); i < chunks; i++ {
		fmt.Printf("Starting downloader %d\n", i)
		// Ensure no off-by-one because of dividing the filesize.
		if i == chunks-1 && (i*chunksz+chunksz) < df.Size {
			chunksz = df.Size - pos
		}
		task := NewDownloadTask(df, pos, chunksz)
		fmt.Println(task.Summary())
		data, err := json.MarshalIndent(task, "", "  ")
		if err != nil {
			panic(err)
		}
		fmt.Println(string(data))
		wg.Add(1)
		go task.DownloadPart(ch, &wg)
		pos += chunksz
	}

	// Wait for all goroutines in wg to finish, then close the done_ch,
	// that way we can exit the loop gracefully by selecting on done_ch.
	go func() {
		wg.Wait()
		close(done_ch)
	}()

loop:
	for {
		select {
		case msg := <-ch:
			// Process message from a download task.
			if msg.Err != nil {
				log.Printf("Error downloading part: %v", msg.Err)
			} else {
				data, err := json.MarshalIndent(msg, "", "  ")
				if err != nil {
					panic(err)
				}
				fmt.Println(string(data))
				log.Printf("Downloaded part %d%% from %d to %d (%d bytes)",
					msg.Progress,
					msg.DownloadTask.From,
					msg.DownloadTask.From+msg.DownloadTask.Size,
					msg.DownloadTask.Size)
			}
		case <-done_ch:
			// This means the wg has finished and all download tasks are done.
			fmt.Println("All download tasks completed.")
			// Use a labeled break to exit the loop, not just the select
			break loop
		}
	}

	fmt.Println("Download completed.")

	return nil
}

func get_file_size(url string) (int64, error) {
	// Make a HEAD request to the URL to get the file size but not download the file
	resp, err := http.Head(url)
	if err != nil {
		return 0, fmt.Errorf("failed to make HEAD request: %w", err)
	}
	// Ensure the response body is closed after we're done with it
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("non-OK HTTP status: %s", resp.Status)
	}

	size := resp.ContentLength
	if size < 0 {
		return 0, fmt.Errorf("could not determine file size")
	}

	return size, nil
}
