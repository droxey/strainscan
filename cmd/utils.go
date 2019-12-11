// Package cmd provides the command line interface to scrape cannaconnection.com
package cmd

import (
	"fmt"
	"math/rand"
	"strings"

	. "github.com/logrusorgru/aurora"
)

// LowercaseAlphabet generates a lowercase alphabet.
func LowercaseAlphabet(numberOfLetters int) []string {
	alphabet := make([]string, numberOfLetters)
	for i := range alphabet {
		letter := 'a' + byte(i)
		alphabet[i] = string(letter)
	}
	return alphabet
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func CreateUI(sep string) {
	alphabet := LowercaseAlphabet(26)
	display := make([]string, 0)
	for i := range alphabet {
		stats := sep + alphabet[i] + sep
		display = append(display, stats)
	}
	fmt.Println("\n" + strings.Join(display, ""))
}

func UpdateUI(sep string) {
	fmt.Printf("%s", Gray(1-1, sep+" "+sep).BgGray(24-1))
}
