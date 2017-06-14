package main

import (
  "fmt"
  "os"
  "encoding/json"
  "sync"

  "github.com/PuerkitoBio/goquery"
  "github.com/pilgreen/slush/mcclatchy"
)

var wg sync.WaitGroup
var list []mcclatchy.Article

func main() {
  if len(os.Args) <= 1 {
    fmt.Printf("Usage: %s [SECTION URL]\n", os.Args[0])
    os.Exit(1)
  }

  var links []string
  doc, _ := goquery.NewDocument(os.Args[1])
  links = mcclatchy.SectionLinks(doc)

  list = make([]mcclatchy.Article, len(links))
  for i, url := range links {
    if len(url) > 0 {
      wg.Add(1)
      go fetch(url, i)
    }
  }

  wg.Wait()

  b, _ := json.Marshal(list)
  os.Stdout.Write(b)
}

func fetch(url string, slot int) {
  doc, _ := goquery.NewDocument(url)
  list[slot] = mcclatchy.ParseArticle(doc)
  wg.Done()
}

func filter(vs []string, f func(string) bool) []string {
  vsf := make([]string, 0)
  for _, v := range vs {
    if f(v) {
      vsf = append(vsf, v)
    }
  }
  return vsf
}
