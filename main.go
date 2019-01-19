package main

import (
  "fmt"
  "os"
  "log"
  "flag"
  "encoding/json"
  "github.com/pilgreen/slush/article"
)

func main() {
  var out []byte
  var err error

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

  args := flag.Args()
  list := make([]article.Article, len(args))
  ch := make(chan article.Article, len(list))

  for i, url := range args {
    a, err := article.Fetch(url, *bodyPtr)
    check(err)

    ch <- a
    list[i] = <-ch
  }

  if len(args) == 1 {
    out, err = json.Marshal(list[0])
    check(err)
  } else {
    out, err = json.Marshal(list)
    check(err)
  }

  os.Stdout.Write(out)
}

func check(e error) {
  if e != nil {
    log.Fatal(e)
  }
}

