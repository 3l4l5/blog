package helpers_test

import (
	"reflect"
	"testing"

	"github.com/3l4l5/blog/src/cmd/converter/core/article"
	"github.com/3l4l5/blog/src/cmd/converter/core/images"
	"github.com/3l4l5/blog/src/cmd/converter/helpers"
	"github.com/3l4l5/blog/src/lib"
)

func TestReplaceImagePath(t *testing.T) {
	tests := []struct {
		name     string
		images   []images.ImageUploaded
		articles []article.ArticlePage
		want     []article.ArticlePageImageReplaced
		wantErr  bool
	}{
		{
			name: "記事に含まれる画像が一つ",
			images: []images.ImageUploaded{
				{
					ArticleId:   "0000000000",
					OriginPath:  "articles/2000/09/01/test.png",
					UploadedUrl: "http://example.com/images/result.png",
				},
			},
			articles: []article.ArticlePage{
				{
					ID:       "0000000000",
					Metadata: lib.ArticleMetaData{},
					Content: `
<h1>h1</h1>
<p>これは画像です</p>
<img src="./test.png" alt="画像" title="title">
<p>ところで、この画像はいい画像です</p>
`,
				},
			},
			want: []article.ArticlePageImageReplaced{
				{
					ID:       "0000000000",
					Metadata: lib.ArticleMetaData{},
					Content: `
<h1>h1</h1>
<p>これは画像です</p>
<img src="http://example.com/images/result.png" alt="画像" title="title">
<p>ところで、この画像はいい画像です</p>
`,
				},
			},
		},
		{
			name: "記事に含まれる画像が複数",
			images: []images.ImageUploaded{
				{
					ArticleId:   "0000000000",
					OriginPath:  "articles/2000/09/01/test1.png",
					UploadedUrl: "http://example.com/images/result1.png",
				},
				{
					ArticleId:   "0000000000",
					OriginPath:  "articles/2000/09/01/test2.jpg",
					UploadedUrl: "http://example.com/images/result2.jpg",
				},
			},
			articles: []article.ArticlePage{
				{
					ID:       "0000000000",
					Metadata: lib.ArticleMetaData{},
					Content: `
<h1>h1</h1>
<img src="./test1.png" alt="画像1">
<p>本文</p>
<img src="./test2.jpg" alt="画像2" title="画像2">
`,
				},
			},
			want: []article.ArticlePageImageReplaced{
				{
					ID:       "0000000000",
					Metadata: lib.ArticleMetaData{},
					Content: `
<h1>h1</h1>
<img src="http://example.com/images/result1.png" alt="画像1">
<p>本文</p>
<img src="http://example.com/images/result2.jpg" alt="画像2" title="画像2">
`,
				},
			},
		},
		{
			name: "複数の記事の画像をそれぞれ置換する",
			images: []images.ImageUploaded{
				{
					ArticleId:   "0000000001",
					OriginPath:  "articles/2000/09/01/test.png",
					UploadedUrl: "http://example.com/images/article1.png",
				},
				{
					ArticleId:   "0000000002",
					OriginPath:  "articles/2000/09/02/test.png",
					UploadedUrl: "http://example.com/images/article2.png",
				},
			},
			articles: []article.ArticlePage{
				{
					ID:       "0000000001",
					Metadata: lib.ArticleMetaData{},
					Content:  `<img src="./test.png" alt="画像">`,
				},
				{
					ID:       "0000000002",
					Metadata: lib.ArticleMetaData{},
					Content:  `<img src="./test.png" alt="画像">`,
				},
			},
			want: []article.ArticlePageImageReplaced{
				{
					ID:       "0000000001",
					Metadata: lib.ArticleMetaData{},
					Content:  `<img src="http://example.com/images/article1.png" alt="画像">`,
				},
				{
					ID:       "0000000002",
					Metadata: lib.ArticleMetaData{},
					Content:  `<img src="http://example.com/images/article2.png" alt="画像">`,
				},
			},
		},
		{
			name:   "画像を含まない記事は変更しない",
			images: []images.ImageUploaded{},
			articles: []article.ArticlePage{
				{
					ID:       "0000000000",
					Metadata: lib.ArticleMetaData{},
					Content: `
<h1>h1</h1>
<p>画像はありません。</p>
`,
				},
			},
			want: []article.ArticlePageImageReplaced{
				{
					ID:       "0000000000",
					Metadata: lib.ArticleMetaData{},
					Content: `
<h1>h1</h1>
<p>画像はありません。</p>
`,
				},
			},
		},
		{
			name: "titleなしの画像も置換する",
			images: []images.ImageUploaded{
				{
					ArticleId:   "0000000000",
					OriginPath:  "articles/2000/09/01/test.png",
					UploadedUrl: "http://example.com/images/result.png",
				},
			},
			articles: []article.ArticlePage{
				{
					ID:       "0000000000",
					Metadata: lib.ArticleMetaData{},
					Content:  `<img src="./test.png" alt="画像">`,
				},
			},
			want: []article.ArticlePageImageReplaced{
				{
					ID:       "0000000000",
					Metadata: lib.ArticleMetaData{},
					Content:  `<img src="http://example.com/images/result.png" alt="画像">`,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := helpers.ReplaceImagePath(tt.images, tt.articles)

			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, but got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf(
					"result mismatch\n\ngot:\n%#v\n\nwant:\n%#v",
					got,
					tt.want,
				)
			}
		})
	}
}
