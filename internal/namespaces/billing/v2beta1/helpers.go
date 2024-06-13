package billing

import (
	"path/filepath"
	"strings"
)

const (
	invoiceDefaultPrefix = "scaleway-invoice"
)

func fileNameWithoutExtTrimSuffix(fileName string) string {
	return strings.TrimSuffix(fileName, filepath.Ext(fileName))
}

func getDirFile(filePath string) (string, string) {
	dir := "."
	dirTmp, file := filepath.Split(filePath)
	if len(dirTmp) > 0 {
		dir = dirTmp
	}

	if file == "." {
		file = ""
	}

	return dir, file
}
