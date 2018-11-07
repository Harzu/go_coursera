package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

type byMode []os.FileInfo
type byName []os.FileInfo

func (f byMode) Len() int           { return len(f) }
func (f byMode) Less(i, j int) bool { return f[i].Mode() < f[j].Mode() }
func (f byMode) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

func (f byName) Len() int           { return len(f) }
func (f byName) Less(i, j int) bool { return f[i].Name() < f[j].Name() }
func (f byName) Swap(i, j int)      { f[i], f[j] = f[j], f[i] }

func getTab(index, lastIndex int) string {
	if index == lastIndex {
		return "\t"
	}

	return "│\t"
}

func fileSize(file os.FileInfo) string {
	if file.IsDir() {
		return ""
	} else if file.Size() <= 0 {
		return "(empty)"
	} else {
		return "(" + fmt.Sprint(file.Size()) + "b)"
	}
}

func getSortFiles(path string, flag bool) []os.FileInfo {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	if !flag {
		sort.Sort(byMode(files))
	} else {
		sort.Sort(byName(files))
	}

	return files
}

func getTree(path, result, tab string, printFiles bool) string {
	sortFiles := getSortFiles(path, printFiles)

	for index, file := range sortFiles {
		var line string
		if !printFiles && !file.IsDir() {
			continue
		}

		if index == len(sortFiles)-1 {
			line += tab + "└───"
		} else {
			line += tab + "├───"
		}

		result += line + file.Name() + " " + fileSize(file) + "\n"
		if file.IsDir() {
			newPath := filepath.Join(path, file.Name())
			result = getTree(newPath, result, tab+getTab(index, len(sortFiles)-1), printFiles)
		}
	}

	return result
}

func dirTree(out io.Writer, path string, printFiles bool) error {
	var result, tab string
	result = getTree(path, result, tab, printFiles)
	fmt.Fprintln(out, result)
	return nil
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}
