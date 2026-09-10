package helpers

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"sync"

	"github.com/3l4l5/blog/src/cmd/converter/core"
	"github.com/3l4l5/blog/src/lib"
)

const articleTemplate = "template/article.html"

func GetArticles(paths []string, parser func([]byte) (string, error)) ([]core.ArticlePage, []error) {
	files := make([]core.ArticlePage, len(paths))
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

		data := struct {
			Content template.HTML
		}{
			Content: template.HTML(body),
		}
		if err := tmpl.Execute(&buf, data); err != nil {
			panic(err)
		}

		id := filepath.Dir(path)

		article := core.ArticlePage{
			ID:       id,
			Content:  core.HtmlString(buf.String()),
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
