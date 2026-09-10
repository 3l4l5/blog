package converter

import (
	"bytes"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/3l4l5/blog/src/cmd/converter/core"
	"github.com/3l4l5/blog/src/cmd/converter/helpers"
	"github.com/yuin/goldmark"
)

const root = "articles"
const distDir = "dist/"

func getMarkdownArticlePaths(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".md" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return files, err
	}
	return files, nil
}

func parseMarkdownToHtml(data []byte) (string, error) {
	var buf bytes.Buffer

	if err := goldmark.Convert(data, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func ConvertMarkdownToHtml() {
	markdownArticlePaths, err := getMarkdownArticlePaths(root)
	if err != nil {
		log.Fatal(err)
	}
	articles, errs := helpers.GetArticles(markdownArticlePaths, parseMarkdownToHtml)
	if len(errs) > 0 {
		for _, err := range errs {
			log.Println(err)
		}
		os.Exit(1)
	}

	var pages []core.PageInterface
	// tagIndexPages, err := helpers.GenerateTagIndexPage(articles)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	for _, articlePage := range articles {
		pages = append(pages, &articlePage)
	}
	// for _, tagIndexPage := range tagIndexPages {
	// 	pages = append(pages, &tagIndexPage)
	// }

	filteredPages := filterTargetArticle(pages)

	err = saveToDist(filteredPages)
	if err != nil {
		log.Fatal(err)
	}
}

func filterTargetArticle(pages []core.PageInterface) []core.PageInterface {
	buf := []core.PageInterface{}
	for _, page := range pages {
		if page.IsPublish() {
			buf = append(buf, page)
		}
	}
	return buf
}

func saveToDist(pages []core.PageInterface) error {
	if len(pages) <= 0 {
		log.Print("No pages")
		return nil
	}
	info, err := os.Stat(distDir)
	if err != nil || !info.IsDir() {
		err = os.Mkdir(distDir, os.FileMode(0777))
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, page := range pages {
		dirpath := filepath.Dir(distDir + page.GetPath())
		if err := os.MkdirAll(dirpath, 0755); err != nil {
			return err
		}
		f, err := os.Create(distDir + page.GetPath())
		if err != nil {
			return err
		}
		defer f.Close()
		d := []byte(page.GetContent())
		n, err := f.Write(d)
		if err != nil {
			return err
		}
		fmt.Printf("%d bytes written", n)
	}
	return nil
	// for _, page := range pages {
	// 	page.GetPath()
	// 	page.GetContent()
	// }
}
