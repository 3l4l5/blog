package main

import (
	"errors"
	"log"
	"os"

	"github.com/3l4l5/blog/src/cmd/converter"
	newarticle "github.com/3l4l5/blog/src/cmd/new-article"
)

func main() {

	args := os.Args

	if len(args) != 2 {
		log.Fatal(errors.New("args length must be 1."))
	}
	if args[1] == "newarticle" {
		newarticle.CreateNewArticle()
	}
	if args[1] == "build" {
		converter.ConvertMarkdownToHtml()
	} else {
		log.Fatal("invalid option")
	}
}
