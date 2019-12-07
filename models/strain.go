package models

// Feature provides named attributes for each strain.
type Feature struct {
	Name  string `selector:".feature-title", json:"name"`
	Value string `selector:".feature-value", json:"value"`
}

// Parents describes the lineage of the strain.
type Parents struct {
	First string `selector:".first", json:"first"`
	Last  string `selector:".last", json:"last"`
}

// Strain stores all data scraped from the website.
type Strain struct {
	Name     string `json:"name"`
	URL  	 string `json:"url"`
	Features []*Feature
	Parents  *Parents
}
