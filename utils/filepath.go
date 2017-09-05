package utils

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func CopyFiles(src string, dst string) error {
	var err error
	files, err := ioutil.ReadDir(src)
	for _, file := range files {
		tsrc := filepath.Join(src, file.Name())
		tdst := filepath.Join(dst, file.Name())

		if file.IsDir() {
			if _, err := os.Stat(tdst); os.IsNotExist(err) {
				os.Mkdir(tdst, 0755)
			}
			err = CopyFiles(tsrc, tdst)
		} else {
			dstFile, err := os.Create(tdst)
			if err != nil {
				return err
			}

			srcFile, err := os.OpenFile(tsrc, os.O_RDONLY, 0666)
			if err != nil {
				dstFile.Close()
				return err
			}
			_, err = io.Copy(dstFile, srcFile)
			srcFile.Close()
			dstFile.Close()
		}
	}
	return err
}
