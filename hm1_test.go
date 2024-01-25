package hm1

import (
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"
	"testing"
	"sort"
)

type sqrtCase struct{ x, out float64 }
type sqrtTest struct {
	name  string
	cases []sqrtCase
}

var sqrtTests = []sqrtTest{
	{
		name:  "Negative",
		cases: []sqrtCase{{-1.0, 0}},
	},
	// Move a value from input to output
	{
		name:  "Zero",
		cases: []sqrtCase{{0, 0}},
	},
	{
		name:  "256",
		cases: []sqrtCase{{256.0, 16.0}},
	},
}

type picCase struct{ dx, dy, outLength int }
type picTest struct {
	name  string
	cases []picCase
}

var picTests = []picTest{
	{
		name:  "Zero",
		cases: []picCase{{0, 0, 0}},
	},
	// Move a value from input to output
	{
		name:  "3х3",
		cases: []picCase{{3, 3, 9}},
	},
	{
		name:  "256x256",
		cases: []picCase{{256, 256, 65_536}},
	},
}

type wordCountCase struct {
	in  string
	out map[string]int
}
type wordCountTest struct {
	name  string
	cases []wordCountCase
}

var wordCountTests = []wordCountTest{
	{
		name: "Just sentence",
		cases: []wordCountCase{
			{"A man a plan a canal panama.", map[string]int{"A": 1, "a": 2, "canal": 1, "man": 1, "panama.": 1, "plan": 1}},
			{"I ate a donut. Then I ate another donut.", map[string]int{"I": 2, "Then": 1, "a": 1, "another": 1, "ate": 2, "donut.": 2}},
		},
	},
	// Move a value from input to output
	{
		name: "Empty sentence",
		cases: []wordCountCase{
			{"", map[string]int{}},
			{"    ", map[string]int{}},
		},
	},
}

type crawlCase struct {
	startUrl  string
	expectToVisit []string
}
type crawlTest struct {
	name  string
	cases []crawlCase
}

var crawlTests = []crawlTest{
	{
		name: "Start from https://golang.org/pkg/",
		cases: []crawlCase{
			{
				"https://golang.org/pkg/",
				 []string{
					"https://golang.org/",
					"https://golang.org/pkg/",
					"https://golang.org/pkg/fmt/",
					"https://golang.org/pkg/os/",
						},
			},
			{
				"https://golang.org/",
				[]string{
					"https://golang.org/",
					"https://golang.org/pkg/",
					"https://golang.org/pkg/fmt/",
					"https://golang.org/pkg/os/",
					},
			},
		},
	},
}


func TestCompute(t *testing.T) {
	for _, test := range sqrtTests {
		t.Run(test.name, func(t *testing.T) { testSqrt(t, test) })
	}

	for _, test := range picTests {
		t.Run(test.name, func(t *testing.T) { testPic(t, test) })
	}

	for _, test := range wordCountTests {
		t.Run(test.name, func(t *testing.T) { testWordCount(t, test) })
	}

	testROT13()

	for _, test := range crawlTests {
		t.Run(test.name, func(t *testing.T) { testCrawlRunCases(t, test) })
	}
}

func testSqrt(t *testing.T, test sqrtTest) {
	for _, c := range test.cases {

		actual := Sqrt(c.x)

		if math.Abs(actual-c.out) > 0.00001 {
			t.Fatalf("Expected f(%v) to be %v, not %v", c.x, c.out, actual)
		}
	}
}

func totalLen(slice [][]uint8) int {
	var result int
	for _, dy := range slice {
		result += len(dy)
	}
	return result
}

func testPic(t *testing.T, test picTest) {
	for _, c := range test.cases {

		actual := Pic(c.dx, c.dy)
		actualTotalLen := totalLen(actual)

		if actualTotalLen != c.outLength {
			t.Fatalf("Expected Pic(%v, %v) to be %v, not %v", c.dx, c.dy, c.outLength, actualTotalLen)
		}
	}
}

func testWordCount(t *testing.T, test wordCountTest) {
	for _, c := range test.cases {

		actual := WordCount(c.in)

		if !reflect.DeepEqual(actual, c.out) {
			t.Fatalf("Expected WordCount(\"%v\") to be %v, not %v", c.in, c.out, actual)
		}
	}
}

func testROT13() {
	s := strings.NewReader("Create new Reader structure from string")
	r := rot13Reader{s}

	b := make([]byte, 8)
	for {
		n, err := r.Read(b)
		fmt.Printf("%q", b[:n])
		if err == io.EOF {
			fmt.Println()
			break
		}
	}

}

func testCrawlRunCases(t *testing.T, test crawlTest) {
	for _, c := range test.cases {
		routesChan := make(chan string)
		var visitedUrls []string
		expectedVisitedUrls := c.expectToVisit

		go Crawl(c.startUrl, 4, fetcher, routesChan)
		for value := range routesChan {
			visitedUrls = append(visitedUrls, value)
		}
		sort.Strings(visitedUrls)
		sort.Strings(expectedVisitedUrls)

		if !reflect.DeepEqual(visitedUrls, expectedVisitedUrls) {
			fmt.Println(len(visitedUrls))
			t.Fatalf("Expected Crawl() to be \n%v, not \n%v", expectedVisitedUrls, visitedUrls)
		}
	}
}

