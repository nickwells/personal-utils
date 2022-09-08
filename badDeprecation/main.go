package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func main() {
	var lineNum int
	lastLine := ""
	depRE := regexp.MustCompile("^// Deprecated:")
	blankLineRE := regexp.MustCompile("^// *$")
	for _, filename := range os.Args[1:] {
		var f *os.File
		var err error
		f, err = os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening: %q : %v\n", filename, err)
			continue
		}
		lineNum = 0
		var scanner *bufio.Scanner = bufio.NewScanner(f)
		for scanner.Scan() {
			lineNum++
			line := scanner.Text()
			if depRE.MatchString(line) {
				if !blankLineRE.MatchString(lastLine) {
					fmt.Println(filename, ":", lineNum, ":", line)
				}
			}
			lastLine = line
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading %q : %v\n", filename, err)
		}
		f.Close()
	}
}
