package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

var wg sync.WaitGroup

type Gazette struct {
	url      string
	filepath string
}

func generateUrls(start int, end int, url string, year int, dir string, gazettes chan<- *Gazette) {
	for i := start; i <= end; i++ {
		gazettes <- &Gazette{
			url:      fmt.Sprintf("%s/%d/%d.pdf", url, year, i),
			filepath: filepath.Join(dir, fmt.Sprintf("%d.pdf", i)),
		}
	}
	close(gazettes)
}

func main() {
	gazettes := make(chan *Gazette)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Unable to get home directory")
	}
	url := flag.String("url", "https://egazette.nic.in/WriteReadData", "Base URL to make requests to (defaults to https://egazette.nic.in/WriteReadData/)")
	year := flag.Int("year", 2016, "Year the articles were published (defaults to 2016)")
	start := flag.Int("start", 160000, "Starting range of article ID (defaults to 160000)")
	end := flag.Int("end", 180000, "Ending range of article ID (defaults to 180000)")
	dir := flag.String("dir", homeDir, "Directory where the PDFs will be saved (defaults to your home directory)")
	flag.Parse()

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go processUrls(gazettes)
	}

	go generateUrls(*start, *end, *url, *year, *dir, gazettes)

	wg.Wait()

}

func processUrls(gazettes <-chan *Gazette) {
	defer wg.Done()
	for gazette := range gazettes {
		downloadFile(gazette.url, gazette.filepath)
	}
}

func downloadFile(url string, filepath string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("error: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		out, err := os.Create(filepath)
		if err != nil {
			fmt.Printf("error: %s\n", err)
		}
		defer out.Close()
		fmt.Printf("download %s to %s\n", url, filepath)
		io.Copy(out, resp.Body)
	} else if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("file not found: %s\n", url)
	} else {
		fmt.Printf("error, got status code: %d for url: %s\n", resp.StatusCode, url)
	}
}
