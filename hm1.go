package hm1

import (
	"fmt"
	"hash/fnv"
	"image"
	"image/color"
	"io"
	"math"
	"strings"
	"sync"

	"golang.org/x/tour/tree"
)

// Краще було створювати новий go package для кожного завдання
// В Go Tour є секція про go packages

// Exercise: Loops and Functions

func Sqrt(x float64) float64 {
	if x-0.0 < 0.00001 {
		return 0.0
	}

	z := x / 2.0
	delta := 10.0

	for i := 0; i < 10 && delta > 0.01; i++ {
		delta = z
		z -= (z*z - x) / (2 * z)
		delta = math.Abs(delta - z)
		fmt.Println(z)
	}

	return z
}

// Exercise: Slices

func Pic(dx, dy int) [][]uint8 {
	result := make([][]uint8, dy)
	for y := range result {
		result[y] = make([]uint8, dx)

		for x := range result[y] {
			result[y][x] = uint8((x + y) / 2)
		}
	}
	return result
}

// Exercise: Maps

func WordCount(s string) map[string]int {
	result := make(map[string]int)
	for _, word := range strings.Fields(s) {
		result[word]++
	}
	return result
}

// Exercise: Fibonacci closure

func Fibonacci() func() int {
	a, b := 0, 1

	fib := func() int {
		result := a
		a, b = b, a+b
		return result
	}

	return fib
}

// Exercise: Stringers

type IPAddr [4]byte

func (ipAddr IPAddr) String() string {
	return fmt.Sprintf("%d.%d.%d.%d", ipAddr[0], ipAddr[1], ipAddr[2], ipAddr[3])
}

// Exercise: Errors

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("cannot Sqrt negative number: %f", float64(e))
}

func SqrtV2(x float64) (float64, error) {
	if x-0.0 < 0.00001 {
		return -1, ErrNegativeSqrt(x)
	}

	z := x / 2.0
	delta := 10.0

	for i := 0; i < 10 && delta > 0.01; i++ {
		delta = z
		z -= (z*z - x) / (2 * z)
		delta = math.Abs(delta - z)
		fmt.Println(z)
	}

	return z, nil
}

// Exercise: Readers

type MyReader struct{}

func (mr MyReader) Read(out []byte) (int, error) {
	out[0] = 65
	return 1, nil
}

// Exercise: rot13Reader

var ROT13 = map[byte]byte{
	97: 110, 98: 111, 99: 112, 100: 113, 101: 114, 102: 115, 103: 116,
	104: 117, 105: 118, 106: 119, 107: 120, 108: 121, 109: 122, 110: 97,
	111: 98, 112: 99, 113: 100, 114: 101, 115: 102, 116: 103, 117: 104,
	118: 105, 119: 106, 120: 107, 121: 108, 122: 109,
	65: 78, 66: 79, 67: 80, 68: 81, 69: 82, 70: 83, 71: 84,
	72: 85, 73: 86, 74: 87, 75: 88, 76: 89, 77: 90, 78: 65,
	79: 66, 80: 67, 81: 68, 82: 69, 83: 70, 84: 71, 85: 72,
	86: 73, 87: 74, 88: 75, 89: 76, 90: 77, 32: 45,
}

type rot13Reader struct {
	r io.Reader
}

func (rot13 rot13Reader) Read(out []byte) (int, error) {
	innerOut := make([]byte, len(out))
	n, e := rot13.r.Read(innerOut)
	if e == nil {
		for i, value := range innerOut {
			out[i] = ROT13[value]
		}
	}
	return n, e
}

// Exercise: Images

type Image struct{}

func (img Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (img Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, 256, 256)
}

func (img Image) At(x, y int) color.Color {
	v := uint8((x + y) / 2)
	return color.RGBA{v, v, 255, 255}
}

// Exercise: Equivalent Binary Trees

func Walk(t *tree.Tree, ch chan int) {
	// канали теж краще закривати по defer якщо є змога тому що він тоді закриється і в разі якшо відбудеться паника
	defer close(ch)

	var walk func(*tree.Tree)
	walk = func(t *tree.Tree) {
		if t == nil {
			return
		}
		walk(t.Left)
		ch <- t.Value
		walk(t.Right)
	}

	walk(t)
}

func Same(t1, t2 *tree.Tree) bool {
	c1 := make(chan int)
	c2 := make(chan int)
	go Walk(t1, c1)
	go Walk(t2, c2)

	for {
		v1, ok1 := <-c1
		v2, ok2 := <-c2

		if ok1 != ok2 {
			return false
		}

		if v1 != v2 {
			return false
		}

		if !ok1 && !ok2 {
			return true
		}
	}
}

// Exercise: Web Crawler

type Fetcher interface {
	Fetch(url string) (body string, urls []string, err error)
}

func Crawl(url string, depth int, fetcher Fetcher, routes chan string) {
	visited := make(map[uint32]bool)
	var mu sync.Mutex

	defer close(routes)

	var crawl func(url string, depth int, fetcher Fetcher)
	crawl = func(url string, depth int, fetcher Fetcher) {
		mu.Lock()
		isVisited := visited[hash(url)]
		mu.Unlock()
		if isVisited {
			return
		}

		mu.Lock()
		visited[hash(url)] = true
		mu.Unlock()

		if depth <= 0 {
			return
		}
		_, urls, err := fetcher.Fetch(url)
		if err != nil {
			// fmt.Println(err)
			return
		}

		// fmt.Printf("found: %s %q\n", url, body)

		routes <- url

		for _, u := range urls {
			crawl(u, depth-1, fetcher)
		}
	}

	crawl(url, depth, fetcher)
}

type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
