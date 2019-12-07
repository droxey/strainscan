// Package cmd provides the command line interface to scrape cannaconnection.com
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/droxey/strainscrape/models"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/queue"
	// "github.com/gocolly/colly/debug"
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
		colly.Async(true),
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

	c.OnHTML(".product .post-content", func(e *colly.HTMLElement) {

	})

	c.OnHTML("#strains_page ul.strains-list li a", func(e *colly.HTMLElement) {
		s := &models.Strain{
			Name:     e.Text,
			Features: make([]*models.Feature, 0),
			Parents:  new(models.Parents),
			URL:      e.Attr("href"),
		}

		strainPages[s.Name] = s
		q.AddURL(s.URL)
	})

	// c.OnHTML(".data-sheet .feature-wrapper", func(e *colly.HTMLElement) {
	// 	feature := &models.Feature{}
	// 	e.Unmarshal(feature)
	// 	feature.Value = strings.Title(strings.ReplaceAll(feature.Value, "-", " "))
	// 	strainPages[0].Features = append(strainPages[0].Features, feature)

	// })

	// c.OnHTML(".data-sheet .multifeature-wrapper:first", func(e *colly.HTMLElement) {
	// 	p := &models.Parents{}
	// 	e.Unmarshal(p)
	// 	p.First = strings.Title(strings.ReplaceAll(p.First, "-", " "))
	// 	p.Last = strings.Title(strings.ReplaceAll(p.Last, "-", " "))
	// 	strainPages[0].Parents = p
	// })

	CreateUI(sep)
	for i := range alphabet {
		pageURL := baseURL + alphabet[i]
		c.Visit(pageURL)
		c.Wait()
		UpdateUI(sep)
	}

	output, _ := json.MarshalIndent(strainPages, "", "  ")
	ioutil.WriteFile(fileName, output, 0644)

	diff := time.Now().Sub(startTime).Seconds()
	fmt.Println(Sprintf("\n\n%s %d strains found in %2.f seconds.", Gray(1-1, "[DONE]").BgGray(24-1), Green(len(strainPages)).Bold(), Green(diff).Bold()))
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
