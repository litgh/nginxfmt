package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var ts int

func main() {
	override := flag.Bool("i", false, "override origin file")
	file := flag.String("f", "", "format nginx conf file path")
	dir := flag.String("d", "", "nginx conf dir")
	flag.IntVar(&ts, "ts", 4, "tabstop")

	flag.Parse()

	if *dir != "" {
		filepath.Walk(*dir, func(path string, info os.FileInfo, e error) error {
			if strings.HasSuffix(path, ".conf") {
				fmtFile(path, *override)
			}
			return nil
		})
	} else if *file != "" {
		fmtFile(*file, *override)
	} else {
		flag.Usage()
	}

}

func fmtFile(file string, override bool) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	buf := format(bufio.NewReader(f))
	if !override {
		fmt.Println(buf.String())
	} else {
		fmt.Println("format", file, "OK")
		ioutil.WriteFile(file, buf.Bytes(), 0666)
	}

}

func format(rd *bufio.Reader) (buf bytes.Buffer) {
	var indent = 0
	var space = strings.Repeat(" ", ts)
	for {
		l, err := rd.ReadString('\n')
		if err == io.EOF {
			break
		}
		l = strings.TrimSpace(l)
		comment := strings.HasPrefix(l, "#")
		if comment {
			buf.WriteString(strings.Repeat(space, indent))
			buf.WriteString(l)
			buf.WriteString("\n")
			continue
		}

		if l == "" {
			continue
		}

		if strings.Contains(l, "{") {
			buf.WriteString("\n")
			buf.WriteString(strings.Repeat(space, indent))
			indent++
		} else if strings.Contains(l, "}") {
			indent--
			buf.WriteString(strings.Repeat(space, indent))
		} else {
			buf.WriteString(strings.Repeat(space, indent))
		}
		buf.WriteString(l)
		buf.WriteString("\n")

	}
	return
}
