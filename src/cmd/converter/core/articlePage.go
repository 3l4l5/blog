package core

import "github.com/3l4l5/blog/src/lib"

type ArticlePage struct {
	ID       string
	Metadata lib.ArticleMetaData
	Content  HtmlString
}

func (a *ArticlePage) GetPath() string {
	return "blog/" + a.ID + ".html"
}

func (a *ArticlePage) GetContent() string {
	return string(a.Content)
}

func (a *ArticlePage) IsPublish() bool {
	return a.Metadata.Publish
}
