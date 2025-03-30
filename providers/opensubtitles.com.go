package providers

import (
  "fmt"
  "net/url"
  "regexp"
  "strings"

  "github.com/ItsMeSamey/subdl_go/common"
  "github.com/ItsMeSamey/subdl_go/dlutils"

  "github.com/ItsMeSamey/go_utils"
  "github.com/PuerkitoBio/goquery"
)

const opensubtitles_com = "https://www.opensubtitles.com"

type OpenSubtitlesSuggestion struct {
  Title          string  `json:"title"`
  Year           string  `json:"year"`
  ID             string  `json:"id"`
  Poster         string  `json:"poster"`
  Rating         float64 `json:"rating"`
  SubtitlesCount int     `json:"subtitles_count"`
  Type           string  `json:"type"`
  Path           string  `json:"path"`
}

type OpenSubtitlesSubtitlesResponse struct {
  Data [][]string `json:"data"`
}

func FetchOpenSubtitlesCom(query string, options common.SearchOptions) (retval []common.MovieListEntry, err error) {
  languageID := "en"
  if options.Language != "" {
    languageID = string(options.Language)
  }

  autocompleteURL := opensubtitles_com + "/en/en/search/autocomplete/" + url.QueryEscape(query) + ".json"

  suggestions, status, err := dlutils.FetchJson[[]OpenSubtitlesSuggestion](autocompleteURL, dlutils.FetchInit{})
  if err != nil { return nil, utils.WithStack(err) }
  if err = status.Error(); err != nil { return }

  for _, e := range suggestions {
    subtitlesURL := opensubtitles_com + "/" + languageID + strings.Replace(strings.Replace(e.Path, "current_locale", languageID, 1), "movies", "features", 1) + "/subtitles.json"
    retval = append(retval, &OpenSubtitlesMovieLink{
      data: common.MovieListData{
        Title:   e.Title,
        Options: options,
      },
      link: subtitlesURL,
    })
  }

  return retval, nil
}

type OpenSubtitlesMovieLink struct {
  data common.MovieListData
  link string
}

func (m *OpenSubtitlesMovieLink) Data() *common.MovieListData { return &m.data }

func (m *OpenSubtitlesMovieLink) ToSubtitleLinks() (retval []common.SubtitleListEntry, err error) {
  subtitlesResponse, status, err := dlutils.FetchJson[OpenSubtitlesSubtitlesResponse](m.link, dlutils.FetchInit{})
  if err != nil { return nil, utils.WithStack(err) }
  if err = status.Error(); err != nil { return }

  for _, s := range subtitlesResponse.Data {
    if len(s) == 0 { continue }

    var filename, language, initialLink string

    if len(s) > 2 {
      doc, _ := goquery.NewDocumentFromReader(strings.NewReader(s[2]))
      filename = strings.Split(strings.TrimSpace(doc.Text()), "\n")[0]
    }

    if len(s) > 1 {
      doc, _ := goquery.NewDocumentFromReader(strings.NewReader(s[1]))
      language, _ = doc.Children().First().Attr("title")
      language = strings.TrimSpace(language)
    }

    doc, _ := goquery.NewDocumentFromReader(strings.NewReader(s[len(s)-1]))
    linkNode := doc.Find(`a[data-remote="true"]`).First()
    initialLink, _ = linkNode.Attr("href")

    if initialLink == "" { continue }

    retval = append(retval, &OpenSubtitlesSubtitleLink{
      data: common.SubtitleListData{
        Parent:   m,
        Filename: filename,
        Language: language,
      },
      initialLink: initialLink,
    })
  }

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

type OpenSubtitlesSubtitleLink struct {
  data common.SubtitleListData
  initialLink string
  finalLink   string
}

func (s *OpenSubtitlesSubtitleLink) Data() *common.SubtitleListData { return &s.data }
func (s *OpenSubtitlesSubtitleLink) IsZip() bool { return false }

var fileDownloadRegex = regexp.MustCompile(`file_download\('([^']*)','([^']*)'`)
func (s *OpenSubtitlesSubtitleLink) DownloadLink() (link string, err error) {
  if s.finalLink != "" { return s.finalLink, nil }

  fetchURL := opensubtitles_com + s.initialLink
  bodyBytes, status, err := dlutils.FetchText(fetchURL, dlutils.FetchInit{
    Headers: []dlutils.Header{
      {Key: "x-csrf-token", Value: "SZHfvYUiNV3uhpKkRPfQPcfhqtrdJVw9hCwxAc+XknB5Wsct+7gZOHlrwJqWElrevrWoZlReTBeJmSPPIVWmzw=="},
      {Key: "x-requested-with", Value: "XMLHttpRequest"},
    },
  })
  if err = utils.WithStack(err); err != nil { return }
  if err = utils.WithStack(status.Error()); err != nil { return }

  body := string(bodyBytes)
  matches := fileDownloadRegex.FindStringSubmatch(body)

  if len(matches) < 3 { return "", utils.WithStack(fmt.Errorf("could not extract download link from script at %s. Body: %s", fetchURL, body)) }

  filename := matches[1]
  downloadLink := matches[2]

  if filename != "" { s.data.Filename = filename }
  s.finalLink = downloadLink

  return s.finalLink, nil
}

