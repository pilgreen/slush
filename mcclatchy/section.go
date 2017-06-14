package mcclatchy

import (
  "github.com/PuerkitoBio/goquery"
)

/**
 * Retrieves a list of published article urls
 */

func SectionLinks(doc *goquery.Document) (list []string) {
  mainstage := doc.Find("[id^=mainstage] article").First()
  mainstage.Each(func(i int, a *goquery.Selection) {
    link := titleLink(a)
    if len(link) > 0 {
      list = append(list, titleLink(a))
    }
  })

  wideRail := doc.Find("#story-list article")
  wideRail.Each(func(i int, a *goquery.Selection) {
    link := titleLink(a)
    if len(link) > 0 {
      list = append(list, titleLink(a))
    }
  })
  return
}

/**
 * Helpers
 */

func titleLink(article *goquery.Selection) (url string) {
  url, _ = article.Find(".title a").Attr("href")
  return
}
