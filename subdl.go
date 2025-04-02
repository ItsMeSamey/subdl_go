package subdl

import (
  "github.com/ItsMeSamey/go_fuzzy"

  "github.com/ItsMeSamey/subdl_go/common"
  "github.com/ItsMeSamey/subdl_go/dlutils"
)

type errNoMovies struct{}
func (e errNoMovies) Error() string { return "No movies found." }
var ErrNoMovies = errNoMovies{}

type errNoSubtitles struct{}
func (e errNoSubtitles) Error() string { return "No subtitles found." }
var ErrNoSubtitles = errNoSubtitles{}

func Download(
  query string,
  searchOptions common.SearchOptions,
  provider func(query string, options common.SearchOptions) ([]common.MovieListEntry, error),
  downloadOptions common.DownloadOptions,
) (retval common.DownloadedSubtitle, err error) {
  movies, err := provider(query, searchOptions)
  if err != nil { return }
  if len(movies) == 0 {
    err = ErrNoMovies
    return 
  }
  if downloadOptions.MovieListQuery != "" {
    downloadOptions.MovieListSorter.SortAny(
      fuzzy.ToSwapper(movies, func(m common.MovieListEntry) string { return m.Data().Title }),
      downloadOptions.MovieListQuery,
    )
  }

  subtitles, err := movies[0].ToSubtitleLinks()
  if err != nil { return }
  if len(subtitles) == 0 {
    err = ErrNoSubtitles
    return 
  }
  if downloadOptions.SubtitleListQuery != "" {
    downloadOptions.SubtitleListSorter.SortAny(
      fuzzy.ToSwapper(subtitles, func(s common.SubtitleListEntry) string { return s.Data().Filename }),
      downloadOptions.SubtitleListQuery,
    )
  }

  retval, err = dlutils.DownloadSubtitleEntry(subtitles[0])
  if err != nil { return }

  return
}

