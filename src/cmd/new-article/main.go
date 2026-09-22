package newarticle

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/3l4l5/blog/src/lib"
	"gopkg.in/yaml.v3"
)

const filename = "article.md"

func createDatePath() string {
	articleDir := "articles"
	now := time.Now()
	year, month, date := now.Date()
	path := filepath.Join(
		articleDir,
		strconv.Itoa(year),
		fmt.Sprintf("%02d", int(month)),
		fmt.Sprintf("%02d", date),
	)
	return path
}

func createArticleId() string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, 10)

	for i := range b {
		v, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))

		b[i] = chars[v.Int64()]
	}
	return string(b)
}

func createArticleTemplate() (string, error) {
	const body = `
# タイトル
	`

	metadata := lib.ArticleMetaData{
		Title:      "タイトル",
		Date:       time.Now(),
		Categories: "tech",
		Tags:       []string{"tag1", "tag2"},
		Publish:    false,
	}
	frontMatter, err := yaml.Marshal(metadata)
	if err != nil {
		return "", err
	}
	content := fmt.Sprintf(
		"---\n%s---\n\n%s",
		frontMatter,
		body,
	)

	return content, nil
}

func CreateNewArticle() {
	targetPath := createDatePath()
	articleId := createArticleId()
	articleDirPath := filepath.Join(targetPath, articleId)

	err := os.MkdirAll(articleDirPath, 0755)
	if err != nil {
		log.Fatal(err)
	}

	targetFilePath := filepath.Join(articleDirPath, filename)

	fileObj, err := os.Create(articleDirPath + "/" + filename)
	if err != nil {
		log.Fatal(err)
	}
	defer fileObj.Close()

	articleTemplate, err := createArticleTemplate()
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(targetFilePath, []byte(articleTemplate), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
