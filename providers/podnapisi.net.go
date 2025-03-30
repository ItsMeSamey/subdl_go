package providers

import (
  "net/url"
  "strings"

  "subtitle_downloader/common"
  "subtitle_downloader/dlutils"

  "github.com/ItsMeSamey/go_utils"
  "github.com/PuerkitoBio/goquery"
)

const podnapisi_net = "https://www.podnapisi.net"

type PodnapisiPoster struct {
  Inline string `json:"inline"`
  Normal string `json:"normal"`
  Small  string `json:"small"`
  Title  string `json:"title"`
}

type PodnapisiSuggestion struct {
  Aliases   []string        `json:"aliases"`
  ID        string          `json:"id"`
  Posters   PodnapisiPoster `json:"posters"`
  Providers []string        `json:"providers"`
  Slug      string          `json:"slug"`
  Title     string          `json:"title"`
  Type      string          `json:"type"`
  Year      int             `json:"year"`
}

type PodnapisiSuggestionResult struct {
  Aggs   map[string]any         `json:"aggs"`
  Data   []PodnapisiSuggestion  `json:"data"`
  Status string                 `json:"status"`
}

func FetchPodnapisiNet(query string, options common.SearchOptions) (retval []common.MovieListEntry, err error) {
  searchURL := podnapisi_net + "/moviedb/search/?keywords=" + url.QueryEscape(query)
  searchResult, status, err := dlutils.FetchJson[PodnapisiSuggestionResult](searchURL, dlutils.FetchInit{
    Headers: []dlutils.Header{
      {Key: "X-Requested-With", Value: "XMLHttpRequest"},
      {Key: "Accept", Value: "*/*"},
    },
  })
  if err != nil { return nil, utils.WithStack(err) }
  if err = status.Error(); err != nil { return }

  for _, e := range searchResult.Data {
    retval = append(retval, &PodnapisiMovieLink{
      data: common.MovieListData{
        Title:   e.Title,
        Options: options,
      },
      movieID: e.ID,
    })
  }

  return retval, nil
}

type PodnapisiMovieLink struct {
  data    common.MovieListData
  movieID string
}

func (m *PodnapisiMovieLink) Data() *common.MovieListData { return &m.data }

func (m *PodnapisiMovieLink) ToSubtitleLinks() (retval []common.SubtitleListEntry, err error) {
  subtitlesSearchURL := podnapisi_net + "/subtitles/search/" + m.movieID

  root, status, err := dlutils.FetchHtml(subtitlesSearchURL, dlutils.FetchInit{
    Headers: []dlutils.Header{{Key: "Accept", Value: "*/*"}},
  })
  if err != nil { return nil, utils.WithStack(err) }
  if err = status.Error(); err != nil { return }

  root.Find("tbody").First().Find("tr").Each(func(i int, tr *goquery.Selection) {
    lang := strings.TrimSpace(tr.Find("abbr").First().Text())
    if lang == "" { return }

    targetLang := string(m.data.Options.Language)
    if targetLang != "" && strings.ToLower(lang) != strings.ToLower(targetLang) {
      return
    }

    linkPath, ok := tr.Find(`a[rel="nofollow"]`).First().Attr("href")
    if !ok || linkPath == "" { return }

    filename := strings.TrimSpace(tr.Find(`span.release`).First().Text())

    retval = append(retval, &PodnapisiSubtitleLink{
      data: common.SubtitleListData{
        Parent:   m,
        Filename: filename,
        Language: lang,
      },
      linkPath: linkPath,
    })
  })

  return retval, nil
}

type PodnapisiSubtitleLink struct {
  data     common.SubtitleListData
  linkPath string
}
func (s *PodnapisiSubtitleLink) Data() *common.SubtitleListData { return &s.data }
func (s *PodnapisiSubtitleLink) IsZip() bool                     { return true }
func (s *PodnapisiSubtitleLink) DownloadLink() (string, error) {
  return podnapisi_net + s.linkPath, nil
}

