package smms

import (
	"encoding/json"
	multipartreader "github.com/lzjluzijie/MultipartReader"
	"io/ioutil"
	"log"
	"net/http"
	"path"
)

type Uploader struct {
}

type UploadResponse struct {
	Code  string
	Data UploadResponseData
}

type UploadResponseData struct {
	FileName  string
	StoreName string
	Size  int64
	URL   string
}

func (uploader Uploader) Upload(src string) (url string, err error) {
	name := path.Base(src)
	resp, err := http.Get(src)
	if err != nil {
		return
	}

	reader := multipartreader.NewMultipartReader()
	reader.AddFormReader(resp.Body, "smfile", name, resp.ContentLength)

	req, err := http.NewRequest("POST", "https://sm.ms/api/upload", reader)
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

	//log.Println(uResp)

	url = uResp.Data.URL
	return
}