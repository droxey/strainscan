package models

// Strain stores all data scraped from the website.
type Strain struct {
	Name        string
	Description string
	Parents     []string
	Image       string
}
