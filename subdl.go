package subdl

import (
  "github.com/ItsMeSamey/go_fuzzy/heuristics"
  "github.com/ItsMeSamey/go_fuzzy/transformers"
  "github.com/ItsMeSamey/go_utils"
  "github.com/ItsMeSamey/subdl_go/common"
  "github.com/ItsMeSamey/subdl_go/dlutils"
)


func DownloadSubtitles(
  query string,
  options common.SearchOptions,
  provider func(query string, options common.SearchOptions) ([]common.MovieListEntry, error),
) (retval common.DownloadedSubtitle, err error) {
  if options.Sorter.ScoreFn == nil {
    options.Sorter.ScoreFn = heuristics.Wrap[float32](heuristics.FrequencySimilarity)
    options.Sorter.Transformer = transformers.Lowercase()
  }

  movies, err := provider(query, options)
  if err = utils.WithStack(err); err != nil { return }
  if len(movies) == 0 { return }
  // TODO: maybe sort movies

  subtitles, err := movies[0].ToSubtitleLinks()
  if err = utils.WithStack(err); err != nil { return }
  if len(subtitles) == 0 { return }
  // TODO: maybe sort subtitles

  retval, err = dlutils.DownloadSubtitleEntry(subtitles[0])
  if err = utils.WithStack(err); err != nil { return }

  return
}

