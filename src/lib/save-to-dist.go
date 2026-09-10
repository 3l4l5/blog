package lib

import "html/template"

type DistInfo struct {
	target  string
	content template.HTML
}

func SaveToDist(distinfo DistInfo) {
}
