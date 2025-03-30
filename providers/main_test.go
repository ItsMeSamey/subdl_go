package providers

import (
  "testing"

  "subtitle_downloader/common"
  "subtitle_downloader/dlutils"

  "github.com/ItsMeSamey/go_fuzzy/heuristics"
  "github.com/ItsMeSamey/go_fuzzy/transformers"
  "github.com/ItsMeSamey/go_utils"
)

func testProvider(t *testing.T, providerFn func(query string, options common.SearchOptions) ([]common.MovieListEntry, error)) {
  utils.SetErrorStackTrace(true)

  query := "The Matrix"
  options := common.SearchOptions{}

  if options.Sorter.ScoreFn == nil {
    options.Sorter.ScoreFn = heuristics.Wrap[float32](heuristics.FrequencySimilarity)
    options.Sorter.Transformer = transformers.Lowercase()
  }

  movies, err := providerFn(query, options)
  if err != nil {
    t.Fatalf("Error fetching movies: %v", err)
  }
  if len(movies) == 0 {
    t.Fatalf("No movies found.")
  }

  subtitles, err := movies[0].ToSubtitleLinks()
  if err != nil {
    t.Fatalf("Error fetching subtitles: %v", err)
  }
  if len(subtitles) == 0 {
    t.Fatalf("No subtitles found.")
  }

  result, err := dlutils.DownloadSubtitleEntry(subtitles[0])
  if err != nil {
    t.Fatalf("Error fetching subtitles: %v", err)
  }

  if len(result.Subtitles) == 0 {
    t.Fatalf("No subtitles found.")
  }

  for _, s := range result.Subtitles {
    println("Filename:", s.Filename, ", Subtitle:\n", string(s.Subtitle[: 100]), "\n")
    t.Logf("Filename: %s, Subtitle:\n%s\n", s.Filename, string(s.Subtitle[: 100]))
  }
}

func TestFetchMovieSubtitlesOrg(t *testing.T) {
  t.Parallel()

  tests := []struct {
    name string
    provider func(query string, options common.SearchOptions) ([]common.MovieListEntry, error)
  }{
    // {"FetchMovieSubtitlesOrg", FetchMovieSubtitlesOrg},
    // {"FetchMoviesubtitlesrtCom", FetchMoviesubtitlesrtCom},
    // {"OpenSubtitlesCom", FetchOpenSubtitlesCom},
    // {"PodnapisiNet", FetchPodnapisiNet},
    // {"SubdlCom", FetchSubdlCom},
    {"YifySubtitlesCh", FetchYifySubtitlesCh},
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      t.Parallel()
      t.Logf("Testing %s", tt.name)
      testProvider(t, tt.provider)
    })
  }
}

