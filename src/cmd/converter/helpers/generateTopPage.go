package helpers

import (
	"bytes"
	"fmt"
	"html/template"
	"sort"

	"github.com/3l4l5/blog/src/cmd/converter/core/article"
	"github.com/google/uuid"
)

const topPageTemplatePath = "./template/home.html"

func GenerateTopPage(articles []article.ArticlePageImageReplaced) (article.TopPage, error) {
	tmpl := template.Must(template.ParseFiles(topPageTemplatePath))

	type articlesData struct {
		Title       string
		Description string
		Path        string
		Date        string
		Tags        []string
	}
	dataList := []articlesData{}
	sort.Slice(articles, func(i, j int) bool { return articles[i].Metadata.Date.After(articles[j].Metadata.Date) })
	for _, article := range articles {
		y, m, d := article.Metadata.Date.Date()
		dataList = append(dataList, articlesData{
			Title:       article.Metadata.Title,
			Description: article.Metadata.Description,
			Date:        fmt.Sprintf("%d/%d/%d", y, m, d),
			Path:        article.GetPath(),
			Tags:        article.Metadata.Tags,
		})
		article.GetPath()
	}
	id, _ := uuid.NewV7()
	data := struct {
		Articles []articlesData
		BuildId  string
	}{
		Articles: dataList,
		BuildId:  id.String(),
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		panic(err)
	}
	return article.TopPage{
		Content: article.HtmlString(buf.String()),
	}, nil
}
