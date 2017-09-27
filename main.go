package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	hnBaseUrl  = "https://hacker-news.firebaseio.com/v0"
	numStories = 30
)

// Error Constants
const (
	ErrTooManyStories constErr = "Too many stories requested"
)

type constErr string

func (c constErr) Error() string {
	return string(c)
}

func main() {
	port := flag.Int("port", 3000, "the port to start the server on")
	indexTmpl := template.Must(template.ParseFiles("index.gohtml"))
	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Starting server at %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, index(indexTmpl)))
}

type indexData struct {
	Stories    []story
	RenderTime time.Duration
}

func index(tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		stories, err := getTopStories(numStories)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := indexData{
			Stories:    stories,
			RenderTime: time.Now().Sub(start),
		}
		w.Header().Set("Content-Type", "text/html")
		tmpl.Execute(w, data)
	}
}

func storiesUrl() string {
	return hnBaseUrl + "/topstories.json"
}

func storyUrl(id int) string {
	return fmt.Sprintf("%s/item/%d.json", hnBaseUrl, id)
}

type story struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	URL   string `json:"url"`
}

// Hostname will attempt to parse the URL and return just
// the domain. Eg "gophercises.com" from the URL
// https://www.gophercises.com/exercises/quiethn
// It is far from perfect and will miss all subdomains
// aside from www
func (s story) Domain() string {
	u, err := url.Parse(s.URL)
	if err != nil {
		return "unknown"
	}
	return strings.TrimPrefix(u.Hostname(), "www.")
}

func (s story) isDiscussion() bool {
	return s.Text != ""
}

func getTopStories(n int) ([]story, error) {
	var ids []int
	err := getJson(storiesUrl(), &ids)
	if err != nil {
		return nil, err
	}
	if n > len(ids) {
		return nil, ErrTooManyStories
	}
	var stories []story
	for _, id := range ids {
		story, err := getStory(id)
		if err != nil {
			continue
		}
		if !story.isDiscussion() {
			stories = append(stories, story)
			if len(stories) >= n {
				break
			}
		}
	}
	return stories, nil
}

func getStory(id int) (story, error) {
	var s story
	err := getJson(storyUrl(id), &s)
	if err != nil {
		return s, err
	}
	return s, nil
}

func getJson(url string, out interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	dec := json.NewDecoder(r.Body)
	err = dec.Decode(out)
	if err != nil {
		return err
	}
	return nil
}
