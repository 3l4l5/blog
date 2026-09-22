package helpers

import (
	"bytes"
	"html/template"

	"github.com/3l4l5/blog/src/cmd/converter/core/article"
	"github.com/3l4l5/blog/src/lib"
)

const tagPageTemplatePath = "template/tagIndexPage.html"

func GenerateTagIndexPage(articles []article.ArticlePageImageReplaced) ([]article.TagIndexPage, error) {

	articleByTag := make(map[string][]article.ArticlePageImageReplaced)
	for _, article := range articles {
		for _, tag := range article.Metadata.Tags {
			articleByTag[tag] = append(articleByTag[tag], article)
		}
	}

	var tagPages []article.TagIndexPage
	for tagName, articles := range articleByTag {
		var pageInfos []PageInfo
		for _, article := range articles {
			pageInfos = append(pageInfos, PageInfo{
				article.GetPath(),
				article.Metadata.Title,
			})
		}
		content, err := applyTemplate(pageInfos)
		if err != nil {
			return []article.TagIndexPage{}, err
		}
		tagPages = append(tagPages, article.TagIndexPage{
			TagName: tagName,
			Body:    article.HtmlString(content),
		})
	}

	return tagPages, nil
}

type PageInfo struct {
	Url   string
	Title string
}

func applyTemplate(pageInfos []PageInfo) (string, error) {
	htmlByte, err := lib.GetFilesAsByte(tagPageTemplatePath)
	if err != nil {
		return "", err
	}
	tagPageTemplate, err := template.New("TagPage").Parse(string(htmlByte))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tagPageTemplate.Execute(&buf, pageInfos); err != nil {
		return "", err
	}
	return buf.String(), nil
}
