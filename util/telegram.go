package util

import (
	"strings"

	"golang.org/x/net/html"
)

func EscapeMarkdownV2(text string) string {
	specialChars := []string{
		"_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!",
	}
	for _, char := range specialChars {
		text = strings.ReplaceAll(text, char, `\`+char)
	}
	return text
}

func ExtractTextFromHTML(htmlStr string) string {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		return htmlStr
	}
	var f func(*html.Node)
	var output strings.Builder
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			output.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return output.String()
}
