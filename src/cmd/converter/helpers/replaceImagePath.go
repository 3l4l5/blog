package helpers

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/3l4l5/blog/src/cmd/converter/core/article"
	"github.com/3l4l5/blog/src/cmd/converter/core/images"
)

func ReplaceImagePath(inputImages []images.ImageUploaded, articles []article.ArticlePage) ([]article.ArticlePageImageReplaced, error) {
	inputImagesMap := make(map[string][]images.ImageUploaded)
	for _, inputImage := range inputImages {
		inputImagesMap[inputImage.ArticleId] = append(inputImagesMap[inputImage.ArticleId], inputImage)
	}
	fmt.Print(inputImagesMap)
	results := []article.ArticlePageImageReplaced{}
	re := regexp.MustCompile(`(<img\b[^>]*\bsrc\s*=\s*["'])([^"']+)(["'][^>]*>)`)
	for _, targetArticle := range articles {
		result := re.ReplaceAllStringFunc(string(targetArticle.Content), func(match string) string {
			parts := re.FindStringSubmatch(match)
			path := parts[2]

			for _, image := range inputImagesMap[targetArticle.ID] {
				if filepath.Base(image.OriginPath) == filepath.Base(path) {
					return parts[1] + image.UploadedUrl + parts[3]
				}
			}

			return match
		})
		results = append(results, article.ArticlePageImageReplaced{
			ID:       targetArticle.ID,
			Metadata: targetArticle.Metadata,
			Content:  article.HtmlString(result),
		})
	}
	return results, nil
}
