package filter

import (
	"testing"
	"time"
)

func TestFilterKeyword(t *testing.T) {
	var tests = []struct {
		keyword string
		expect  string
		str     string
	}{
		{
			keyword: "",
			expect:  "str",
			str:     "str",
		},
		{
			keyword: "key",
			str:     "has key",
			expect:  "",
		},
		{
			keyword: "nokey",
			str:     "has key",
			expect:  "has key",
		},
		{
			keyword: "key|text",
			str:     "has key",
			expect:  "",
		},
		{
			keyword: "key|text",
			str:     "has text",
			expect:  "",
		},
	}

	for _, test := range tests {
		if test.expect != FilterKeyword(test.keyword)(test.str) {
			t.Errorf("filter key error:%v expect:%v", test.keyword, test.expect)
		}
	}
}

func TestFilterTimeLimiter(t *testing.T) {
	var tests = []struct {
		duration time.Duration
		count    int64
	}{
		{
			duration: time.Second,
			count:    10,
		},
		{
			duration: time.Second,
			count:    60,
		},
		{
			duration: time.Minute,
			count:    10,
		},
	}

	for _, test := range tests {
		limit := FilterTimeLimiter(test.duration, test.count)
		if "test" != limit.Filter("test") {
			t.Errorf("time limiter filter error")
		}

		time.Sleep(time.Duration(int64(test.duration) / (test.count * 2)))

		if "" != limit.Filter("test") {
			t.Errorf("time limiter filter error")
		}

		time.Sleep(time.Duration(int64(test.duration) / (test.count * 2)))

		if "test" != limit.Filter("test") {
			t.Errorf("time limiter filter error")
		}
	}
}
