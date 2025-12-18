package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)


type Article struct {
	Source struct {
		Name string `json:"name"`
	} `json:"source"`

	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	PublishedAt string `json:"publishedAt"`
}

type Response struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}


func FetchNews(query, category string, limit int) (*Response, error) {
	baseURL := "https://newsapi.org/v2/"

	apiKey := os.Getenv("APIKey")
	if apiKey == "" {
		return nil, fmt.Errorf("APIKey not set in environment")
	}

	var endpoint string
	var params string

	if query != "" {
		endpoint = "everything?"
		params = fmt.Sprintf(
			"q=%s&sortBy=publishedAt&language=en&pageSize=%d",
			strings.ReplaceAll(query, " ", "+"),
			limit,
		)
	} else {
		endpoint = "top-headlines?"
		params = fmt.Sprintf(
			"country=us&category=%s&pageSize=%d",
			category,
			limit,
		)
	}

	url := baseURL + endpoint + params + "&apiKey=" + apiKey

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var newsResp Response
	if err := json.Unmarshal(body, &newsResp); err != nil {
		return nil, fmt.Errorf("JSON parse error: %v\n%s", err, string(body))
	}

	if newsResp.Status != "ok" {
		return nil, fmt.Errorf("API returned status: %s", newsResp.Status)
	}

	// Sort newest first
	sort.Slice(newsResp.Articles, func(i, j int) bool {
		ti, _ := time.Parse(time.RFC3339, newsResp.Articles[i].PublishedAt)
		tj, _ := time.Parse(time.RFC3339, newsResp.Articles[j].PublishedAt)
		return ti.After(tj)
	})

	return &newsResp, nil
}


func DisplayNews(articles []Article) {
	if len(articles) == 0 {
		fmt.Println("No news found.")
		return
	}

	color.Cyan("Found %d articles:\n", len(articles))

	for i, a := range articles {
		titleColor := color.New(color.FgYellow).SprintFunc()
		lower := strings.ToLower(a.Title)

		if strings.Contains(lower, "crypto") ||
			strings.Contains(lower, "bitcoin") ||
			strings.Contains(lower, "stock") {
			titleColor = color.New(color.FgGreen).SprintFunc()
		}

		fmt.Printf("%d. %s\n", i+1, titleColor(a.Title))
		color.White("   [%s] - %s", a.Source.Name, a.PublishedAt[:10])
		fmt.Printf("   %s\n", strings.TrimSpace(a.Description))
		color.Blue("   Link: %s\n\n", a.URL)
	}
}


func main() {
	// Load .env once
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	app := &cli.App{
		Name:  "finnews",
		Usage: "CLI for financial news (stocks / crypto / markets)",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "q",
				Usage: "Search query (e.g. 'bitcoin', 'AAPL stock')",
			},
			&cli.StringFlag{
				Name:  "category",
				Value: "business",
				Usage: "Top headlines category",
			},
			&cli.IntFlag{
				Name:  "limit",
				Value: 10,
				Usage: "Number of articles (max 100)",
			},
		},
		Action: func(c *cli.Context) error {
			limit := c.Int("limit")
			if limit > 100 {
				limit = 100
			}

			news, err := FetchNews(
				c.String("q"),
				c.String("category"),
				limit,
			)
			if err != nil {
				return err
			}

			DisplayNews(news.Articles)
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
