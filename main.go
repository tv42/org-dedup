package main

import (
	"bytes"
	"crypto/sha1"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  %s <FILE >FILE\n", os.Args[0])
	flag.PrintDefaults()
}

func ComputeSHA1(data []byte) [sha1.Size]byte {
	h := sha1.New()
	r := bytes.NewReader(data)
	_, err := io.Copy(h, r)
	if err != nil {
		log.Fatalf("SHA-1 failed: %s", err)
	}
	var sum [sha1.Size]byte
	h.Sum(sum[0:0])
	return sum
}

func main() {
	prog := path.Base(os.Args[0])
	log.SetFlags(0)
	log.SetPrefix(prog + ": ")

	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 0 {
		usage()
		os.Exit(1)
	}

	data, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Cannot read data: %s", err)
	}

	re := regexp.MustCompile(`(?m)^\*([^*].*)?\n((\*\*|[^*]).*\n)*`)
	matches := re.FindAll(data, -1)

	var seen = map[[sha1.Size]byte]struct{}{}

	for _, match := range matches {
		hash := ComputeSHA1(match)
		_, present := seen[hash]
		if !present {
			r := bytes.NewReader(match)
			_, err := io.Copy(os.Stdout, r)
			if err != nil {
				log.Fatalf("Writing to standard output failed: %s", err)
			}
		}
		seen[hash] = struct{}{}
	}
}
