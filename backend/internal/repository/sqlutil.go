package repository

import "strings"

func rebind(driver, query string) string {
	if driver != "pgx" {
		return query
	}
	var b strings.Builder
	n := 1
	for _, r := range query {
		if r == '?' {
			b.WriteByte('$')
			b.WriteString(itoa(n))
			n++
		} else {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func itoa(n int) string {
	if n < 10 {
		return string(rune('0' + n))
	}
	return itoa(n/10) + string(rune('0'+n%10))
}

func boolInt(v bool) int {
	if v {
		return 1
	}
	return 0
}
