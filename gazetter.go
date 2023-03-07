package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {

	url := flag.String("url", "https://egazette.nic.in/WriteReadData", "Base URL to make requests to (defaults to https://egazette.nic.in/WriteReadData/)")
	year := flag.Int("year", 2016, "Year the articles were published (defaults to 2016)")
	start := flag.Int("start", 160000, "Starting range of article ID (defaults to 160000)")
	end := flag.Int("end", 180000, "Ending range of article ID (defaults to 180000)")

	flag.Parse()

	for i := *start; i <= *end; i++ {
		requestUrl := fmt.Sprintf("%s/%d/%d.pdf", *url, *year, i)
		filename := fmt.Sprintf("%d.pdf", i)
		err := downloadFile(filename, requestUrl)
		if err != nil {
			fmt.Printf("error: %s\n", requestUrl)
		}
	}

}

func downloadFile(filepath string, url string) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		out, err := os.Create(filepath)
		if err != nil {
			return err
		}
		defer out.Close()
		_, err = io.Copy(out, resp.Body)
		return err
	} else if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("file not found: %s\n", url)
		return nil
	} else {
		return fmt.Errorf("error, got status code: %d for url: %s", resp.StatusCode, url)
	}
}
