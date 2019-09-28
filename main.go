package main

import (
	"io"
	"os"
	"sort"
	"strconv"
)

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

func dirTree(out io.Writer, path string, printFiles bool) error {
	dirList(out, path, printFiles, "", "", true)
	return nil
}

func dirList(out io.Writer, path string, printFiles bool, prefix string, last string, isFirstLevel bool) error {
	dir, err := os.Open(path)
	if err != nil {
		panic(err.Error())
	}
	defer dir.Close()

	fileInfo, err := dir.Stat()
	if err != nil {
		panic(err.Error())
	}

	if fileInfo.IsDir() {
		if !isFirstLevel {
			out.Write([]byte(prefix + last + fileInfo.Name() + "\n"))
		}

		subFilesTmp, err := dir.Readdir(-1)
		if err != nil {
			return nil
		}

		var subFiles []string
		for _, subFile := range subFilesTmp {
			if !(printFiles || subFile.IsDir()) {
				continue
			}
			subFiles = append(subFiles, subFile.Name())
		}

		sort.Strings(subFiles)

		for n, subFName := range subFiles {
			sPref := ""
			sLast := ""
			if isFirstLevel {
				sPref = ""
			} else if last == "└───" {
				sPref = "\t"
			} else {
				sPref = "│\t"
			}
			if n == len(subFiles)-1 {
				sLast = "└───"
			} else {
				sLast = "├───"
			}

			dirList(out, path+string(os.PathSeparator)+subFName, printFiles, prefix+sPref, sLast, false)
		}
	} else if printFiles && !isFirstLevel {
		fSize := strconv.FormatInt(fileInfo.Size(), 10) + "b"
		if fSize == "0b" {
			fSize = "empty"
		}
		out.Write([]byte(prefix + last + fileInfo.Name() + " (" + fSize + ")" + "\n"))
	}

	return nil
}
