package article

type TagIndexPage struct {
	TagName string
	Body    HtmlString
}

func (t *TagIndexPage) GetPath() string {
	return "tags/" + t.TagName
}

func (t *TagIndexPage) GetContent() string {
	return string(t.Body)
}
