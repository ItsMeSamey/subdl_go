package dlutils

import (
  "bytes"

  "github.com/bytedance/sonic"
  "github.com/valyala/fasthttp"
  "github.com/ItsMeSamey/go_utils"
  "github.com/PuerkitoBio/goquery"
)

type FetchResponse struct {
  Status int
  Body   []byte
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

func Fetch(url string, init FetchInit) (out FetchResponse, err error) {
  req := init.AsFasthttpRequest()
  resp := fasthttp.Response{}
  if err = utils.WithStack(fasthttp.DoRedirects(&req, &resp, init.MaxRedirects)); err != nil { return }

  return FetchResponse{Status: resp.StatusCode(), Body: resp.Body()}, nil
}

func (resp *FetchResponse) OK() bool {
  return resp.Status >= 200 && resp.Status < 300
}
func (resp *FetchResponse) JSON(out any) (err error) {
  return utils.WithStack(sonic.Unmarshal(resp.Body, out))
}
func (resp *FetchResponse) HTML() (out *goquery.Document, err error) {
  doc, err := goquery.NewDocumentFromReader(bytes.NewReader(resp.Body))
  return doc, utils.WithStack(err)
}

func Download(url string, init FetchInit, processResponse func(req *fasthttp.Response) error) (err error) {
  req := init.AsFasthttpRequest()
  resp := fasthttp.Response{}
  if err = utils.WithStack(fasthttp.DoRedirects(&req, &resp, init.MaxRedirects)); err != nil { return }

  return utils.WithStack(processResponse(&resp))
}

