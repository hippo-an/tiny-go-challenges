package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/dev-hippo-an/tiny-go-challenges/qhn_05/client"
)

func main() {
	var port, numStories int

	flag.IntVar(&port, "port", 3001, "the port to start web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")

	flag.Parse()

	tpl := template.Must(template.ParseFiles("./index.gohtml"))

	http.HandleFunc("/", handler(numStories, tpl))

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {

	sc := storyCache{
		numStories: numStories,
	}

	go func() {
		ticker := time.NewTicker(4 * time.Second)

		for {
			temp := storyCache{
				numStories: numStories,
			}
			temp.stories()
			sc.mutex.Lock()
			sc.cache = temp.cache
			sc.expiration = temp.expiration
			sc.mutex.Unlock()
			<-ticker.C
		}
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		stories, err := sc.stories()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := templateData{
			Stories: stories,
			Time:    time.Since(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	}
}

type storyCache struct {
	numStories int
	cache      []item
	expiration time.Time
	mutex      sync.Mutex
}

func (sc *storyCache) stories() ([]item, error) {
	sc.mutex.Lock()
	defer sc.mutex.Unlock()
	if time.Since(sc.expiration) < 0 {
		log.Println("items from cached")

		return sc.cache, nil
	}

	stories, err := getTopStories(sc.numStories)
	if err != nil {
		return nil, err
	}

	sc.expiration = time.Now().Add(5 * time.Second)
	sc.cache = stories

	return sc.cache, nil
}

func getStories(ids []int) []item {
	type result struct {
		idx  int
		item item
		err  error
	}

	resultCh := make(chan result)

	for idx, id := range ids {
		go func(idx, id int) {
			var c client.Client
			hnItem, err := c.GetItem(id)
			if err != nil {
				resultCh <- result{idx: idx, err: err}
			}

			resultCh <- result{idx: idx, item: parseHNItem(hnItem)}
		}(idx, id)
	}

	var results []result

	for range ids {
		results = append(results, <-resultCh)
	}

	close(resultCh)

	sort.Slice(results, func(i, j int) bool {
		return results[i].idx < results[j].idx
	})

	var stories []item

	for i := 0; i < len(ids); i++ {
		if results[i].err != nil {
			continue
		}

		if isStoryLink(results[i].item) {
			stories = append(stories, results[i].item)
		}
	}

	return stories
}

func getTopStories(numStories int) ([]item, error) {
	var c client.Client
	ids, err := c.TopItems()
	if err != nil {
		return nil, errors.New("failed to load top stories")
	}

	var stories []item
	from := 0
	idsLength := len(ids)

	for i := 0; i < 5; i++ {
		need := numStories - len(stories)
		to := from + need

		if to > idsLength {
			to = idsLength
		}

		stories = append(stories, getStories(ids[from:to])...)
		from = to
		if len(stories) > numStories || from >= idsLength {
			break
		}
	}

	if len(stories) >= numStories {
		return stories[:numStories], nil
	}

	return stories, nil
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem client.Item) item {
	ret := item{Item: hnItem}
	fullUrl, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(fullUrl.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	client.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}
