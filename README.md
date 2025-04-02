# subdl_go

`subdl_go` is a Go package designed to facilitate the search and download of movie subtitles from various online providers. It provides a common interface for fetching movie and subtitle data, making it easy to integrate subtitle downloading functionality into your Go applications.

## Table of Contents

- [Installation](#installation)
- [Providers](#providers)
- [Usage](#usage)
- [Types](#types)
- [Functions](#functions)
- [Contributing](#contributing)

## Installation

To install `subdl_go`, use `go get`:

```sh
go get github.com/ItsMeSamey/subdl_go
```

## Providers

The package supports multiple subtitle providers. Each provider has its own implementation for fetching movie and subtitle data. The currently supported providers are:

- `moviesubtitles.org`
- `moviesubtitlesrt.com`
- `opensubtitles.com`
- `podnapisi.net`
- `subdl.com`
- `yifysubtitles.ch`

## Basic Usage

Here's a basic example of how to use `subdl_go` to search for and download subtitles form opensubtitles.com:

```go
package main

import (
  "fmt"

  "github.com/ItsMeSamey/subdl_go"
  "github.com/ItsMeSamey/subdl_go/common"
  "github.com/ItsMeSamey/subdl_go/providers"
)

func main() {
  options := common.SearchOptions{
    Language: common.LangEN,
  }

  result, err := subdl.Download("The Matrix", options, providers.FetchOpenSubtitlesCom, common.DownloadOptions{})
  if err != nil {
    fmt.Println("Error fetching subtitles:", err)
    return
  }

  fmt.Println("File: ", result.Subtitles[0].Filename)
  fmt.Println("Subtitle: ", string(result.Subtitles[0].Subtitle[:100])) // First 100 characters only for readability
}
```

## Advanced Usage

```go
package main

import (
  "fmt"

  "github.com/ItsMeSamey/go_fuzzy"
  "github.com/ItsMeSamey/go_fuzzy/heuristics"
  "github.com/ItsMeSamey/go_fuzzy/transformers"

  "github.com/ItsMeSamey/subdl_go/common"
  "github.com/ItsMeSamey/subdl_go/dlutils"
  "github.com/ItsMeSamey/subdl_go/providers"
)

func main() {
  options := common.SearchOptions{
    Language: common.LangEN,
    Sorter: fuzzy.Sorter[float32, string, string]{
      Scorer: fuzzy.Scorer[float32, string, string]{
        ScoreFn: heuristics.LevenshteinSimilarityPercentage[float32, string, string],
        Transformer: transformers.Lowercase(),
      },
    },
  }

  movies, err := providers.FetchMovieSubtitlesOrg("The Matrix", options)
  if err != nil {
    fmt.Println("Error fetching movie list:", err)
    return
  }

  // See documentation for "github.com/ItsMeSamey/go_fuzzy"
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

```

## Types

#### `SearchOptions`

```go
type SearchOptions struct {
  Language common.LanguageID
  Sorter   fuzzy.Sorter[float32, string, string] // fuzzy is "github.com/ItsMeSamey/go_fuzzy"
}
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request on GitHub.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

---

This README provides a basic overview of the `subdl_go` package and its functionality. For more detailed information, please refer to the source code and comments within the package.
