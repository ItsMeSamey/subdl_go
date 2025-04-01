package subdl

import (
  "encoding/json"
  "fmt"
  "testing"

  "github.com/ItsMeSamey/go_utils"

  "github.com/ItsMeSamey/subdl_go/common"
  "github.com/ItsMeSamey/subdl_go/dlutils"
  "github.com/ItsMeSamey/subdl_go/providers"

  "github.com/ItsMeSamey/go_fuzzy"
  "github.com/ItsMeSamey/go_fuzzy/heuristics"
  "github.com/ItsMeSamey/go_fuzzy/transformers"
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
  for _, m := range movies { fmt.Println(m.Data().Title) }

  subtitles, err := movies[0].ToSubtitleLinks()
  if err != nil {
    t.Fatalf("Error fetching subtitles: %v", err)
  }
  if len(subtitles) == 0 {
    t.Fatalf("No subtitles found.")
  }
  for _, s := range subtitles {
    oldp := s.Data().Parent
    s.Data().Parent = nil
    output, _ := json.Marshal(s.Data())
    s.Data().Parent = oldp
    fmt.Println(string(output))
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
      fmt.Printf("Testing %s\n", tt.name)
      testProvider(t, tt.provider)
    })
  }
}

func TestReadmeBasic(t *testing.T) {
  t.Parallel()
  options := common.SearchOptions{
    Language: common.LangEN,
  }

  result, err := Download("The Matrix", options, providers.FetchOpenSubtitlesCom, common.DownloadOptions{})
  if err != nil { t.Fatal("Error fetching subtitles:", err) }

  fmt.Println("File: ", result.Subtitles[0].Filename)
  fmt.Println("Subtitle: ", string(result.Subtitles[0].Subtitle[:100])) // First 100 characters only for readability
}

func TestReadmeAdvanced(t *testing.T) {
  return
  t.Parallel()
  options := common.SearchOptions{
    Language: common.LangEN,
    Sorter: fuzzy.Sorter[float32, string, string]{
      Scorer: fuzzy.Scorer[float32, string, string]{
        ScoreFn: heuristics.Wrap[float32](heuristics.LevenshteinSimilarityPercentage),
        Transformer: transformers.Lowercase(),
      },
    },
  }

  movies, err := providers.FetchMovieSubtitlesOrg("The Matrix", options)
  if err != nil {
    fmt.Println("Error fetching movie list:", err)
    return
  }

  // See sorter documentation
  options.Sorter.SortAny(fuzzy.ToSwapper(movies, func (m common.MovieListEntry) string {return m.Data().Title}), "matrix reloaded")

  reloadedMovie := movies[0]
  subtitles, err := reloadedMovie.ToSubtitleLinks()
  if err != nil {
    fmt.Println("Error fetching subtitle list:", err)
    return
  }

  fmt.Println("Found", len(subtitles), "subtitles for", reloadedMovie.Data().Title)
  result, err := dlutils.DownloadSubtitleEntry(subtitles[0])

  fmt.Println("File: ", result.Subtitles[0].Filename)
  fmt.Println("Subtitle: ", string(result.Subtitles[0].Subtitle[:100])) // First 100 characters only for readability
}

