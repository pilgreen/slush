package article

import (
  "net/http"

  "image"
  _ "image/jpeg"

  "github.com/PuerkitoBio/goquery"
  "github.com/tdewolff/minify"
  "github.com/tdewolff/minify/html"
)

type Article struct {
  Url string `json:"url"`
  Title string `json:"title"`
  Summary string `json:"summary"`
  Published string `json:"published"`
  Body string `json:"body,omitempty"`
  Kicker Kicker `json:"kicker"`
  Photo Photo `json:"photo"`
}

type Photo struct {
  Url string `json:"url"`
  Width int `json:"width"`
  Height int `json:"height"`
}

type Kicker struct {
  Url string `json:"url"`
  Title string `json:"title"`
}

/**
 * Fetches an article
 */

func Fetch(url string, body bool) (article Article, err error) {
  doc, err := goquery.NewDocument(url)
  return ParseArticle(doc, body), err
}

/**
 * Gets Meta information from a goquery document
 */

func ParseArticle(doc *goquery.Document, body bool) (article Article) {
  article.Url, _ = doc.Find("link[rel=canonical]").Attr("href")
  article.Title, _ = doc.Find("meta[property='og:title']").Attr("content")
  article.Summary, _ = doc.Find("meta[name=description]").Attr("content")
  article.Published = doc.Find("#publish_date").First().Text()
  article.Kicker = ParseKicker(doc)
  article.Photo = ParsePhoto(doc)

  if(body) {
    article.Body = ParseBody(doc)
  }

  return
}

/**
 * Retrieves the lead asset
 */

func ParsePhoto(doc *goquery.Document) (photo Photo) {
  photo.Url, _ = doc.Find("meta[property='og:image']").Attr("content")
  if len(photo.Url) > 0 {
    photo.Width, photo.Height = photo.Dimensions()
  }

  return
}

/**
 * Retrieves the kicker info
 */

func ParseKicker(doc *goquery.Document) (kicker Kicker) {
  k := doc.Find(".lead-story-title .kicker-macro a")
  kicker.Url, _ = k.Attr("href")
  kicker.Title, _ = k.Attr("title")

  return
}

/**
 * Retrieves the body of the story
 */

func ParseBody(doc *goquery.Document) string {
  raw, _ := doc.Find(".content-body").Html()

  m := minify.New()
  m.Add("text/html", &html.Minifier{ KeepEndTags: true })

  min, _ := m.String("text/html", raw)
  return min
}

/**
 * Fetches dimensions for a remote photo
 */

func (p Photo) Dimensions() (int, int) {
  resp, _ := http.Get(p.Url)
  img, _, _ := image.DecodeConfig(resp.Body)
  return img.Width, img.Height
}
