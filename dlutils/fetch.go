package dlutils

import (
  "fmt"
  "bytes"

  "github.com/bytedance/sonic"
  "github.com/valyala/fasthttp"
  "github.com/ItsMeSamey/go_utils"
  "github.com/PuerkitoBio/goquery"
)

type ResponseStatus int
func (status ResponseStatus) OK() bool { return status >= 200 && status < 300 }
func (status ResponseStatus) Error() error {
  if status.OK() { return nil }
  return utils.WithStack(fmt.Errorf("HTTP error: Bad Status: %d", status))
}

type Header struct {
  Key string
  Value string
}
type FetchInit struct {
  Method       string
  Headers      []Header
  Body         []byte
  MaxRedirects int
}

func (init FetchInit) AsFasthttpRequest() (req fasthttp.Request) {
  if init.Method != "" { req.Header.SetMethod(init.Method) }
  for _, header := range init.Headers { req.Header.Set(header.Key, header.Value) }
  if init.Body != nil { req.SetBody(init.Body) }
  if init.MaxRedirects == 0 { init.MaxRedirects = 1 << 4 }
  return
}

func FetchText(url string, init FetchInit) (body []byte, status ResponseStatus, err error) {
  req := init.AsFasthttpRequest()
  resp := fasthttp.Response{}
  if err = utils.WithStack(fasthttp.DoRedirects(&req, &resp, init.MaxRedirects)); err != nil { return }

  return resp.Body(), ResponseStatus(resp.StatusCode()), nil
}

func FetchJson[T any](url string, init FetchInit) (json T, status ResponseStatus, err error) {
  body, status, err := FetchText(url, init)
  if err != nil { return }
  err = utils.WithStack(sonic.Unmarshal(body, &json))
  return 
}

func FetchHtml(url string, init FetchInit) (out *goquery.Document, status ResponseStatus, err error) {
  body, status, err := FetchText(url, init)
  if err != nil { return }
  doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
  return doc, status, utils.WithStack(err)
}

func Download(url string, init FetchInit, processResponse func(req *fasthttp.Response) error) (err error) {
  req := init.AsFasthttpRequest()
  resp := fasthttp.Response{}
  if err = utils.WithStack(fasthttp.DoRedirects(&req, &resp, init.MaxRedirects)); err != nil { return }

  return utils.WithStack(processResponse(&resp))
}

