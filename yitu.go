package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/lzjluzijie/MultipartReader"

	"github.com/urfave/cli"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	app := &cli.App{
		Name:        "yitu",
		Usage:       "A simple tool that helps you replace images in markdown files.",
		Description: "https://github.com/lzjluzijie/yitu",
		Version:     "initial",
		Author:      "Halulu",
		Email:       "lzjluzijie@gmail.com",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "path",
				Usage: "filepath",
			},
		},
		Action: func(c *cli.Context) (err error) {
			p := c.String("p")
			if p == "" {
				panic("path is empty!")
			}

			log.Printf("Start to read %s", p)

			f, err := ioutil.ReadFile(p)
			if err != nil {
				log.Println(err.Error())
				return
			}
			old := string(f)

			replaced, err := Replace(old)
			if err != nil {
				log.Println(err.Error())
				return
			}

			err = ioutil.WriteFile(p, []byte(replaced), 0600)
			if err != nil {
				log.Println(err.Error())
				return
			}

			log.Println("done")
			return
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf(err.Error())
	}
}

type Image struct {
	Name  string
	Size  int64
	Short string
	URL   string
}

func Replace(old string) (replaced string, err error) {
	replaced = old

	re := regexp.MustCompile(`!\[(.*)\]\((.*)\)`)
	matches := re.FindAllStringSubmatch(old, -1)

	for _, match := range matches {
		url := match[2]
		name := path.Base(url)
		resp, err := http.Get(url)
		if err != nil {
			log.Println(err)
			continue
		}

		reader := multipartreader.NewMultipartReader()
		reader.AddFormReader(resp.Body, "tu", name, resp.ContentLength)

		req, err := http.NewRequest("POST", "https://6tu.halu.lu/api/upload", reader)
		if err != nil {
			log.Println(err)
			continue
		}

		reader.SetupHTTPRequest(req)
		resp, err = http.DefaultClient.Do(req)

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println(string(data))

		image := &Image{}
		err = json.Unmarshal(data, image)
		if err != nil {
			log.Printf("json unmarshal %s error: %s", string(data), err.Error())
			continue
		}

		replaced = strings.Replace(replaced, url, image.URL, -1)
		log.Printf("replace %s with %s", url, image.URL)
	}
	return
}
