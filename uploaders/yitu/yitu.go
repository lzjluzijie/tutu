package yitu

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Uploader struct {
	W string
}

type UploadResponse struct {
	Name string
	Size int64
	URL  string
}

func (uploader Uploader) Upload(src string) (url string, err error) {
	resp, err := http.PostForm("https://t.halu.lu/api/upload", map[string][]string{"url": {src}})
	if err != nil {
		return
	}

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

	url = uResp.URL + uploader.W
	return
}
