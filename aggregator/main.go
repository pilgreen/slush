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

type Article struct {
  Meta mcclatchy.MetaData
  Photo mcclatchy.Photo
}

var wg sync.WaitGroup
var list []Article

func main() {
  tmpPtr := flag.String("template", "", "path to template file")
  flag.Parse()

  flag.Usage = func() {
    fmt.Fprintf(os.Stderr, "Usage: %s [options] url1 url2 ... urlN\n", os.Args[0])
    flag.PrintDefaults()
  }

  /**
   * Section Flag
   */

  links := flag.Args()

  if len(links) == 0 {
    flag.Usage()
    os.Exit(1)
  }

  list = make([]Article, len(links))

  for i, url := range links {
    wg.Add(1)
    go fetch(url, i)
  }

  wg.Wait()

  if len(*tmpPtr) == 0 {
    b, _ := json.Marshal(list)
    os.Stdout.Write(b)
  } else {
    tmpl, err := template.ParseFiles(*tmpPtr)
    check(err)

    for _, story := range list {
      tmpl.Execute(os.Stdout, story)
    }
  }
}

func fetch(url string, slot int) {
  defer wg.Done()
  doc, _ := goquery.NewDocument(url)
  article := Article{mcclatchy.ArticleMeta(doc), mcclatchy.ArticlePhoto(doc)}
  list[slot] = article
}

func check(e error) {
  if e != nil {
    log.Fatal(e)
  }
}

