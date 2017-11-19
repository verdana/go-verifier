package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/fatih/color"
)

const (
	appname = "go-verifier"
)

var (
	hashType = flag.String("hash", "md5", "Specify hashtype, values: md5, sha1, sha256")
	helpText = flag.Bool("help", false, "Show this help information")
	nopath   = flag.Bool("nopath", false, "Without full path name")
	upper    = flag.Bool("upper", false, "Get hash in uppercase")
	verify   = flag.String("verify", "", "Read and verify the checksum file")
)

var hasher *Hasher

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", appname)
	fmt.Fprintf(os.Stderr, "\t%s [flags]                Runs on all files in current directory\n", appname)
	fmt.Fprintf(os.Stderr, "\t%s [flags] [file|dir]     Compute for given files or directories\n", appname)
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *helpText {
		flag.Usage()
		return
	}

	hasher = NewHasher()
	hasher.SetHashType(*hashType)

	if flag.NArg() == 0 {
		if *verify != "" {
			verifyChecksum()
			return
		}

		hashDir(".")
	} else {
		var files []string
		for _, arg := range flag.Args() {
			file, err := os.Stat(arg)
			if err != nil {
				continue // ignore error
			}

			switch mode := file.Mode(); {
			case mode.IsDir():
				if fpath, err := filepath.Abs(arg); err == nil {
					files = append(files, scanDir(fpath)...)
				}
			case mode.IsRegular():
				if fpath, err := filepath.Abs(arg); err == nil {
					files = append(files, fpath)
				}
			}
		}

		files = cleanSlices(files)
		for _, fpath := range files {
			printResult(hasher.HashFile(fpath), fpath)
		}
	}
}

func cleanSlices(elements []string) []string {
	duplicates := map[string]bool{}
	result := []string{}

	for v := range elements {
		if duplicates[elements[v]] == true {
		} else {
			duplicates[elements[v]] = true
			result = append(result, elements[v])
		}
	}

	sort.Strings(result)
	return result
}

func scanDir(dirname string) []string {
	finfo, err := ioutil.ReadDir(dirname)
	if err != nil {
		log.Fatal(err)
	}

	var files []string
	for _, file := range finfo {
		if file.IsDir() {
			continue
		}
		files = append(files, filepath.Join(dirname, file.Name()))
	}
	return files
}

func hashDir(dirname string) {
	files := scanDir(dirname)
	for _, fpath := range files {
		printResult(hasher.HashFile(fpath), fpath)
	}
}

func printResult(checksum string, fpath string) {
	if *nopath == true {
		fpath = filepath.Base(fpath)
	}
	fmt.Printf("%s  %s\n", checksum, fpath)
}

func verifyChecksum() {
	if file, err := os.Stat(*verify); os.IsNotExist(err) || file.IsDir() {
		fmt.Printf("%s is not a regular file\n", *verify)
		return
	}

	file, err := os.Open(*verify)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		subs := strings.Fields(scanner.Text())

		if hasher.Verify(subs[0], subs[1]) {
			fmt.Fprintf(color.Output, "%s  %s\n", subs[0], color.GreenString("OK"))
		} else {
			fmt.Fprintf(color.Output, "%s  %s\n", subs[0], color.RedString("ERR"))
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
