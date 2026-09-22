package article

import (
	"github.com/3l4l5/blog/src/lib"
)

type ArticlePage struct {
	ID       string
	Metadata lib.ArticleMetaData
	Content  HtmlString
}

type ArticlePageImageReplaced ArticlePage

func (a *ArticlePageImageReplaced) GetPath() string {
	path := a.Metadata.Date.Format("2006/01/02")
	return "blog/" + path + "/" + a.ID + ".html"
}

func (a *ArticlePageImageReplaced) GetContent() string {
	return string(a.Content)
}

func (a *ArticlePageImageReplaced) IsPublish() bool {
	return a.Metadata.Publish
}
