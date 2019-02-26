package tu

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path"

	"github.com/lzjluzijie/MultipartReader"
)

type Uploader struct {
}

type UploadResponse struct {
	Name  string
	Size  int64
	Short string
	URL   string
}

func (uploader Uploader) Upload(src string) (url string, err error) {
	name := path.Base(src)
	resp, err := http.Get(src)
	if err != nil {
		return
	}

	reader := multipartreader.NewMultipartReader()
	reader.AddFormReader(resp.Body, "tu", name, resp.ContentLength)

	req, err := http.NewRequest("POST", "https://6tu.halu.lu/api/upload", reader)
	if err != nil {
		return
	}

	reader.SetupHTTPRequest(req)
	resp, err = http.DefaultClient.Do(req)

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	log.Println(string(data))

	uResp := &UploadResponse{}
	err = json.Unmarshal(data, uResp)
	if err != nil {
		return
	}

	url = uResp.URL
	return
}
