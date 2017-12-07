package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	bblfsh "gopkg.in/bblfsh/client-go.v2"
)

func main() {
	backend := flag.String("b", "localhost:9432", "address of the bblfsh endpoint")
	lang := flag.String("l", "", "language of the file (required)")
	file := flag.String("f", "", "path of the file to parse; when empty parses stdin")
	flag.Parse()

	if *lang == "" {
		fmt.Fprintln(os.Stderr, "Missing language name.")
		flag.Usage()
		os.Exit(1)
	}

	if err := parse(*backend, *lang, *file); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func parse(backend, lang, file string) error {
	c, err := bblfsh.NewClient(backend)
	if err != nil {
		return fmt.Errorf("could not connect %s: %v", backend, err)
	}

	req := c.NewParseRequest().Language(lang)
	if file != "" {
		req = req.ReadFile(file)
	} else {
		b, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("could not read stdin: %v", err)
		}
		req = req.Content(string(b))
	}

	res, err := req.Do()
	if err != nil {
		return fmt.Errorf("received error from backend: %v", err)
	}

	fmt.Println(res)
	return nil
}
