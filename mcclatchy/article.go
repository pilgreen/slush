package mcclatchy

import (
  "regexp"
  "time"
  "net/http"
  "strconv"

  "image"
  _ "image/jpeg"

  "github.com/PuerkitoBio/goquery"
)

type Article struct {
  Meta MetaData `json:"meta"`
  Photo Photo `json:"photo"`
  Body string `json:"body"`
}

type MetaData struct {
  Url string `json:"url"`
  Title string `json:"title"`
  Author string `json:"author"`
  Summary string `json:"summary"`
  Published string `json:"published"`
  PubDate time.Time `json:"pubdate"`
}

type Photo struct {
  Photographer string `json:"photographer"`
  Caption string `json:"caption"`
  Credit string `json:"credit"`
  Url string `json:"url"`
  Width int `json:"width"`
  Height int `json:"height"`
  Sizes []PhotoSource `json:"sizes"`
}

type PhotoSource struct {
  Url string `json:"url"`
  Width int `json:"width"`
}

/**
 * Gets Meta information from a goquery document
 */

func ParseArticle(doc *goquery.Document) (article Article) {
  article.Meta = ArticleMeta(doc)
  article.Photo = ArticlePhoto(doc)
  article.Body = ArticleBody(doc)
  return
}

/**
 * Retrieves the article meta information
 */

func ArticleMeta(doc *goquery.Document) (meta MetaData) {
  meta.Url, _ = doc.Find("link[rel=canonical]").Attr("href")
  meta.Title, _ = doc.Find("meta[property='og:title']").Attr("content")
  meta.Author = doc.Find(".byline .ng_byline_name").Text()
  meta.Summary, _ = doc.Find("meta[name=description]").Attr("content")
  meta.Published = doc.Find("#story-header .published-date").Text()

  // parse date for the computers
  const dateForm = "January 02, 2006 3:04 PM"
  meta.PubDate, _ = time.Parse(dateForm, meta.Published)

  return
}

/**
 * Retrieves the lead asset
 */

func ArticlePhoto(doc *goquery.Document) (photo Photo) {
  photo.Url, _ = doc.Find("meta[property='og:image']").Attr("content")
  if len(photo.Url) > 0 {
    photo.Width, photo.Height = PhotoDimensions(photo.Url)
  }

  leadArt := doc.Find(".lead-item picture")
  if leadArt != nil {
    caption := doc.Find(".lead-item .caption")
    photo.Photographer = caption.Find(".photographer").Text()
    photo.Caption = caption.Find(".caption-text").Text()
    photo.Credit = caption.Find(".credits").Text()

    source := leadArt.Find("source")
    source.Each(func(i int, sel *goquery.Selection) {
      var s PhotoSource
      s.Url, _ = sel.Attr("srcset")
      s.Width = sourceWidth(s.Url)
      photo.Sizes = append(photo.Sizes, s)
    })
  }

  return
}

/**
 * Retrieves the body of the story
 */

func ArticleBody(doc *goquery.Document) string {
  html, _ := doc.Find("[id^=content-body]").Html()
  return html
}

/**
 * Fetches dimensions for a remote photo
 */

func PhotoDimensions(url string) (int, int) {
  resp, _ := http.Get(url)
  img, _, _ := image.DecodeConfig(resp.Body)
  return img.Width, img.Height
}

/**
 * Helpers
 */

func sourceWidth(url string) (w int) {
  r, _ := regexp.Compile("/FREE_([0-9]+)/")
  match := r.FindStringSubmatch(url)
  w, _ = strconv.Atoi(match[1])
  return
}

