package helpers

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"sync"

	"github.com/3l4l5/blog/src/cmd/converter/core/article"
	"github.com/3l4l5/blog/src/lib"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/parser"
	"go.abhg.dev/goldmark/frontmatter"
)

const articleTemplate = "template/article.html"

func GetArticles(paths []string, parser func([]byte) (string, error)) ([]article.ArticlePage, []error) {
	files := make([]article.ArticlePage, len(paths))
	errCh := make(chan error, len(paths))

	var wg sync.WaitGroup
	wg.Add(len(paths))
	tmpl := template.Must(template.ParseFiles(articleTemplate))

	routine := func(index int, path string) {
		defer wg.Done()
		content, err := lib.GetFilesAsByte(path)
		if err != nil {
			errCh <- err
			return
		}
		metadata, err := GetMetadataFromMarkdown(string(content))
		if err != nil {
			errCh <- err
			return
		}
		bodyRemovedMetadata, err := removeMeatadata(string(content))
		if err != nil {
			errCh <- err
			return
		}
		body, err := parser([]byte(bodyRemovedMetadata))
		if err != nil {
			errCh <- err
			return
		}

		var buf bytes.Buffer

		y, m, d := metadata.Date.Date()
		data := struct {
			Title   string
			Content template.HTML
			Tags    []string
			Date    string
		}{
			Title:   metadata.Title,
			Content: template.HTML(body),
			Tags:    metadata.Tags,
			Date:    fmt.Sprintf("%d/%d/%d", y, m, d),
		}
		if err := tmpl.Execute(&buf, data); err != nil {
			panic(err)
		}

		id := filepath.Base(filepath.Dir(path))

		article := article.ArticlePage{
			ID:       id,
			Content:  article.HtmlString(buf.String()),
			Metadata: metadata,
		}
		files[index] = article
	}
	for i, path := range paths {
		go routine(i, path)
	}
	wg.Wait()
	close(errCh)
	var errors []error
	for err := range errCh {
		errors = append(errors, err)
	}
	if len(errors) > 0 {
		return files, errors
	}
	return files, nil
}

func removeMeatadata(body string) (string, error) {
	parts := strings.SplitN(body, "---", 3)

	if len(parts) < 3 {
		return "", fmt.Errorf("Invalid Markdown shape.")
	}

	result := strings.TrimPrefix(parts[2], "\n")
	return result, nil
}

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
