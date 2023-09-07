package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

func crawlPage(url string) []string {
	response, err := http.Get(url)
	if err != nil {
		log.Printf("Error fetching URL %s: %s", url, err)
		return nil
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		log.Printf("Failed to fetch URL %s: Status code %d", url, response.StatusCode)
		return nil
	}

	var links []string

	tokenizer := html.NewTokenizer(response.Body)
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return links
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()

			if token.Data == "a" {
				for _, attr := range token.Attr {
					if attr.Key == "href" {
						link := attr.Val
						if strings.HasPrefix(link, "http") {
							fmt.Println("Found link:", link)
							links = append(links, link)
						}
					}
				}
			}
		}
	}
}

func saveLinksToJSON(links []string) {
	jsonData, err := json.MarshalIndent(links, "", " ")
	if err != nil {
		log.Println("Error marshaling data:", err)
		return
	}

	err = ioutil.WriteFile("output.json", jsonData, 0644)
	if err != nil {
		log.Println("Error writing to file:", err)
		return
	}

	log.Println("Links saved to output.json")
}

func main() {
	urlPtr := flag.String("u", "", "URL to crawl")
	flag.Parse()

	if *urlPtr == "" {
		log.Println("please provide a URL using -u flag")
		return
	}

	links := crawlPage(*urlPtr)
	if len(links) > 0 {
		saveLinksToJSON(links)
	}
}
