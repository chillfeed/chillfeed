package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
	"gopkg.in/yaml.v2"
)

const defaultArticlesPerPage = 20
const defaultFetchWeeks = 4
const defaultRepo = "chillfeed/chillfeed"
const defaultConfigFile = "config.yaml.example"

type Feed struct {
	URL   string `yaml:"url"`
	Title string `yaml:"title,omitempty"`
}

type Config struct {
	Feeds           []Feed `yaml:"feeds"`
	ArticlesPerPage int    `yaml:"articlesPerPage,omitempty"`
	FetchWeeks      int    `yaml:"fetchWeeks,omitempty"`
	Repo            string `yaml:"repo,omitempty"`
}

type Article struct {
	Title      string    `json:"title"`
	Summary    string    `json:"summary"`
	Link       string    `json:"link"`
	Published  time.Time `json:"published"`
	FeedTitle  string    `json:"feedTitle"`
	FeedAuthor string    `json:"feedAuthor"`
	Homepage   string    `json:"homepage"`
}

type Metadata struct {
	TotalPages   int       `json:"totalPages"`
	LastFetched  time.Time `json:"lastFetched"`
	FetchedWeeks int       `json:"fetchedWeeks"`
	Repo         string    `json:"repo"`
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

func getConfigFile() (string, error) {
	if _, err := os.Stat("config.yaml"); err == nil {
		return "config.yaml", nil
	} else if os.IsNotExist(err) {
		if _, err := os.Stat(defaultConfigFile); err == nil {
			fmt.Printf("Warning: Using sample %s file. Please create and commit your own config.yaml file to track your feeds.\n", defaultConfigFile)
			return defaultConfigFile, nil
		}
	}
	return "", fmt.Errorf("config file not found")
}

func main() {
	// Use sample config if config.yaml doesn't exist
	configFile, err := getConfigFile()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Read and parse the YAML file
	yamlFile, err := os.ReadFile(configFile)
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

	// Set default values if not provided
	articlesPerPage := defaultArticlesPerPage
	if config.ArticlesPerPage != 0 {
		articlesPerPage = config.ArticlesPerPage
	}
	fetchWeeks := defaultFetchWeeks
	if config.FetchWeeks != 0 {
		fetchWeeks = config.FetchWeeks
	}

	repo := os.Getenv("GITHUB_REPOSITORY")
	if repo == "" {
		repo = defaultRepo
	}

	var articles []Article
	parser := gofeed.NewParser()
	ageLimit := time.Now().AddDate(0, 0, -7*fetchWeeks)

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
		feedHomepage := parsedFeed.Link

		// Check if the feed has any items within the last month
		hasRecentItems := false
		for _, item := range parsedFeed.Items {
			if item.PublishedParsed != nil && item.PublishedParsed.After(ageLimit) {
				hasRecentItems = true
				fmt.Printf("Retrieving [%s]...\n", feedTitle)
				break
			}
		}

		if !hasRecentItems {
			fmt.Printf("Skipping [%s]: No updates within the last %d week(s).\n", feedTitle, fetchWeeks)
			continue
		}

		for _, item := range parsedFeed.Items {
			if item.PublishedParsed == nil {
				continue // Skip items without a valid publication date
			}

			if item.PublishedParsed.Before(ageLimit) {
				continue // Skip items older than one month
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
				Homepage:   feedHomepage,
			})
		}
	}

	// Sort articles by published date, newest first
	sort.Slice(articles, func(i, j int) bool {
		return articles[i].Published.After(articles[j].Published)
	})

	// Calculate the number of pages
	totalPages := int(math.Ceil(float64(len(articles)) / float64(articlesPerPage)))

	// Create a directory for the paginated JSON files
	err = os.MkdirAll("web/articles", 0755)
	if err != nil {
		fmt.Printf("Error creating articles directory: %v\n", err)
		return
	}

	// Remove existing JSON files
	files, err := os.ReadDir("web/articles")
	if err != nil {
		fmt.Printf("Error reading articles directory: %v\n", err)
		return
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			err = os.Remove(fmt.Sprintf("web/articles/%s", file.Name()))
			if err != nil {
				fmt.Printf("Error removing existing JSON file: %v\n", err)
				return
			}
		}
	}

	// Create paginated JSON files
	for page := 1; page <= totalPages; page++ {
		start := (page - 1) * articlesPerPage
		end := start + articlesPerPage
		if end > len(articles) {
			end = len(articles)
		}

		pageArticles := articles[start:end]

		file, err := os.Create(fmt.Sprintf("web/articles/page_%d.json", page))
		if err != nil {
			fmt.Printf("Error creating file for page %d: %v\n", page, err)
			continue
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		err = encoder.Encode(pageArticles)
		if err != nil {
			fmt.Printf("Error encoding JSON for page %d: %v\n", page, err)
		}
	}

	// Create a metadata file with page count, last fetch time, and number of weeks fetched
	metadataFile, err := os.Create("web/articles/metadata.json")
	if err != nil {
		fmt.Printf("Error creating metadata file: %v\n", err)
		return
	}
	defer metadataFile.Close()

	metadata := Metadata{
		TotalPages:   totalPages,
		LastFetched:  time.Now().UTC(),
		FetchedWeeks: fetchWeeks,
		Repo:         repo,
	}
	encoder := json.NewEncoder(metadataFile)
	err = encoder.Encode(metadata)
	if err != nil {
		fmt.Printf("Error encoding metadata JSON: %v\n", err)
	}

	fmt.Println("Articles fetched, sorted, and saved successfully.")
}
