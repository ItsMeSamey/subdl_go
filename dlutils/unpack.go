package dlutils

import (
	"io"
	"bytes"
	"strings"
	"archive/zip"

	"subtitle_downloader/common"

	"github.com/ItsMeSamey/go_fuzzy"
	"github.com/ItsMeSamey/go_utils"
)

type sortable []*zip.File
func (s sortable) Len() int { return len(s) }
func (s sortable) Get(i int) string { return s[i].Name }
func (s sortable) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// An internal method that tries to unpack a zipped subtitle file
func UnpackZipped(sorter fuzzy.Sorter[float32, string, string], target string, data []byte) (out []common.DownloadedSubtitle, err error) {
  r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
  if err != nil { return nil, utils.WithStack(err) }

  var allFiles []*zip.File
  for _, f := range r.File {
    if !f.FileInfo().IsDir() { allFiles = append(allFiles, f) }
  }

  var srtFiles []*zip.File
  for _, f := range allFiles {
    if strings.HasSuffix(f.Name, ".srt") { srtFiles = append(srtFiles, f) }
  }

  if len(srtFiles) == 0 { srtFiles = allFiles }
  if len(srtFiles) > 1 {
    sorter.SortAny(sortable(srtFiles), target)
  }
  for _, f := range srtFiles {
    rc, err := f.Open()
    if err != nil { return nil, utils.WithStack(err) }
    defer rc.Close()

    content, err := io.ReadAll(rc)
    if err != nil { return nil, utils.WithStack(err) }
    out = append(out, common.DownloadedSubtitle{Subtitle: content, Filename: f.Name})
  }

  return out, nil
}

