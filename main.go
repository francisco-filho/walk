package main

import (
	"io"
	"log"
	"os"
	"flag"
	"path/filepath"
)

type config struct {
	ext string
	archive string
	size int64
	list bool
	del bool
	logWriter io.Writer
}

func main(){
	root := flag.String("root", ".", "The root file")
	ext := flag.String("ext", "", "This is the extension")
	list := flag.Bool("list", false, "List the files")
	size := flag.Int64("size", 0, "The min size of the files")
	del := flag.Bool("del", false, "Delete the files")
	logfile := flag.String("log", "", "File to save the logs")
	archive := flag.String("archive", "", "Diretory to archive the files")

	flag.Parse()

	f := os.Stdout

	if *logfile != "" {
		f, err := os.OpenFile(*logfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		defer f.Close()
	}

	cfg := config{
		ext: *ext,
		list: *list,
		size: *size,
		del: *del,
		archive: *archive,
		logWriter: f,
	}

	if err := run(*root, os.Stdout, cfg); err != nil {
		log.Fatal(err)
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

			if cfg.archive != "" {
				if err := archiveFile(cfg.archive, root, path); err != nil {
					return err
				}
				
			}

			if cfg.del {
				logger := log.New(cfg.logWriter, "Deleted File: ", log.LstdFlags)
				return delFile(path, logger)
			}

			return listFiles(path, out)
		})
}

