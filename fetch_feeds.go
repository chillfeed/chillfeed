package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
	"gopkg.in/yaml.v2"
)

type Feed struct {
	URL   string `yaml:"url"`
	Title string `yaml:"title,omitempty"`
}

type Config struct {
	Feeds []Feed `yaml:"feeds"`
}

type Article struct {
	Title      string    `json:"title"`
	Summary    string    `json:"summary"`
	Link       string    `json:"link"`
	Published  time.Time `json:"published"`
	FeedTitle  string    `json:"feedTitle"`
	FeedAuthor string    `json:"feedAuthor"`
}

// Helper function to strip HTML tags and limit to a few sentences
func limitSummary(input string, sentenceLimit int) string {
	// Strip HTML tags
	doc, err := html.Parse(strings.NewReader(input))
	if err != nil {
		return ""
	}
	var textBuilder strings.Builder
	var extractText func(*html.Node)
	extractText = func(n *html.Node) {
		if n.Type == html.TextNode {
			textBuilder.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractText(c)
		}
	}
	extractText(doc)
	text := textBuilder.String()

	// Split into sentences (simple implementation)
	sentences := strings.FieldsFunc(text, func(r rune) bool {
		return r == '.' || r == '!' || r == '?'
	})

	// Limit to specified number of sentences
	if len(sentences) > sentenceLimit {
		sentences = sentences[:sentenceLimit]
	}

	// Join sentences and trim spaces
	summary := strings.Join(sentences, ". ")
	summary = strings.TrimSpace(summary)

	// Add ellipsis if we've truncated the summary
	if len(sentences) < len(strings.FieldsFunc(text, func(r rune) bool {
		return r == '.' || r == '!' || r == '?'
	})) {
		summary += "..."
	}

	return summary
}

func main() {
	// Read and parse the YAML file
	yamlFile, err := os.ReadFile("feeds.yaml")
	if err != nil {
		fmt.Printf("Error reading YAML file: %v\n", err)
		return
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		fmt.Printf("Error parsing YAML: %v\n", err)
		return
	}

	var articles []Article
	parser := gofeed.NewParser()
	oneMonthAgo := time.Now().AddDate(0, -1, 0)

	for _, feed := range config.Feeds {
		parsedFeed, err := parser.ParseURL(feed.URL)
		if err != nil {
			fmt.Printf("Error parsing feed %s: %v\n", feed.URL, err)
			continue
		}

		feedTitle := parsedFeed.Title
		if feed.Title != "" {
			feedTitle = feed.Title
		}
		feedAuthor := ""
		if parsedFeed.Author != nil {
			feedAuthor = parsedFeed.Author.Name
		}

		for _, item := range parsedFeed.Items {
			if item.PublishedParsed == nil {
				continue // Skip items without a valid publication date
			}

			if item.PublishedParsed.Before(oneMonthAgo) {
				continue // Skip items older than two months
			}

			summary := item.Description
			if summary == "" {
				summary = item.Content
			}
			limitedSummary := limitSummary(summary, 3) // Limit to 3 sentences

			articles = append(articles, Article{
				Title:      item.Title,
				Summary:    limitedSummary,
				Link:       item.Link,
				Published:  *item.PublishedParsed,
				FeedTitle:  feedTitle,
				FeedAuthor: feedAuthor,
			})
		}
	}

	// Sort articles by published date, newest first
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Published.After(articles[j].Published)
	})

	file, err := os.Create("web/articles.json")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(articles)
	if err != nil {
		fmt.Printf("Error encoding JSON: %v\n", err)
		return
	}

	fmt.Println("Articles fetched, sorted, and saved successfully.")
}
