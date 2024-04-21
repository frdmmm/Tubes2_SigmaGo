package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/PuerkitoBio/goquery"
)

// Page represents a Wikipedia page with its title and links
type Page struct {
	Title string
	Links []string
}

var (
	visited  = make(map[string]bool)
	edgeTo   = make(map[string]string)
	mutex    = &sync.Mutex{}
	maxDepth = 6 // Maximum depth for BFS
)

func main() {
	r := gin.Default()

	r.POST("/solve", solveHandler)

	r.Run(":8080")
}

func solveHandler(c *gin.Context) {
	startTitle := c.PostForm("start")
	endTitle := c.PostForm("end")

	fmt.Println("Received start title:", startTitle)
	fmt.Println("Received end title:", endTitle)

	// Rest of the code...

	if startTitle == "" || endTitle == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both start and end titles are required"})
		return
	}

	fmt.Printf("Start Title: %s, End Title: %s\n", startTitle, endTitle)

	solution := solveWikiRace(startTitle, endTitle)

	c.JSON(http.StatusOK, gin.H{
		"solution": solution,
	})
}

func solveWikiRace(start, end string) []string {
	fmt.Println("Solving Wiki Race...")
	// Reset global variables
	visited = make(map[string]bool)
	edgeTo = make(map[string]string)

	if start == end {
		return []string{start}
	}

	queue := []string{start}
	visited[start] = true
	found := false

	for len(queue) > 0 {
		currTitle := queue[0]
		queue = queue[1:]

		if currTitle == end {
			found = true
			break
		}

		fmt.Printf("Scraping links for %s...\n", currTitle)
		links, err := scrapeLinks(currTitle)
		if err != nil {
			log.Printf("Error scraping links for %s: %v", currTitle, err)
			continue
		}

		fmt.Printf("Found %d links for %s\n", len(links), currTitle)

		for _, link := range links {
			if !visited[link] {
				visited[link] = true
				edgeTo[link] = currTitle
				queue = append(queue, link)
			}
		}
	}

	if !found {
		return []string{"No solution found"}
	}

	path := getPath(end)
	return path
}

func scrapeLinks(title string) ([]string, error) {
	url := fmt.Sprintf("https://en.wikipedia.org/wiki/%s", title)
	fmt.Println("Scraping URL:", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var links []string
	doc.Find("#bodyContent #mw-content-text a").Each(func(i int, s *goquery.Selection) {
		link, exists := s.Attr("href")
		if exists && strings.HasPrefix(link, "/wiki/") {
			linkTitle := strings.TrimPrefix(link, "/wiki/")
			links = append(links, linkTitle)
		}
	})

	return links, nil
}

func getPath(end string) []string {
	var path []string
	for end != "" {
		path = append([]string{end}, path...)
		end = edgeTo[end]
	}
	return path
}
