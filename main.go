package main

import (
	"io"
	"os"
	"flag"
	"path/filepath"
)

type config struct {
	ext string
	size int64
	list bool
	del bool
}

func main(){
	root := flag.String("root", ".", "The root file")
	ext := flag.String("ext", "", "This is the extension")
	list := flag.Bool("list", true, "List the files")
	size := flag.Int64("size", 0, "The min size of the files")
	del := flag.Bool("del", false, "Delete the files")

	flag.Parse()

	cfg := config{
		ext: *ext,
		list: *list,
		size: *size,
		del: *del,
	}

	if err := run(*root, os.Stdout, cfg); err != nil {
		os.Exit(1)
	}
}

func run(root string, out io.Writer, cfg config) error {
	return filepath.Walk(root, 
		func (path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if filterOut(path, cfg.ext, cfg.size, info) {
				return nil
			}

			if cfg.list {
				return listFiles(path, out)
			}

			if cfg.del {
				return delFile(path)
			}

			return listFiles(path, out)
		})
}

