package models

type StrainLinks struct {
	URL           string
	LeaflyURL     string `selector:"a.btn:nth-child(1):not([href*='https://duckduckgo.com/?q='])" attr:"href"`
	AllBudURL     string `selector:"a.btn:nth-child(2):not([href*='https://duckduckgo.com/?q='])" attr:"href"`
	WikileafURL   string `selector:"a.btn:nth-child(3):not([href*='https://duckduckgo.com/?q='])" attr:"href"`
	KannapediaURL string `selector:"a.btn:nth-child(4):not([href*='https://duckduckgo.com/?q='])" attr:"href"`
}

// Strain stores all data scraped from the website.
type Strain struct {
	Name     string `selector:"td:nth-child(1) > a"`
	URL      string `selector:"td:nth-child(1) > a" attr:"href"`
	Race     string `selector:"td:nth-child(2)"`
	Validity string `selector:"td:nth-child(3)"`
	Slug     string `selector:"td:nth-child(4)"`
	Links    []string
}
