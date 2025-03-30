package dlutils

import (
  "io"
  "bytes"
  "regexp"
  "strings"
  "archive/zip"

  "subtitle_downloader/common"

  "github.com/valyala/fasthttp"
  "github.com/ItsMeSamey/go_fuzzy"
  "github.com/ItsMeSamey/go_utils"
)

type sortable []*zip.File
func (s sortable) Len() int { return len(s) }
func (s sortable) Get(i int) string { return s[i].Name }
func (s sortable) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// An internal method that tries to unpack a zipped subtitle file
func unpackZipped(sorter fuzzy.Sorter[float32, string, string], target string, data []byte) (out []common.DownloadedSubtitleEntry, err error) {
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
    out = append(out, common.DownloadedSubtitleEntry{Subtitle: content, Filename: f.Name})
  }

  return out, nil
}

var filenameRegex = regexp.MustCompile(`filename="([^"]+?)"`)
func Download(entry common.SubtitleListEntry) (retval common.DownloadedSubtitle, err error) {
  url, err := entry.DownloadLink()
  if err = utils.WithStack(err); err != nil { return }

  req := fasthttp.Request{}
  req.SetRequestURI(url)
  resp := fasthttp.Response{}
  if err = utils.WithStack(fasthttp.DoRedirects(&req, &resp, 1 << 16)); err != nil { return }

  data := resp.Body()
  var filename string

  filenames := filenameRegex.FindSubmatch(resp.Header.Peek("Content-Disposition"))
  if len(filenames) >= 1 {
    filename = string(filenames[0])
  } else {
    filename = entry.Data().Filename
  }

  retval.Parent = entry
  if entry.IsZip() {
    retval.Subtitles, err = unpackZipped(entry.Data().Parent.Data().Options.Sorter, filename, data)
    if err = utils.WithStack(err); err != nil { return }
  } else {
    retval.Subtitles = []common.DownloadedSubtitleEntry{{
      Subtitle: data,
      Filename: filename,
    }}
  }

  return
}

