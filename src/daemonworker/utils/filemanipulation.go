package utils

import (
	"archive/tar"
	"io"
	"log"
	"os"
	"path/filepath"
)

func CheckFilePath(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func CopyFile(source io.Reader, destPath string, filemode os.FileMode) error {
	targetFile, err := os.OpenFile(destPath, os.O_CREATE|os.O_WRONLY, filemode)
	if err != nil {
		log.Printf("Failed to create target uploaded: %+v\n", err)
		return err
	}
	defer targetFile.Close()
	_, err = io.Copy(targetFile, source)
	if err != nil {
		log.Printf("Failed to write source to target: %+v\n", err)
		return err
	}
	return nil
}

func Untar(source io.Reader, destPath string) error {
	reader := tar.NewReader(source)
	for {
		fh, err := reader.Next()
		if fh == nil || err != nil {
			if err == io.EOF {
				log.Println("Finished untarred file due to read to EOF.")
			} else if err != nil {
				log.Printf("Error occurred while untaring file: %+v", err)
			} else {
				log.Printf("Failed to untar file.")
			}
			break
		}
		targetPath := filepath.Join(destPath, fh.Name)
		switch fh.Typeflag {
		case tar.TypeDir:
			err = CheckFilePath(targetPath)
			if err != nil {
				log.Printf("Failed to mkdir for path: %s", targetPath)
				return err
			}
		case tar.TypeReg:
			err = CopyFile(reader, targetPath, os.FileMode(fh.Mode))
			if err != nil {
				log.Printf("Failed to write file: %s", fh.Name)
				return err
			}
		}
	}
	return nil
}
