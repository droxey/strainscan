// Package cmd provides the command line interface to scrape cannaconnection.com
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/droxey/strainscraper/models"
	"github.com/gocolly/colly"
	"github.com/gosimple/slug"
	. "github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var scraper = &cobra.Command{
	Use:   "scrape [strain name]",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		strainSlug := slug.Make(args[0])
		strainURL := "https://www.cannaconnection.com/strains/" + strainSlug
		strainData := make([]*models.Strain, 0)

		c := colly.NewCollector()

		c.OnHTML(".product .post-content", func(e *colly.HTMLElement) {
			newStrain := &models.Strain{
				Features: make([]*models.Feature, 0),
				Parents:  new(models.Parents),
			}

			e.Unmarshal(newStrain)
			strainData = append(strainData, newStrain)
		})

		c.OnHTML(".data-sheet .feature-wrapper", func(e *colly.HTMLElement) {
			feature := &models.Feature{}
			e.Unmarshal(feature)
			strainData[0].Features = append(strainData[0].Features, feature)
		})

		c.OnHTML(".data-sheet .multifeature-wrapper:first", func(e *colly.HTMLElement) {
			p := &models.Parents{}
			e.Unmarshal(p)
			strainData[0].Parents = p
		})

		err := c.Visit(strainURL)
		if err != nil {
			fmt.Printf("%s '%s' could not be found. Please check your spelling and try again.\n", White("[ERROR]").Bold().BgRed(), Red(args[0]).Bold())
			fmt.Printf("%s URL visited: %s", White("[DEBUG]").Bold().BgBlue(), Yellow(strainURL))
			os.Exit(1)
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.Encode(strainData)
	},
}

func init() {
	rootCmd.AddCommand(scraper)
}
