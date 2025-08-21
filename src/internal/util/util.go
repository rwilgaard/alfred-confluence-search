package util

import (
	"fmt"
	"regexp"
	"strings"
)

type ParsedQuery struct {
	Text      string
	Spaces    []string
	FullQuery string
}

func ParseQuery(query string) *ParsedQuery {
	q := new(ParsedQuery)
	spaceRegex := regexp.MustCompile(`^@\w+`)
	q.FullQuery = query

	for _, w := range strings.Split(query, " ") {
		switch {
		case spaceRegex.MatchString(w):
			s := strings.ReplaceAll(w[1:], "_", " ")
			q.Spaces = append(q.Spaces, s)
		default:
			q.Text = q.Text + w + " "
		}
	}

	return q
}

func Autocomplete(query string) string {
	for _, w := range strings.Split(query, " ") {
		switch w {
		case "@":
			return "spaces"
		}
	}
	return ""
}

func BuildJQL(query *ParsedQuery) (jql string) {
	var conditions []string
	typeCQL := "type = page"

	if text := strings.TrimSpace(query.Text); text != "" {
		text = strings.ReplaceAll(text, "å", "å")
		conditions = append(conditions, fmt.Sprintf("siteSearch ~ '%s'", text))
    }

	addClause := func(field string, values []string) {
		if len(values) > 0 {
			conditions = append(conditions, fmt.Sprintf("%s in ('%s')", field, strings.Join(values, "','")))
		}
	}

	addClause("space.key", query.Spaces)

	if len(conditions) == 0 {
		return typeCQL
	}

    return strings.Join(conditions, " AND ") + " AND " + typeCQL
}
