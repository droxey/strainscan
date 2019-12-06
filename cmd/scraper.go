/*
Copyright © 2019 Dani Roxberry <dani@bitoriented.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"github.com/gocolly/colly"
)

var scraper = &cobra.Command{
	Use:   "scrape",
	Short: "",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		type feature struct {
			Name string `selector:".feature-title"`
			Value string `selector:".feature-value"`
		}

		type parents struct {
			First string `selector:".first"`
			Last string `selector:".last"`
		}

		type strain struct {
			Name string `selector:"h1"`
			Features []*feature
			Parents *parents
		}

		strains := make([]*strain, 0)

		c := colly.NewCollector()

		c.OnHTML(".product .post-content", func(e *colly.HTMLElement) {
			newStrain := &strain{
				Features: make([]*feature, 0),
				Parents: new(parents),
			}

			e.Unmarshal(newStrain)
			strains = append(strains, newStrain)
		})

		c.OnHTML(".data-sheet .feature-wrapper", func(e *colly.HTMLElement) {
			feature := &feature{}
			e.Unmarshal(feature)
			strains[0].Features = append(strains[0].Features, feature)
		})

		c.OnHTML(".data-sheet .multifeature-wrapper:first", func(e *colly.HTMLElement) {
			p := &parents{}
			e.Unmarshal(p)
			strains[0].Parents = p
		})

		c.Visit("https://www.cannaconnection.com/strains/pineapple-kush")

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		fmt.Print("\n-----\n\n")
		enc.Encode(strains)
	},
}

func init() {
	rootCmd.AddCommand(scraper)
}
