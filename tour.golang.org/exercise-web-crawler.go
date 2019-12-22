package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

type Synchronization struct {
	urls_map  map[string]bool
	mutex sync.Mutex
	queue []string
	queue_len int
	wg sync.WaitGroup
}

// thread function
// gets url from queue, fetches additional urls, appends new urls to queue
func Crawl(fetcher Fetcher, urls_meta *Synchronization) {
	defer urls_meta.wg.Done()
	
	urls_meta.mutex.Lock()
	if (urls_meta.queue_len == 0) {
		urls_meta.mutex.Unlock()
		return	
	}
	url := urls_meta.queue[0]
	urls_meta.queue = urls_meta.queue[1:] // pop from queue
	urls_meta.queue_len -= 1
	urls_meta.mutex.Unlock()

	body, urls, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("found: %s %q\n", url, body)
	for _, u := range urls {
		urls_meta.mutex.Lock()
 		if _, v := urls_meta.urls_map[u]; !v  {
			urls_meta.urls_map[u] = true
			urls_meta.queue_len += 1
			urls_meta.queue = append(urls_meta.queue, u)
		}
		urls_meta.mutex.Unlock()
	}
	return
}

func main() {
	n_concurrent_urls := 3
	seed_url := "https://golang.org/"
	urls_meta := Synchronization{urls_map: make(map[string]bool), queue: make([]string, 0), queue_len: 0}	
	urls_meta.queue = append(urls_meta.queue,seed_url)
	urls_meta.urls_map[seed_url] = true
	urls_meta.queue_len += 1
	
	for {
		urls_meta.mutex.Lock()
		if (urls_meta.queue_len == 0) {
			urls_meta.mutex.Unlock()
			return
		}
		urls_meta.mutex.Unlock()
		
		urls_meta.wg.Add(n_concurrent_urls)
		for i := 0; i < n_concurrent_urls; i++ {
			go Crawl(fetcher, &urls_meta)
		}
		urls_meta.wg.Wait() // join all threads
	}																	
}

// fakeFetcher is Fetcher that returns canned results.
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

// fetcher is a populated fakeFetcher.
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
