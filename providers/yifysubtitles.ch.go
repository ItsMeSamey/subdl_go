package providers

import (
  "net/url"
  "strings"

  "subtitle_downloader/common"
  "subtitle_downloader/dlutils"

  "github.com/ItsMeSamey/go_fuzzy"
  "github.com/ItsMeSamey/go_utils"
  "github.com/PuerkitoBio/goquery"
)

const yifysubtitles_ch = "https://yifysubtitles.ch"

type YifySubtitlesSuggestion struct {
  Movie string `json:"movie"`
  Imdb  string `json:"imdb"`
}

func FetchYifySubtitlesCh(query string, options common.SearchOptions) (retval []common.MovieListEntry, err error) {
  searchURL := yifysubtitles_ch + "/ajax/search/?mov=" + url.QueryEscape(query)
  searchResult, status, err := dlutils.FetchJson[[]YifySubtitlesSuggestion](searchURL, dlutils.FetchInit{})
  if err != nil { return nil, utils.WithStack(err) }
  if err = status.Error(); err != nil { return }

  for _, e := range searchResult {
    retval = append(retval, &YifySubtitlesMovieLink{
      data: common.MovieListData{
        Title:   e.Movie,
        Options: options,
      },
      link: yifysubtitles_ch + "/movie-imdb/" + e.Imdb,
    })
  }

  return retval, nil
}

type YifySubtitlesMovieLink struct {
  data common.MovieListData
  link string
}
func (m *YifySubtitlesMovieLink) Data() *common.MovieListData { return &m.data }
func (m *YifySubtitlesMovieLink) ToSubtitleLinks() (retval []common.SubtitleListEntry, err error) {
  root, status, err := dlutils.FetchHtml(m.link, dlutils.FetchInit{})
  if err != nil { return nil, utils.WithStack(err) }
  if err = status.Error(); err != nil { return }

  table := root.Find("table").First()
  table.Find("tbody").Children().Each(func(i int, s *goquery.Selection) {
    linkElem := s.Find("a > span[class='text-muted']").Parent()
    filenameWithNewline := strings.TrimSpace(linkElem.Text())
    filenameParts := strings.SplitN(filenameWithNewline, "\n", 2)
    var filename string
    if len(filenameParts) > 0 {
      filename = filenameParts[0]
      if index := strings.Index(filename, " "); index != -1 {
        filename = filename[index+1:] + ".zip"
      } else {
        filename += ".zip"
      }
    }

    href, ok := linkElem.Attr("href")
    if !ok || href == "" {
      return
    }

    langElem := s.Find("span[class='sub-lang']").First()
    lang := strings.TrimSpace(langElem.Text())

    subtitleLink := &YifySubtitlesSubtitleLink{
      data: common.SubtitleListData{
        Parent:   m,
        Filename: filename,
        Language: lang,
      },
      link: href,
    }
    retval = append(retval, subtitleLink)
  })

  if m.data.Options.Language != "" {
    languageName, ok := common.LanguageNameMap[m.data.Options.Language]
    if !ok { return }
    m.data.Options.Sorter.SortAny(fuzzy.ToSwapper(retval, func (s common.SubtitleListEntry) string { return s.Data().Filename }), languageName)
  }

  return
}

type YifySubtitlesSubtitleLink struct {
  data common.SubtitleListData
  link string
}

func (s *YifySubtitlesSubtitleLink) Data() *common.SubtitleListData { return &s.data }
func (s *YifySubtitlesSubtitleLink) IsZip() bool { return true }
func (s *YifySubtitlesSubtitleLink) DownloadLink() (string, error) {
  return strings.Replace(yifysubtitles_ch + s.link + ".zip", "/subtitles/", "/subtitle/", 1), nil
}

