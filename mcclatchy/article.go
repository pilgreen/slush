package mcclatchy

import (
  "regexp"
  "time"
  "encoding/json"
  "net/http"
  "strconv"

  "image"
  _ "image/jpeg"

  "github.com/PuerkitoBio/goquery"
  "github.com/tdewolff/minify"
  "github.com/tdewolff/minify/html"
)

type Article struct {
  Meta MetaData `json:"meta"`
  Photo Photo `json:"photo"`
  Video Video `json:"video"`
  Body string `json:"body,omitempty"`
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
  Sizes []PhotoSource `json:"sizes,omitempty"`
}

type PhotoSource struct {
  Url string `json:"url"`
}

type Video struct {
  Id string `json:"id"`
  Url string `json:"url"`
  Duration string `json:"duration"`
  Image string `json:"image"`
  Videographer string `json:"videographer"`
  Credit string `json:"credit"`
  Title string `json:"description"`
  Caption string `json:"displayDescription"`
}

/**
 * Gets Meta information from a goquery document
 */

func ParseArticle(doc *goquery.Document, body bool) (article Article) {
  article.Meta = ArticleMeta(doc)
  article.Photo = ArticlePhoto(doc)
  article.Video = ArticleVideo(doc)
  if(body) {
    article.Body = ArticleBody(doc)
  }
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
    photo.Width, photo.Height = photo.Dimensions()
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
      photo.Sizes = append(photo.Sizes, s)
    })
  }

  return
}

/**
 * Retrieves the lead video asset
 */

 func ArticleVideo(doc *goquery.Document) (video Video) {
   // This is totally hacked in Vim and is not stable
   r := regexp.MustCompile(`(?m)return video\((\{.*\}\]\}\]\}), `)

   wrapper := doc.Find(".video-media[data-id]")
   match := r.FindStringSubmatch(wrapper.Text())
   json.Unmarshal([]byte(match[1]), &video)
   return
 }

/**
 * Retrieves the body of the story
 */

func ArticleBody(doc *goquery.Document) string {
  raw, _ := doc.Find("[id^=content-body]").Html()

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

/**
 * Gets the url for the smallest image size
 */

func (p Photo) Smallest() PhotoSource {
  match := p.Sizes[0]
  for i := 1; i < len(p.Sizes); i++ {
    s := p.Sizes[i]
    if s.Width() < match.Width() {
      match = s
    }
  }
  return match
}

/**
 * Finds the width of a photo via the url
 */

func (ps PhotoSource) Width() (w int) {
  r, _ := regexp.Compile(`/FREE_([0-9]+)/`)
  match := r.FindStringSubmatch(ps.Url)
  w, _ = strconv.Atoi(match[1])
  return
}

