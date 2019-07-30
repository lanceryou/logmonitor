package filter

import (
	"strings"
	"time"
)

type Filter interface {
	Filter(string) string
}

func Filters(str string, fs []Filter) string {
	for _, f := range fs {
		str = f.Filter(str)
		if str == "" {
			return str
		}
	}
	return str
}

type FilterFunc func(string) string

func (f FilterFunc) Filter(str string) string {
	return f(str)
}

// support avoid send by keywords
func FilterKeyword(keyword string) FilterFunc {
	keys := strings.Split(keyword, "|")
	return func(str string) string {
		for _, key := range keys {
			if key != "" && strings.Contains(str, key) {
				return ""
			}
		}
		return str
	}
}

func FilterTimeLimiter(duration time.Duration, count int64) FilterFunc {
	next := time.Duration(int64(duration) / count)
	runTime := time.Now()
	return func(str string) string {
		now := time.Now()
		if now.Before(runTime) {
			return ""
		}

		runTime = now.Add(next)
		return str
	}
}
