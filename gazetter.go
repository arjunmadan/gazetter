package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
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

func generateTnExtraordinaryUrls(url string, year int, dir string, gazettes chan<- *Gazette) {
	// http://www.stationeryprinting.tn.gov.in/extraordinary/extraord_list2022.php
	response, err := http.Get(fmt.Sprintf("%s%d.php", url, year))
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	document.Find("a").Each(func(_ int, s *goquery.Selection) {
		gazId, _ := s.Attr("href")
		if strings.Contains(gazId, ".pdf") {
			gazettes <- &Gazette{
				url:      fmt.Sprintf("http://www.stationeryprinting.tn.gov.in/extraordinary/%s", gazId),
				filepath: filepath.Join(dir, strings.Split(gazId, "/")[1]),
			}
		}
	})
	close(gazettes)
}

func generateTnWeeklyUrls(url string, year int, dir string, gazettes chan<- *Gazette) {
	// http://www.stationeryprinting.tn.gov.in/gazette/gazette_list2022.php
	response, err := http.Get(fmt.Sprintf("%s%d.php", url, year))
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	document.Find("a").Each(func(_ int, s *goquery.Selection) {
		url, _ := s.Attr("href")
		if strings.Contains(url, "issueno") {
			//http://www.stationeryprinting.tn.gov.in/gazette/gazette_det2022.php?issueno=52
			fmt.Println(s.Attr("href"))
			response, err := http.Get(fmt.Sprintf("http://www.stationeryprinting.tn.gov.in/gazette/%s", url))
			if err != nil {
				log.Fatal(err)
			}
			defer response.Body.Close()
			document, err := goquery.NewDocumentFromReader(response.Body)
			if err != nil {
				log.Fatal("Error loading HTTP response body. ", err)
			}

			document.Find("a").Each(func(_ int, s *goquery.Selection) {
				gazId, _ := s.Attr("href")
				if strings.Contains(gazId, ".pdf") {
					gazettes <- &Gazette{
						url:      fmt.Sprintf("http://www.stationeryprinting.tn.gov.in/gazette/%s", gazId),
						filepath: filepath.Join(dir, strings.Split(gazId, "/")[1]),
					}
				}
			})
		}
	})

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
	gaz := flag.String("gaz", "", "The gazette from which articles need to be downloaded")

	flag.Parse()

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go processUrls(gazettes)
	}

	if *gaz == "" {
		go generateUrls(*start, *end, *url, *year, *dir, gazettes)
	} else if *gaz == "TN" {
		go generateTnWeeklyUrls("http://www.stationeryprinting.tn.gov.in/gazette/gazette_list", *year, *dir, gazettes)
		go generateTnExtraordinaryUrls("http://www.stationeryprinting.tn.gov.in/extraordinary/extraord_list", *year, *dir, gazettes)
	}

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
