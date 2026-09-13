package helpers

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/3l4l5/blog/src/cmd/converter/core"
)

const topPageTemplatePath = "./template/home.html"

func GenerateTopPage(articles []core.ArticlePage) (core.TopPage, error) {
	tmpl := template.Must(template.ParseFiles(topPageTemplatePath))

	type articlesData struct {
		Title       string
		Description string
		Path        string
		Date        string
	}
	dataList := []articlesData{}
	for _, article := range articles {
		y, m, d := article.Metadata.Date.Date()
		dataList = append(dataList, articlesData{
			Title:       article.Metadata.Title,
			Description: article.Metadata.Description,
			Date:        fmt.Sprintf("%d/%d/%d", y, m, d),
			Path:        article.GetPath(),
		})
		article.GetPath()
	}
	data := struct {
		Articles []articlesData
	}{
		Articles: dataList,
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		panic(err)
	}
	return core.TopPage{
		Content: core.HtmlString(buf.String()),
	}, nil
}
