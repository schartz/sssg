package helpers

import "strings"

func MakeTitle(in_title string) string {
	in_title = strings.Replace(in_title, "_", " ", -1)
	in_title = strings.Replace(in_title, ".md", " ", 1)
	in_title = strings.Title(strings.ToLower(in_title))

	return in_title
}