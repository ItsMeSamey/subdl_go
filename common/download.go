package common

import "github.com/ItsMeSamey/go_fuzzy"

type SearchOptions struct {
  Language LanguageID
  Sorter   fuzzy.Sorter[float32, string, string]
}
type DownloadOptions struct {
  // Set this to "" to not sort the movie list
  MovieListQuery     string
  MovieListSorter    fuzzy.Sorter[float32, string, string]

  // Set this to "" to not sort the subtitle list
  SubtitleListQuery  string
  SubtitleListSorter fuzzy.Sorter[float32, string, string]
}

type MovieListData struct {
  // The title of the movie
  Title   string

  // The Options when frtching the list of subtitles for this movie
  Options SearchOptions
}
type MovieListEntry interface {
  Data() *MovieListData
  ToSubtitleLinks() ([]SubtitleListEntry, error)
}

type SubtitleListData struct {
  Parent MovieListEntry
  // Name of the subtitle file
  Filename string
  // Language of the subtitle file
  Language string
  // Setting this to "" will force refetching when calling DownloadLink()
  Target string
}
type SubtitleListEntry interface {
  Data() *SubtitleListData
  // Weather the downloaded file is a zip file or not
  IsZip() bool
  // Returns the Download link from where we can fetch the subtitle file
  DownloadLink() (string, error)
}


type DownloadedSubtitleEntry struct {
  Subtitle []byte
  Filename  string
}
type DownloadedSubtitle struct {
  Parent SubtitleListEntry

  // May contain 0, or 1 or more subtitles
  Subtitles []DownloadedSubtitleEntry
}

