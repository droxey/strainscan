// Package cmd provides the command line interface to scrape cannaconnection.com
package cmd

import (
	"math/rand"
	"strings"

	"github.com/droxey/strainscan/models"
	cmap "github.com/orcaman/concurrent-map"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandomString ...
func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// GetStrainForURL ...
func GetStrainForURL(strainMap cmap.ConcurrentMap, url string) *models.Strain {
	slug := strings.Replace(url, strainURL, "", 1)
	if tmp, ok := strainMap.Get(slug); ok {
		return tmp.(*models.Strain)
	}

	return nil
}
