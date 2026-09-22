package lib

import (
	"path/filepath"
	"strings"
)

func GetTopDir(path string) string {
	parts := strings.Split(path, string(filepath.Separator))
	topDir := parts[0]
	return topDir
}

func DivideTopDir(path string) (string, string) {
	parts := strings.Split(path, string(filepath.Separator))
	return parts[0], strings.Join(parts[1:], string(filepath.Separator))
}
