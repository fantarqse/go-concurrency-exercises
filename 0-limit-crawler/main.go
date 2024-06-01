//////////////////////////////////////////////////////////////////////
//
// Your task is to change the code to limit the crawler to at most one
// page per second, while maintaining concurrency (in other words,
// Crawl() must be called concurrently)
//
// @hint: you can achieve this by adding 3 lines
//

package main

import (
	"fmt"
	"sync"
	"time"
)

// Crawl uses `fetcher` from the `mockfetcher.go` file to imitate a
// real crawler. It crawls until the maximum depth has reached.
func Crawl(url string, depth int, wg *sync.WaitGroup, rate <-chan time.Time) {
	defer wg.Done()

	if depth <= 0 {
		return
	}

	<-rate
	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	// NOTE:
	// '%q' allows to safely escape a string and add quotes to it.
	fmt.Printf("found: %s %q\n", url, body)

	wg.Add(len(urls))
	for _, u := range urls {
		// Do not remove the `go` keyword, as Crawl() must be
		// called concurrently
		go Crawl(u, depth-1, wg, rate)
	}
	return
}

func main() {
	// NOTE:
	// 'WaitGroup' is a synchronization tool
	// that allows us to sync goroutines by adding them to the collection (wg.Add(n)).
	// 'wg.Wait()' blocks the program and waits when the counter is a zero.
	var wg sync.WaitGroup

	// NOTE:
	// 'time.Tick(d)' is a wrapper for 'time.NewTicker(d)', but provides access only to channel.
	// time.Tick(d) == time.NewTicker(d).C
	// 'time.NewTicker().C' has a '<-chan time.Time' type which means it is a receive-only channel.
	//
	// 1) 'chan Type' is a bidirectional channel.
	// 2) 'chan<- Type' is a send-only channel.
	// 3) '<-chan Type' is a receive-only channel.
	rate := time.Tick(time.Second * 1)

	wg.Add(1)
	Crawl("http://golang.org/", 4, &wg, rate)
	wg.Wait()
}
