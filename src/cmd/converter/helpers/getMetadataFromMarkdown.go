package helpers

import (
	"bytes"

	"github.com/3l4l5/blog/src/lib"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

func GetMetadataFromMarkdown(content string) (lib.ArticleMetaData, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			&frontmatter.Extender{},
		),
	)

	var buf bytes.Buffer
	ctx := parser.NewContext()
	if err := md.Convert([]byte(content), &buf, parser.WithContext(ctx)); err != nil {
		return lib.ArticleMetaData{}, err
	}
	var metadata lib.ArticleMetaData
	if err := frontmatter.Get(ctx).Decode(&metadata); err != nil {
		return lib.ArticleMetaData{}, err
	}
	return metadata, nil
}
