package util

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GoqueryImmediateText(s *goquery.Selection, sep string) string {
	elements := s.Contents().Map(func(_ int, s *goquery.Selection) string {
		if goquery.NodeName(s) == "#text" {
			return strings.TrimSpace(s.Text())
		}
		return ""
	})
	return strings.Join(elements, sep)
}
