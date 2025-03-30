package providers

import (
  "net/url"
  "strings"

  "github.com/ItsMeSamey/subdl_go/common"
  "github.com/ItsMeSamey/subdl_go/dlutils"

  "github.com/ItsMeSamey/go_utils"
  "github.com/PuerkitoBio/goquery"
)

const subdl_com_api = "https://api.subdl.com"
const subdl_com_site = "https://subdl.com"

type SubdlSuggestion struct {
  Type         string `json:"type"`
  Name         string `json:"name"`
  PosterURL    string `json:"poster_url"`
  Year         int    `json:"year"`
  Link         string `json:"link"`
  OriginalName string `json:"original_name"`
}

type SubdlSuggestionResult struct {
  Results []SubdlSuggestion `json:"results"`
}

func FetchSubdlCom(query string, options common.SearchOptions) (retval []common.MovieListEntry, err error) {
  searchURL := subdl_com_api + "/auto?query=" + url.QueryEscape(query)
  searchResult, status, err := dlutils.FetchJson[SubdlSuggestionResult](searchURL, dlutils.FetchInit{})
  if err != nil {
    return nil, utils.WithStack(err)
  }
  if err = status.Error(); err != nil {
    return
  }

  for _, e := range searchResult.Results {
    retval = append(retval, &SubdlMovieLink{
      data: common.MovieListData{
        Title:   e.Name,
        Options: options,
      },
      link: e.Link,
    })
  }

  return retval, nil
}

type SubdlMovieLink struct {
  data common.MovieListData
  link string
}

func (m *SubdlMovieLink) Data() *common.MovieListData { return &m.data }

func (m *SubdlMovieLink) ToSubtitleLinks() (retval []common.SubtitleListEntry, err error) {
  root, status, err := dlutils.FetchHtml(subdl_com_site + m.link, dlutils.FetchInit{})
  if err != nil { return nil, utils.WithStack(err) }
  if err = status.Error(); err != nil { return }

  root.Find(`div[class="flex flex-col mt-4 select-none"]`).Each(func(i int, section *goquery.Selection) {
    langElem := section.Children().First().Children().First().Children().First()
    lang := strings.ToLower(strings.TrimSpace(langElem.Text()))
    if lang == "" { return }
    if m.data.Options.Language != "" && lang != string(m.data.Options.Language) { return }

    section.Find("li").Each(func(j int, s *goquery.Selection) {
      linkElem := s.Children().Last().Children().Last()
      linkPath, ok := linkElem.Attr("href")
      if !ok || linkPath == "" {
        return
      }

      filenameElem := s.Children().First().Children().First()
      filename := strings.TrimSpace(filenameElem.Text())

      retval = append(retval, &SubdlSubtitleLink{
        data: common.SubtitleListData{
          Parent:   m,
          Filename: filename,
          Language: lang,
        },
        downloadPath: linkPath,
      })
    })
  })

  return retval, nil
}

type SubdlSubtitleLink struct {
  data         common.SubtitleListData
  downloadPath string
}
func (s *SubdlSubtitleLink) Data() *common.SubtitleListData { return &s.data }
func (s *SubdlSubtitleLink) IsZip() bool { return true }
func (s *SubdlSubtitleLink) DownloadLink() (string, error) {
  return s.downloadPath, nil
}

