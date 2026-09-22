package article

type TopPage struct {
	Content HtmlString
}

func (a *TopPage) GetPath() string {
	return "index.html"
}

func (a *TopPage) GetContent() string {
	return string(a.Content)
}

func (a *TopPage) IsPublish() bool {
	return true
}
