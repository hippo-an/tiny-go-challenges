package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
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
		var client Client
		ids, err := client.TopItems()
		if err != nil {
			http.Error(w, "Failed to load top stories", http.StatusInternalServerError)
			return
		}
		var stories []item
		itemChan := make(chan *item, numStories)
		wg := sync.WaitGroup{}

		for _, id := range ids {
			wg.Add(1)

			go func() {
				defer wg.Done()

				hnItem, err := client.GetItem(id)
				if err != nil {
					itemChan <- nil
				}

				item := parseHNItem(hnItem)
				itemChan <- &item
			}()

			item := <-itemChan

			if item != nil && isStoryLink(*item) {
				stories = append(stories, *item)
				if len(stories) >= numStories {
					break
				}
			}
		}

		wg.Wait()

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

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}
	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}
