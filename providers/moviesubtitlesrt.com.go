package providers

import (
  "net/url"
  "strings"

  "github.com/ItsMeSamey/subdl_go/common"
  "github.com/ItsMeSamey/subdl_go/dlutils"

  "github.com/ItsMeSamey/go_utils"
  "github.com/PuerkitoBio/goquery"
)

const moviesubtitlesrt_com = "https://moviesubtitlesrt.com"

func FetchMoviesubtitlesrtCom(query string, options common.SearchOptions) (retval []common.MovieListEntry, err error) {
  searchURL := moviesubtitlesrt_com + "/?s=" + url.QueryEscape(query)
  root, status, err := dlutils.FetchHtml(searchURL, dlutils.FetchInit{})
  if err != nil { return nil, utils.WithStack(err) }
  if err = status.Error(); err != nil { return }

  root.Find(`div[class="inside-article"] > header > h2 > a`).Each(func(i int, s *goquery.Selection) {
    link, ok := s.Attr("href")
    if !ok || link == "" { return }
    title := strings.TrimSpace(s.Text())
    if title == "" { return }

    retval = append(retval, &MoviesubtitlesrtMovieLink{
      data: common.MovieListData{
        Title:   title,
        Options: options,
      },
      link: link,
    })
  })

  return retval, nil
}

type MoviesubtitlesrtMovieLink struct {
  data common.MovieListData
  link string
}
func (m *MoviesubtitlesrtMovieLink) Data() *common.MovieListData { return &m.data }
func (m *MoviesubtitlesrtMovieLink) ToSubtitleLinks() (retval []common.SubtitleListEntry, err error) {
  root, status, err := dlutils.FetchHtml(m.link, dlutils.FetchInit{})
  if err != nil { return nil, utils.WithStack(err) }
  if err = status.Error(); err != nil { return }

  langName := ""
  langNode := root.Find("tbody > tr:nth-child(2) > td:last-child").First()
  if langNode.Length() > 0 { langName = strings.TrimSpace(langNode.Text()) }

  langCode := common.FindLanguageCode(langName, m.data.Options.Sorter)

  downloadLinkNode := root.Find("center > a").First()
  downloadLink, ok := downloadLinkNode.Attr("href")
  if !ok || downloadLink == "" { // If the primary link isn't found, skip adding this entry
    return retval, nil // Return empty slice, not an error
  }

  retval = append(retval, &MoviesubtitlesrtSubtitleLink{
    data: common.SubtitleListData{
      Parent:   m,
      Filename: m.data.Title,
      Language: string(langCode),
    },
    link: downloadLink,
  })

  targetLang := string(m.Data().Options.Language)
  if targetLang != "" {
    filteredRetval := []common.SubtitleListEntry{}
    for _, entry := range retval {
      if strings.ToLower(entry.Data().Language) == targetLang {
        filteredRetval = append(filteredRetval, entry)
      }
    }
    retval = filteredRetval
  }

  return retval, nil
}

type MoviesubtitlesrtSubtitleLink struct {
  data common.SubtitleListData
  link string
}
func (s *MoviesubtitlesrtSubtitleLink) Data() *common.SubtitleListData { return &s.data }
func (s *MoviesubtitlesrtSubtitleLink) IsZip() bool { return true }
func (s *MoviesubtitlesrtSubtitleLink) DownloadLink() (string, error) {
  return s.link, nil
}

