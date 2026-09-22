package main

import (
	"errors"
	"log"
	"os"

	"github.com/3l4l5/blog/src/cmd/converter"
	"github.com/3l4l5/blog/src/cmd/converter/helpers"
	newarticle "github.com/3l4l5/blog/src/cmd/new-article"
)

func main() {

	args := os.Args

	if len(args) != 2 {
		log.Fatal(errors.New("args length must be 1."))
	} else if args[1] == "newarticle" {
		newarticle.CreateNewArticle()
	} else if args[1] == "build" {
		uploader := helpers.ImageUploaderFactory("build")
		converter.ConvertMarkdownToHtml(uploader)
	} else if args[1] == "publish" {
		uploader := helpers.ImageUploaderFactory("publish")
		converter.ConvertMarkdownToHtml(uploader)
	} else {
		log.Fatal("invalid option")
	}
}
