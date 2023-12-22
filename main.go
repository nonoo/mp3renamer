package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const testMode = false

func mergeDir(dst, src string) error {
	fmt.Println("merging", src, "into", dst)

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, e := range entries {
		fmt.Println("  moving from", src+"/"+e.Name(), "to", dst+"/"+e.Name())

		dstExists := false
		if _, err := os.Stat(dst + "/" + e.Name()); !os.IsNotExist(err) {
			dstExists = true
		}

		if dstExists && e.IsDir() {
			fmt.Println("   ", dst+"/"+e.Name(), "is a directory and exists, merging")
			err = mergeDir(dst+"/"+e.Name(), src+"/"+e.Name())
			if err != nil {
				return err
			}
		} else {
			if !testMode {
				err = os.Rename(src+"/"+e.Name(), dst+"/"+e.Name())
				if err != nil {
					return err
				}
			}
		}
	}

	if testMode {
		return nil
	}

	return os.Remove(src)
}

func processDir(path string) (merged bool, err error) {
	dirs := []string{}
	files := []string{}

	entries, err := os.ReadDir(path)
	if err != nil {
		return false, err
	}

	fmt.Println("processing subdirs in", path)

	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e.Name())
		} else {
			files = append(files, e.Name())
		}
	}

	for _, dir1 := range dirs {
		dir1matcher := regexp.QuoteMeta(dir1)
		dir1matcher = strings.ReplaceAll(dir1matcher, "_", ".")
		dir1matcher = "^" + dir1matcher + "$"
		for _, dir2 := range dirs {
			if dir1 == dir2 {
				continue
			}

			match, _ := regexp.MatchString(dir1matcher, dir2)
			if match {
				err = mergeDir(filepath.Join(path, dir2), filepath.Join(path, dir1))
				return true, err
			}
		}
	}

	fmt.Println("processing files in", path)

	for _, file1 := range files {
		file1matcher := regexp.QuoteMeta(file1)
		file1matcher = strings.ReplaceAll(file1matcher, "_", ".")
		file1matcher = "^" + file1matcher + "$"
		for _, file2 := range files {
			if file1 == file2 {
				continue
			}

			match, _ := regexp.MatchString(file1matcher, file2)
			if match {
				file1Path := filepath.Join(path, file1)
				fmt.Println("  removing", file1Path)
				if !testMode {
					err = os.Remove(file1Path)
					if err != nil {
						return false, err
					}
				}
			}
		}
	}

	return false, err
}

func main() {
	// Check if argument is provided
	if len(os.Args) < 2 {
		log.Fatal("please provide a path")
	}

	root := os.Args[1]

	// Check if root is a valid directory.
	if _, err := os.Stat(root); os.IsNotExist(err) {
		log.Fatal("please provide a valid path")
	}

	mergedError := fmt.Errorf("merged")
	for walkErr := mergedError; walkErr == mergedError; {
		walkErr = filepath.WalkDir(root, func(path string, dirEntry fs.DirEntry, err error) error {
			if err != nil {
				fmt.Println(err)
				return err
			}

			if dirEntry.IsDir() {
				merged, err := processDir(path)
				if err != nil {
					return err
				}
				if merged { // Restart the walk.
					return mergedError
				}
			}
			return err
		})

		if walkErr != nil && walkErr != mergedError {
			log.Fatal(walkErr)
		}
	}

	fmt.Println("done")
}
