package main

import (
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
  tmpPtr := flag.String("template", "template.html", "path to template file")
  jsonFlag := flag.Bool("json", false, "output JSON instead")
  flag.Parse()

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

  if *jsonFlag == true {
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
