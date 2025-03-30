package providers

import (
  "net/url"
  "strings"

  "github.com/ItsMeSamey/subdl_go/common"
  "github.com/ItsMeSamey/subdl_go/dlutils"

  "github.com/ItsMeSamey/go_utils"
  "github.com/PuerkitoBio/goquery"
)

const moviesubtitles_org = "https://www.moviesubtitles.org"

func FetchMovieSubtitlesOrg(query string, options common.SearchOptions) (retval []common.MovieListEntry, err error) {
  formData := url.Values{}
  formData.Set("q", query)

  root, status, err := dlutils.FetchHtml(moviesubtitles_org+"/search.php", dlutils.FetchInit{
    Method: "POST",
    Headers: []dlutils.Header{{Key: "Content-Type", Value: "application/x-www-form-urlencoded"}},
    Body:   []byte(formData.Encode()),
  })

  _ = status // Always returns 500 for some reason.
  if err = utils.WithStack(err); err != nil { return }

  root.Find(`div[style="width:500px"] > a`).Each(func(i int, s *goquery.Selection) {
    link, ok := s.Attr("href")
    if !ok || link == ""{ return }
    title := strings.TrimSpace(s.Text())
    if title == "" { return }

    retval = append(retval, &MoviesSubtitlesMovieLink{
      data: common.MovieListData{
        Title:   title,
        Options: options,
      },
      link: link,
    })
  })

  return retval, nil
}

type MoviesSubtitlesMovieLink struct {
  data common.MovieListData
  link   string
}
func (m *MoviesSubtitlesMovieLink) Data() *common.MovieListData { return &m.data }
func (m *MoviesSubtitlesMovieLink) ToSubtitleLinks() (retval []common.SubtitleListEntry, err error) {
  root, status, err := dlutils.FetchHtml(moviesubtitles_org + m.link, dlutils.FetchInit{})
  if err != nil { return nil, utils.WithStack(err) }
  if err = status.Error(); err != nil { return }

  root.Find(`div[style="margin-bottom:0.5em; padding:3px;"]`).Each(func(i int, s *goquery.Selection) {
    link, ok := s.Find("a").Attr("href")
    if !ok || link == "" { return } // Skip if href is missing

    languageSrc, ok := s.Find("img").Attr("src")
    if !ok || languageSrc == "" { return }

    languageParts := strings.Split(languageSrc, "/")
    languageNameParts := strings.Split(languageParts[len(languageParts)-1], ".")

    retval = append(retval, &MoviesSubtitlesSubtitleLink{
      data: common.SubtitleListData{
        Parent:   m,
        Filename: strings.Split(s.Text(), "\n")[0],
        Language: languageNameParts[0],
      },
      link: strings.Replace(link, "subtitle", "download", 1),
    })
  })

  lang := string(m.Data().Options.Language)
  if lang != "" {
    for i := 0; i < len(retval); i += 1 {
      if strings.ToLower(retval[i].Data().Language) != lang {
        retval[i], retval[len(retval)-1] = retval[len(retval)-1], retval[i]
        retval = retval[:len(retval)-1]
      }
    }
  }

  return
}

type MoviesSubtitlesSubtitleLink struct {
  data common.SubtitleListData
  link string
}
func (s *MoviesSubtitlesSubtitleLink) Data() *common.SubtitleListData { return &s.data }
func (s *MoviesSubtitlesSubtitleLink) IsZip() bool { return true }
func (s *MoviesSubtitlesSubtitleLink) DownloadLink() (string, error) {
  return moviesubtitles_org + s.link, nil
}

