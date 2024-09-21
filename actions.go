package main

import (
	"log"
	"os"
	"io"
	"fmt"
	"path/filepath"
	"compress/gzip"
)

func filterOut(path, ext string, minSize int64, info os.FileInfo) bool {
	if info.IsDir() || info.Size() < minSize {
		return true
	}

	if ext != "" && filepath.Ext(path) != ext {
		return true
	}

	return false
}

func listFiles(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
	return err
}

func delFile(path string, delLogger *log.Logger) error {
	if err := os.Remove(path); err != nil {
		return err
	}
	delLogger.Println(path)
	return nil
}

func archiveFile(dest, root, path string) error {
	s, err := os.Stat(dest)
	if err != nil {
		return err
	}
	if !s.IsDir() {
		return fmt.Errorf("%s is no a dir", path)
	}

	relDir, err := filepath.Rel(root, filepath.Dir(path))
	destFile := fmt.Sprintf("%s.gz", filepath.Base(path))
	targetPath := filepath.Join(dest, relDir, destFile)

	if err := os.MkdirAll(filepath.Base(targetPath), 0755); err != nil {
		return err
	}

	out, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer in.Close()

	gw := gzip.NewWriter(out)
	gw.Name = filepath.Base(path)

	if _, err = io.Copy(gw,in); err != nil {
		return err
	}
	if err = gw.Close(); err != nil {
		return err
	}

	return nil
}

