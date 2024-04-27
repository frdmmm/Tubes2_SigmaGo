package main

import (
	"fmt"
	"log"
	"net/http"
  "context"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

// Page represents a Wikipedia page with its title and links
type Page struct {
  Title string
  Links []string
}

var (
  visited      map[string]bool
  edgeTo       map[string]string
  maxDepth     int = 6
  searchAlgo    string
)

func main() {
  r := gin.Default()
  r.Use(func(c *gin.Context) {
    c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:7000")
    c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")
    if c.Request.Method == "OPTIONS" {
        c.AbortWithStatus(http.StatusNoContent)
        return
    }
    c.Next()
})

  r.POST("/solve", solveHandler)

  r.Run(":8080")
}

func solveHandler(c *gin.Context) {
  startTitle := c.PostForm("start")
  endTitle := c.PostForm("end")
  searchAlgo = c.PostForm("algo") // Get the chosen search algorithm

  fmt.Println("Received start title:", startTitle)
  fmt.Println("Received end title:", endTitle)
  fmt.Println("Search algorithm:", searchAlgo)

  if startTitle == "" || endTitle == "" || searchAlgo != "bfs" && searchAlgo != "ids" {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Both start and end titles are required. Valid search algorithms are bfs and ids"})
    return
  }

  startTime := time.Now()
  visited = make(map[string]bool)
  edgeTo = make(map[string]string)

  var solution []string
  var articlesChecked int
  var pathLength int
  if searchAlgo == "bfs" {
    solution, articlesChecked, pathLength = solveBFS(startTitle, endTitle)
  } else {
    solution, articlesChecked, pathLength = solveIDS(startTitle, endTitle)
  }
  totalTime := time.Since(startTime)

  c.JSON(http.StatusOK, gin.H{
    "solution":       solution,
    "articlesChecked": articlesChecked,
    "pathLength":      pathLength, // Only applicable for IDS
    "timeTaken":       totalTime.String(),
  })
}

func solveBFS(start, end string) ([]string, int, int) {
  fmt.Println("Solving Wiki Race using BFS...")

  var mutex sync.Mutex

  if start == end {
      return []string{start}, 0, 0
  }

  queue := make(chan string, 10000)
  visited := make(map[string]bool)
  articlesChecked := 1
  visited[start] = true
  queue <- start

  ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
  defer cancel()

  g, _ := errgroup.WithContext(ctx)
  g.SetLimit(500)


  defer func() {

      go func() {
          g.Wait()
          close(queue)
      }()
  }()

  for {
      select {
      case <-ctx.Done():
          return []string{"No solution found"}, articlesChecked, 0
      case currTitle, ok := <-queue:
          if !ok {
              continue
          }

          if currTitle == end {
              return getPath(end), articlesChecked, len(getPath(end))
          }

          fmt.Printf("Scraping links for %s...\n", currTitle)

          curr := currTitle
          articlesChecked++
          g.Go(func() error {
              links, err := scrapeLinks(curr)
              if err != nil {
                  log.Printf("Error scraping links for %s: %v", curr, err)
                  return err
              }

              mutex.Lock()
              defer mutex.Unlock()

              for _, link := range links {
                  if !visited[link] {
                      visited[link] = true
                      edgeTo[link] = curr
                      select {
                      case queue <- link:
                      default:

                      }
                  }
              }

              return nil
          })
      }
  }
}





func solveIDS(start, end string) ([]string, int, int) {
  fmt.Println("Solving Wiki Race using IDS...")

  for depth := 1; depth <= maxDepth; depth++ {
    visited = make(map[string]bool)
    edgeTo = make(map[string]string)
    solution, articlesChecked := solveLimitedDepthDFS(start, end, depth)
    if solution != nil {
      return solution, articlesChecked, len(solution) // Path length for IDS is the solution length
    }
  }
  return []string{"No solution found within max depth"}, 0, 0
}

func solveLimitedDepthDFS(start, end string, depth int) ([]string, int) {
  visited[start] = true
  articlesChecked := 1

  if depth == 0 || start == end {
    if start == end {
      return []string{start}, articlesChecked
    }
    return nil, articlesChecked // No solution found within current depth
  }

  fmt.Printf("Scraping links for %s...\n", start)
  articlesChecked++;
  links, err := scrapeLinks(start)
  if err != nil {
    log.Printf("Error scraping links for %s: %v", start, err)
    return nil, articlesChecked
  }


  for _, link := range links {
    solution, articlesChecked := solveLimitedDepthDFS(link, end, depth-1)
    if solution != nil {
      edgeTo[link] = start
      return append([]string{link}, solution...), articlesChecked
    }
  }

  return nil, articlesChecked // No solution found within current depth at this branch
}

func scrapeLinks(title string) ([]string, error) {
    url := fmt.Sprintf("https://en.wikipedia.org/wiki/%s", title)
    // fmt.Println("Scraping URL:", url)
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
        if exists && strings.HasPrefix(link, "/wiki/") && !strings.HasSuffix(link, getMediaSuffix()) {
            linkTitle := strings.TrimPrefix(link, "/wiki/")
            if !strings.Contains(linkTitle, ":") {
                // Lock the mutex before accessing the visited map
                // Unlock the mutex when done

                // Only add the link if it hasn't been visited before
                links = append(links, linkTitle)
            }
        }
    })

    return links, nil
}
func getMediaSuffix() string {
    return `\.jpg|\.png|\.gif|\.bmp|\.mov|\.avi|\.mp4|\.pdf|\.docx|\.xlsx|\.pptx|\.jpeg`
}


func getPath(end string) []string {
    var path []string
    for end != "" {
      path = append([]string{end}, path...)
      end = edgeTo[end]
    }
    return path
  }
  