// Package cmd provides the command line interface to scrape cannaconnection.com
package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/droxey/strainscan/models"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/queue"
	. "github.com/logrusorgru/aurora"
	cmap "github.com/orcaman/concurrent-map"
	"github.com/spf13/cobra"
)

const (
	visitExternalLinks  = false
	debugging           = true
	fileName            = "results.json"
	queryParams         = "?q=&p="
	baseURL             = "https://sdb.openthc.org"
	searchURL           = baseURL + "/search"
	strainURL           = baseURL + "/strain/"
	minValidity         = 2
	minStrainNameLength = 2
	threads             = 8
	maxDepth            = 3
	filePermissions     = 0644
	runAsync            = true
)

func scrapeAll(cmd *cobra.Command, args []string) {
	startTime := time.Now()
	strainMap := cmap.New()
	hasLastPage := false
	firstPageURL := searchURL + queryParams + "1"
	strainCount := 0

	c := setupCollector()
	q, _ := queue.New(threads, &queue.InMemoryQueueStorage{MaxSize: 10000})

	// OpenTHC (List): Get Total # Of Pages
	c.OnHTML("body > div > div.page-list-control-wrap > div > a:nth-child(20)", func(e *colly.HTMLElement) {
		if hasLastPage {
			return
		}

		lastPageButtonURL := e.Attr("href")
		lastPageStr := strings.Replace(lastPageButtonURL, "/search"+queryParams, "", 1)
		lastPageNumber, _ := strconv.Atoi(lastPageStr)

		for i := 2; i <= lastPageNumber; i++ {
			url := searchURL + queryParams + strconv.Itoa(i)
			q.AddURL(url)
		}

		hasLastPage = true
		q.Run(c)
	})

	// OpenTHC (List): Get Data From Search Table
	c.OnHTML("body > div > div.table-responsive.mt-4 > table > tbody > tr", func(e *colly.HTMLElement) {
		s := &models.Strain{}
		e.Unmarshal(s)

		s.URL = baseURL + s.URL
		validity, _ := strconv.Atoi(s.Validity)

		isCleanURL := strings.Contains(s.URL, "/strain/") && (len(s.Slug) >= minStrainNameLength)
		isValidURL := validity >= minValidity
		if isValidURL && isCleanURL {
			strainCount++
			strainMap.Set(s.Slug, s)
			c.Visit(s.URL)
		}
	})

	// OpenTHC (Details): Get External Links From Strain Page
	c.OnHTML("body > div > section:nth-child(3) > div", func(e *colly.HTMLElement) {
		l := &models.StrainLinks{}
		e.Unmarshal(l)

		links := []string{l.LeaflyURL, l.AllBudURL, l.WikileafURL, l.KannapediaURL}
		s := GetStrainForURL(strainMap, e.Request.URL.String())

		for _, link := range links {
			isExtLink := !strings.HasPrefix(link, baseURL) && strings.HasPrefix(link, "http")
			if isExtLink {
				s.Links = append(s.Links, link)
			}
		}
	})

	fmt.Println(Bold("\nStarting scan...\n"))

	c.Visit(firstPageURL)
	c.Wait()

	output, _ := json.MarshalIndent(strainMap, "", "  ")
	ioutil.WriteFile(fileName, output, filePermissions)

	fmt.Println(
		Sprintf("%s %d strains found in %2.2f minutes and saved to %s.",
			Inverse("[DONE]").Bold(),
			Green(strainCount).Bold(),
			Green(time.Now().Sub(startTime).Minutes()).Bold(),
			Blue(fileName).Bold()))
	os.Exit(1)
}

func setupCollector() *colly.Collector {
	c := colly.NewCollector(
		colly.MaxDepth(maxDepth),
		colly.Async(runAsync),
		colly.CacheDir("./.cache"),
	)

	extensions.RandomUserAgent(c)
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: threads,
		RandomDelay: 5 * time.Second})

	c.WithTransport(&http.Transport{
		DisableKeepAlives: true,
		DialContext: (&net.Dialer{
			Timeout: 20 * time.Second,
		}).DialContext,
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", RandomString())

		if debugging {
			fmt.Println(Sprintf("%s %s", Green("[REQ]"), r.URL.String()))
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		if debugging {
			fmt.Println(Sprintf("%s %d %s: %s", Red("[ERR]"), Red(r.StatusCode).Bold(), Red(err.Error()).Bold(), r.Request.URL.String()))
		}
	})

	return c
}

func log(msg string) {
	fmt.Println(Sprintf("%s %s", White("[DEBUG]").BgGray(8-1), Blue(msg)))
}

func init() {
	var listCmd = &cobra.Command{
		Use:   "list",
		Short: "",
		Long:  ``,
		Run:   scrapeAll,
	}

	rootCmd.AddCommand(listCmd)
}
