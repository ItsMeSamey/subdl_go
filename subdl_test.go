package subdl

import (
  "fmt"
  "testing"

  "github.com/ItsMeSamey/subdl_go/common"
  "github.com/ItsMeSamey/subdl_go/dlutils"
  "github.com/ItsMeSamey/subdl_go/providers"

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
    println("Filename:", s.Filename, ", Subtitle:\n", string(s.Subtitle[:100]), "\n")
    t.Logf("Filename: %s, Subtitle:\n%s\n", s.Filename, string(s.Subtitle[:100]))
  }
}

func TestFetchMovieSubtitlesOrg(t *testing.T) {
  t.Parallel()

  tests := []struct {
    name     string
    provider func(query string, options common.SearchOptions) ([]common.MovieListEntry, error)
  }{
    // {"FetchMovieSubtitlesOrg", providers.FetchMovieSubtitlesOrg},
    // {"FetchMoviesubtitlesrtCom", providers.FetchMoviesubtitlesrtCom},
    // {"OpenSubtitlesCom", providers.FetchOpenSubtitlesCom},
    // {"PodnapisiNet", providers.FetchPodnapisiNet},
    // {"SubdlCom", providers.FetchSubdlCom},
    // {"YifySubtitlesCh", providers.FetchYifySubtitlesCh},
  }

  for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
      t.Parallel()
      t.Logf("Testing %s", tt.name)
      testProvider(t, tt.provider)
    })
  }
}

func TestReadmeBasic(t *testing.T) {
  t.Parallel()
  options := common.SearchOptions{
    Language: common.LangEN,
  }

  result, err := Download("The Matrix", options, providers.FetchOpenSubtitlesCom)
  if err != nil { t.Fatal("Error fetching subtitles:", err) }

  fmt.Println("File: ", result.Subtitles[0].Filename)
  fmt.Println("Subtitle: ", string(result.Subtitles[0].Subtitle[:100])) // First 100 characters only for readability
}


