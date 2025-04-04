package dlutils

import (
  "io"
  "bytes"
  "strings"
  "archive/zip"

  "github.com/ItsMeSamey/subdl_go/common"

  "github.com/ItsMeSamey/go_fuzzy"
  "github.com/ItsMeSamey/go_utils"
  "github.com/valyala/fasthttp"
)

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
    sorter.SortAny(fuzzy.ToSwapper(srtFiles, func (f *zip.File) string { return f.Name }), target)
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

func extractFilename(header string) (string) {
  const filenamePrefix = `filename="`
  startIndex := strings.Index(header, filenamePrefix)
  if startIndex == -1 { return "" }
  startIndex += len(filenamePrefix)

  var filename strings.Builder
  escaped := false
  for i := startIndex; i < len(header); i++ {
    char := header[i]

    if char == '"' && !escaped { break } // Closing double quote

    if escaped {
      switch char {
      case '"':
        filename.WriteByte('"')
      case '\\':
        filename.WriteByte('\\')
      default:
        filename.WriteByte('\\')
        filename.WriteByte(char)
      }
      escaped = false
    } else if char == '\\' {
      escaped = true
    } else {
      filename.WriteByte(char)
    }
  }

  // if escaped { return "" } // unterminated escape sequence in filename
  return filename.String()
}

func DownloadSubtitleEntry(entry common.SubtitleListEntry) (retval common.DownloadedSubtitle, err error) {
  url, err := entry.DownloadLink()
  if err = utils.WithStack(err); err != nil { return }

  req := fasthttp.Request{}
  req.SetRequestURI(url)
  resp := fasthttp.Response{}
  if err = utils.WithStack(fasthttp.DoRedirects(&req, &resp, 1 << 16)); err != nil { return }

  data := resp.Body()
  filename := extractFilename(string(resp.Header.Peek("Content-Disposition")))
  if filename == "" { filename = entry.Data().Filename }
  if filename == "" { filename = entry.Data().Parent.Data().Title }
  if filename == "" {
    uriPath := req.URI().Path()
    if len(uriPath) > 256 { uriPath = uriPath[:256] }
    filename = string(uriPath)
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

