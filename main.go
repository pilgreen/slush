package main

import (
  "fmt"
  "os"
  "log"
  "flag"
  "encoding/json"
  "text/template"
  "sync"

  "github.com/PuerkitoBio/goquery"
  "github.com/pilgreen/slush/mcclatchy"
)

var wg sync.WaitGroup
var list []mcclatchy.Article

func main() {
  var links []string

  tmpPtr := flag.String("template", "", "path to golang template file")
  secPtr := flag.Bool("section", false, "pull all stories for a section")
  bodyPtr := flag.Bool("body", false, "include the body html")
  flag.Parse()

  flag.Usage = func() {
    fmt.Fprintln(os.Stderr, "Usage: slush [options] url1 url2 ... urlN")
    flag.PrintDefaults()
  }

  if len(os.Args) < 2 {
    flag.Usage()
    os.Exit(1)
  }

  links = flag.Args()
  if *secPtr {
    doc, _ := goquery.NewDocument(links[0])
    links = mcclatchy.SectionLinks(doc)
  }

  list = make([]mcclatchy.Article, len(links))
  for i, url := range links {
    wg.Add(1)
    go fetch(url, i, *bodyPtr)
  }

  wg.Wait()

  if len(*tmpPtr) == 0 {
    enc := json.NewEncoder(os.Stdout)
    enc.SetEscapeHTML(false)
    enc.Encode(list)
  } else {
    tpl := template.Must(template.New("main").ParseGlob(*tmpPtr))
    tpl.ExecuteTemplate(os.Stdout, *tmpPtr, list)
  }
}

func fetch(url string, slot int, body bool) {
  defer wg.Done()
  doc, _ := goquery.NewDocument(url)
  list[slot] = mcclatchy.ParseArticle(doc, body)
}

func check(e error) {
  if e != nil {
    log.Fatal(e)
  }
}

