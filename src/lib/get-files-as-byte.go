package lib

import "os"

func GetFilesAsByte(path string) ([]byte, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return bytes, nil
}
