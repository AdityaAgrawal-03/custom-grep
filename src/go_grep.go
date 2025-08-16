package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	recursivePtr := flag.Bool("r", false, "recursive search flag")
	invertPtr := flag.Bool("v", false, "inverted search flag")
	casePtr := flag.Bool("i", false, "performs case insensitive search")
	var re *regexp.Regexp
	var reErr error

	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		log.Fatalln("Usage: gogrep [-i] [-r] [-v] <pattern> [filename|directory]")
	}

	pattern := args[0]

	// Compile regex once

	if *casePtr {
		re, reErr = regexp.Compile("(?i)" + pattern)
	} else {
		re, reErr = regexp.Compile(pattern)
	}

	if reErr != nil {
		log.Fatalf("Invalid regular expression: %s", reErr)
	}

	fi, _ := os.Stdin.Stat()

	// check if input is piped
	if (fi.Mode() & os.ModeCharDevice) == 0 {
		// input is coming from command before pipe
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			line := scanner.Text()
			isMatched := re.MatchString(line)

			if *invertPtr && !isMatched {
				fmt.Printf("%s\n", line)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Error reading stdin: %v", err)
		}

		return
	}

	// file/dir mode
	if len(args) < 2 {
		log.Fatalln("Usage: gogrep [-r] [-v] <pattern> <filename|directory>")
	}

	target := args[1]

	if *recursivePtr {
		fmt.Println("first process", target, pattern)
		recursiveErr := filepath.WalkDir(target, func(path string, d os.DirEntry, err error) error {

			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			if strings.ToLower(filepath.Ext(path)) != ".txt" {
				return nil
			}

			searchFile(path, re)

			return nil

		})

		if recursiveErr != nil {
			fmt.Printf("Error walking the path %q: %v\n", target, recursiveErr)
		}
	} else {
		searchFile(target, re)
	}

}

func searchFile(path string, re *regexp.Regexp) {
	file, fileErr := os.Open(path)

	if fileErr != nil {
		log.Fatalf("failed to open file %v", fileErr)
		return
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if re.MatchString(line) {
			fmt.Printf("%s: %s\n", path, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading file %s: %v", path, err)
	}
}
