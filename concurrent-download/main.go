package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"example.com/downloader/downloader"
)

const fileURL = "https://sampletestfile.com/wp-content/uploads/2023/07/1.5-MB-PDF.pdf"

func main() {
	url := flag.String("url", fileURL, "URL of the file to check size")
	output := flag.String("output", "file.data", "Output file to save the size")
	chunks := flag.Int64("chunks", 4, "Number of chunks in parallel")

	// Validate required flags
	if *url == "" {
		flag.Usage()
		log.Fatal("Missing required flag: --url")
	}

	// Parse command line flags
	flag.Parse()

	download := downloader.NewDownloadFile(*url, *output)
	defer download.Close()

	fmt.Println(download.Summary())
	data, err := json.MarshalIndent(download, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(data))

	if err := download.StartDownload(*chunks); err != nil {
		log.Fatalf("Error starting download: %v", err)
	}
}
