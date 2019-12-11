// Package cmd provides the command line interface to scrape cannaconnection.com
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/droxey/strainscan/models"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/queue"
	. "github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

const (
	fileName        = "results.json"
	maxDepth        = 2
	debugging       = false
	baseURL         = "https://www.cannaconnection.com/strains?show_char="
	numberOfLetters = 26
	threads         = 4
	sep             = " "
)

func scrapeAll(cmd *cobra.Command, args []string) {
	startTime := time.Now()
	alphabet := LowercaseAlphabet(numberOfLetters)
	strainPages := make(map[string]*models.Strain, 0)
	c := colly.NewCollector(
		colly.MaxDepth(maxDepth),
		colly.Async(false),
		colly.CacheDir("./.cache"),
	)

	extensions.RandomUserAgent(c)
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: threads})
	c.WithTransport(&http.Transport{
		DisableKeepAlives: true,
	})

	q, _ := queue.New(threads, &queue.InMemoryQueueStorage{MaxSize: 10000})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", RandomString())
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("[ERR] ", r.Request.URL, "\t", r.StatusCode, "\n\tError:\n\t", err)
	})

	c.OnHTML("meta", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()
		if !strings.Contains(url, "/strains/") {
			return
		}

		if e.Attr("property") == "og:description" {
			strainPages[url].Description = e.Attr("content")
		}

		if e.Attr("property") == "og:image" {
			strainPages[url].Image = e.Attr("content")
		}
	})

	c.OnHTML(".product .post-content", func(e *colly.HTMLElement) {
		url := e.Request.URL.String()
		s := strainPages[url]

		parentSelector := "div.product-top > div.pb-right-column > div > div > div:nth-child(2) > div.multi-feature.feature-value"
		parentFirst := strings.Trim(e.DOM.Find(parentSelector+".first a").Text(), " ")
		parentLast := strings.Trim(e.DOM.Find(parentSelector+".last a").Text(), " ")
		s.Parents = append(s.Parents, parentFirst)
		s.Parents = append(s.Parents, parentLast)
	})

	c.OnHTML("#strains_page ul.strains-list li a", func(e *colly.HTMLElement) {
		name := e.Text
		url := e.Attr("href")
		s := &models.Strain{
			Name:    name,
			Parents: make([]string, 0),
		}
		strainPages[url] = s
		q.AddURL(url)
	})

	CreateUI(sep)
	for i := range alphabet {
		pageURL := baseURL + alphabet[i]
		c.Visit(pageURL)
		c.Wait()
		UpdateUI(sep)
	}

	q.Run(c)

	fmt.Println("\n---\n[FILE] Outputting file.")
	output, _ := json.MarshalIndent(strainPages, "", "  ")
	ioutil.WriteFile(fileName, output, 0644)

	diff := time.Now().Sub(startTime).Seconds()
	fmt.Println(Sprintf("%s %d strains found in %2.f seconds.", Gray(1-1, "[DONE]").BgGray(24-1), Green(len(strainPages)).Bold(), Green(diff).Bold()))
	os.Exit(1)
}

func init() {
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "",
		Long:  ``,
		// Use the function reference instead of a pass-through function for better testability.
		Run: scrapeAll,
	}

	rootCmd.AddCommand(listCmd)
}
