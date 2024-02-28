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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		stories, err := getTopStories(numStories)
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
	})
}

func getTopStories(numStories int) ([]item, error) {
	var client client.Client
	ids, err := client.TopItems()
	if err != nil {
		return nil, errors.New("failed to load top stories")
	}

	type result struct {
		idx  int
		item item
		err  error
	}

	resultCh := make(chan result)

	for idx, id := range ids {

		go func(idx, id int) {
			hnItem, err := client.GetItem(id)
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

	for _, res := range results {
		if res.err != nil {
			continue
		}

		if isStoryLink(res.item) {
			if len(stories) >= numStories {
				break
			}

			stories = append(stories, res.item)
		}
	}

	return stories, nil
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem client.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
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
