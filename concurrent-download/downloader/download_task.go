package downloader

import (
	"fmt"
	"io"
	"net/http"
	"sync"
)

type DownloadTask struct {
	DownloadFile *DownloadFile `json:"download_file"`
	From         int64         `json:"from"`
	Size         int64         `json:"size"`
}

type DownloadTaskResult struct {
	DownloadTask *DownloadTask `json:"download_task"`
	Progress     int64         `json:"progress"`
	Err          error         `json:"error,omitempty" `
}

func NewDownloadTask(df *DownloadFile, from, size int64) *DownloadTask {
	return &DownloadTask{
		DownloadFile: df,
		From:         from,
		Size:         size,
	}
}

func (dt *DownloadTask) Summary() string {
	return fmt.Sprintf("Task for %s from %d to %d (%d bytes)", dt.DownloadFile.URL, dt.From, dt.From+dt.Size, dt.Size)
}

func (dt *DownloadTask) downloadPart() ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", dt.DownloadFile.URL, nil)
	if err != nil {
		err = fmt.Errorf("failed to create request: %w", err)
		return nil, err
	}

	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", dt.From, dt.From+dt.Size-1))
	resp, err := client.Do(req)
	if err != nil {
		err = fmt.Errorf("failed to download part: %w", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusPartialContent {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}

func (dt *DownloadTask) DownloadPart(ch chan<- DownloadTaskResult, wg *sync.WaitGroup) {
	// Ensure that we signal the WaitGroup when this function completes
	defer wg.Done()

	ch <- DownloadTaskResult{
		DownloadTask: dt,
		Err:          nil,
		Progress:     0,
	}

	// TODO: store data in the file
	data, err := dt.downloadPart()
	dt.DownloadFile.File.WriteAt(data, dt.From)

	if err != nil {
		ch <- DownloadTaskResult{
			DownloadTask: dt,
			Err:          err,
			Progress:     0,
		}
	} else {
		ch <- DownloadTaskResult{
			DownloadTask: dt,
			Err:          nil,
			Progress:     100,
		}
	}
}
