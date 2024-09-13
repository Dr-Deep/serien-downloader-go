package sto

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func removeSliceElement[T any](idx int, slice []T) []T {
	var newSlice = make([]T, len(slice)-1)

	for i, e := range slice {
		if i == idx {
			continue
		}
		newSlice[i-1] = e
	}

	return newSlice
}

func removeDuplicates[T any](slice []T) []T {
	var (
		newslice = make([]T, 0, len(slice))
		seen     = make(map[any]bool)
	)

	for _, e := range slice {
		if _, ok := seen[e]; !ok {
			newslice = append(newslice, e)
			seen[e] = true
		}
	}

	return newslice
}

func removeHTMLFragmentsFromString(htmlString string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlString))
	if err != nil {
		return htmlString
	}

	return doc.Text()
}
