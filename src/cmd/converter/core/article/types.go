package article

type HtmlString string

type PageInterface interface {
	GetPath() string
	GetContent() string
	IsPublish() bool
}
