package summarizer

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
)

func fetchHTML(url string) (*http.Response, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Unable to fetch via url: %s err: %w", url, err)
	}

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("Unable to get a proper response: %d %s", res.StatusCode, res.Status)
	}

	return res, nil
}

func ScrapeBerlinEvents(month string) (EventsSummary, error) {
	url := fmt.Sprintf("https://www.visitberlin.de/en/blog/%s-berlin", month)

	res, err := fetchHTML(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Parsed HTML successfully.")
	selector := "#main-content > div > article > div"

	var summary []map[string]string
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		fmt.Println("Looking at children or main article div")
		insideContent := false
		currentIndex := -1

		s.Children().Each(func(i int, child *goquery.Selection) {
			node := child.Get(0)

			if node.Type == html.ElementNode {
				switch node.Data {
				case "h2":
					summary = append(summary, map[string]string{"title": strings.TrimSpace(child.Text())})
					currentIndex++
					insideContent = true
					break
				case "p":
					if insideContent {
						targetIndex := summary[currentIndex]
						targetIndex["details"] += " " + strings.TrimSpace(child.Text())
					}
					break
				case "figure":
					break
				default:
					insideContent = false
				}
			}
		})
	})

	summarizer := ChatGPT{
		ApiKey: os.Getenv("CHATGPT_API_KEY"),
	}

	eventsInBerlinSummary, err := summarizer.SummarizeEvents(summary)
	if err != nil {
		return eventsInBerlinSummary, fmt.Errorf("Unable to build summary: %w", err)
	}

	return eventsInBerlinSummary, nil
}
