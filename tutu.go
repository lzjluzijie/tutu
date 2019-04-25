package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lzjluzijie/tutu/uploaders"
	"github.com/lzjluzijie/tutu/uploaders/smms"
	"github.com/lzjluzijie/tutu/uploaders/yitu"

	"github.com/urfave/cli"
)

var uploader uploaders.Uploader

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	app := &cli.App{
		Name:        "tutu",
		Usage:       "A simple tool that helps you replace images in markdown files.",
		Description: "https://github.com/lzjluzijie/tutu",
		Version:     "initial",
		Author:      "Halulu",
		Email:       "lzjluzijie@gmail.com",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "p",
				Usage: "file or dir",
			},
			cli.StringFlag{
				Name:  "u",
				Usage: "uploader",
			},
		},
		Action: func(c *cli.Context) (err error) {
			switch c.String("u") {
			case "yitu":
				uploader = yitu.Uploader{}
				break
			case "yitu-webp":
				uploader = yitu.Uploader{W: "/webp"}
				break
			case "yitu-fhd":
				uploader = yitu.Uploader{W: "/fhd"}
				break
			case "yitu-fhdwebp":
				uploader = yitu.Uploader{W: "/fhdwebp"}
				break

			default:
				uploader = smms.Uploader{}
				break
			}

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

func Replace(old string) (replaced string, err error) {
	replaced = old

	re := regexp.MustCompile(`!\[(.*)\]\((.*)\)`)
	matches := re.FindAllStringSubmatch(old, -1)

	for _, match := range matches {
		src := match[2]
		url, err := uploader.Upload(src)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		replaced = strings.Replace(replaced, src, url, -1)
		log.Printf("replace %s with %s", src, url)
	}
	return
}
