package helpers

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/3l4l5/blog/src/cmd/converter/core/images"
	"github.com/3l4l5/blog/src/lib"
)

func GetImages() ([]images.Image, error) {
	filepath, err := listImageFilePath()
	if err != nil {
		return []images.Image{}, err
	}
	return readImageFiles(filepath)
}

func listImageFilePath() ([]string, error) {
	paths, err := lib.ListFiles("./articles")
	if err != nil {
		return []string{}, err
	}
	return filterImagePath(paths), nil
}

func filterImagePath(paths []string) []string {
	var result []string
	for _, path := range paths {
		filename := filepath.Base(path)
		re := regexp.MustCompile(`(?i)\.(jpg|jpeg|png|gif|webp|bmp|svg|ico|tiff?)$`)
		if re.MatchString(filename) {
			result = append(result, path)
		}
	}
	return result
}

func readImageFiles(filepaths []string) ([]images.Image, error) {
	results := []images.Image{}
	for _, filepath := range filepaths {
		re := regexp.MustCompile(`^articles/\d{4}/\d{2}/\d{2}/([A-Za-z0-9]{10})/`)
		var id string
		matches := re.FindStringSubmatch(filepath)
		if len(matches) > 1 {
			id = matches[1]
			results = append(results,
				images.Image{
					ArticleId:  id,
					OriginPath: filepath,
				},
			)
		} else {
			return nil, fmt.Errorf("invalid directory format.")
		}
	}
	return results, nil
}
