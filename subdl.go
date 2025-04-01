package subdl

import (
  "github.com/ItsMeSamey/go_fuzzy/heuristics"
  "github.com/ItsMeSamey/go_fuzzy/transformers"
  "github.com/ItsMeSamey/go_utils"
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
  options common.SearchOptions,
  provider func(query string, options common.SearchOptions) ([]common.MovieListEntry, error),
) (retval common.DownloadedSubtitle, err error) {
  if options.Sorter.ScoreFn == nil {
    options.Sorter.ScoreFn = heuristics.Wrap[float32](heuristics.FrequencySimilarity)
    options.Sorter.Transformer = transformers.Lowercase()
  }

  movies, err := provider(query, options)
  // for _, m := range movies { println("Movie: ", m.Data().Title) }
  if err = utils.WithStack(err); err != nil { return }
  if len(movies) == 0 {
    err = ErrNoMovies
    return 
  }
  // TODO: maybe sort movies

  subtitles, err := movies[0].ToSubtitleLinks()
  // for _, s := range subtitles { println("Subtitle: ", s.Data().Filename) }
  if err != nil { return }
  if len(subtitles) == 0 {
    err = ErrNoSubtitles
    return 
  }
  // TODO: maybe sort subtitles

  retval, err = dlutils.DownloadSubtitleEntry(subtitles[0])
  if err != nil { return }

  return
}

