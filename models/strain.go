package models

// Feature provides named attributes for each strain.
type Feature struct {
	Name  string `selector:".feature-title"`
	Value string `selector:".feature-value"`
}

// Parents describes the lineage of the strain.
type Parents struct {
	First string `selector:".first"`
	Last  string `selector:".last"`
}

// Strain stores all data scraped from the website.
type Strain struct {
	Name     string `selector:"h1"`
	Features []*Feature
	Parents  *Parents
}
