package helpers

import (
	"bytes"
	"html/template"

	"github.com/3l4l5/blog/src/cmd/converter/core"
	"github.com/3l4l5/blog/src/lib"
)

const tagPageTemplatePath = "template/tagIndexPage.html"

func GenerateTagIndexPage(articles []core.ArticlePage) ([]core.TagIndexPage, error) {

	articleByTag := make(map[string][]core.ArticlePage)
	for _, article := range articles {
		for _, tag := range article.Metadata.Tags {
			articleByTag[tag] = append(articleByTag[tag], article)
		}
	}

	var tagPages []core.TagIndexPage
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
			return []core.TagIndexPage{}, err
		}
		tagPages = append(tagPages, core.TagIndexPage{
			TagName: tagName,
			Body:    core.HtmlString(content),
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
