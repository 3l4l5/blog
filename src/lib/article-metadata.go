package lib

import "time"

type categories string

const (
	CategoriesTech categories = "tech"
	CategoriesLife categories = "life"
)

type ArticleMetaData struct {
	Title       string     `yaml:"title"`
	Date        time.Time  `yaml:"date"`
	Description string     `yaml:"description"`
	Categories  categories `yaml:"categories"`
	Tags        []string   `yaml:"tagas"`
	Publish     bool       `yaml:"publish"`
}
