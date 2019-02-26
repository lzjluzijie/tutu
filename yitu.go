package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
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
				Name:  "p",
				Usage: "filepath",
			},
		},
		Action: func(c *cli.Context) (err error) {
			p := c.String("p")
			if p == "" {
				panic("path is empty!")
			}

			fi, err := os.Stat(p)
			if err != nil {
				panic(err.Error())
			}

			if !fi.IsDir() {
				log.Printf("%s is a file", p)
				err = HandleFile(p)
				if err != nil {
					log.Println(err.Error())
				}
			} else {
				log.Printf("%s is a dir", p)
				list := make([]string, 0)
				err = filepath.Walk(p, func(p string, fi os.FileInfo, err error) error {
					if err != nil {
						panic(err)
					}

					if IsMarkdown(p) {
						list = append(list, p)
						return nil
					}

					log.Printf("%s is not a markdown file, skipped", p)
					return nil
				})
				if err != nil {
					panic(err)
				}

				log.Printf("markdown files: %v", list)
				for _, p := range list {
					err = HandleFile(p)
					if err != nil {
						log.Println(err.Error())
					}
				}
			}

			log.Println("done")
			return
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Println(err.Error())
	}
	return
}

func IsMarkdown(p string) (is bool) {
	is = false
	if path.Ext(p) == ".md" ||
		path.Ext(p) == ".markdown" ||
		path.Ext(p) == ".mdown" ||
		path.Ext(p) == ".mkdn" ||
		path.Ext(p) == ".mkd" ||
		path.Ext(p) == ".mdwn" ||
		path.Ext(p) == ".mdtxt" ||
		path.Ext(p) == ".mdtext" {
		is = true
	}
	return
}

func HandleFile(p string) (err error) {
	f, err := ioutil.ReadFile(p)
	if err != nil {
		return
	}
	old := string(f)

	replaced, err := Replace(old)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(p, []byte(replaced), 0600)
	if err != nil {
		return
	}
	return
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
