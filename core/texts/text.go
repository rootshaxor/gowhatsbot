package texts

import (
	"fmt"
	"strings"
)

func AddTab(c int, s ...string) string {
	var tmp []string
	for i := 0; i < c; i++ {
		tmp = append(tmp, "  ")
	}

	tmp = append(tmp, s...)
	return strings.Join(tmp, " ")
}

func ArrayS(s ...string) []string {
	return s
}

func QuoteBy(text string, qouter ...QouteChar) string {
	for _, quote := range qouter {
		text = fmt.Sprintf("%s%s%s", quote, strings.TrimSpace(text), quote)
	}
	return text
}

func ToUpperFirst(s string) string {
	sa := strings.Split(s, " ")
	for i := range sa {
		temps := strings.Split(strings.ToLower(sa[i]), "")
		for i := range temps {
			if i == 0 {
				temps[i] = (strings.ToUpper(temps[i]))
			}
		}
		sa[i] = strings.Join(temps, "")
	}

	return strings.Join(sa, " ")
}
