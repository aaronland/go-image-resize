package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-image-decode"
	"github.com/aaronland/go-image-encode"
	"github.com/aaronland/go-image-resize"
	"github.com/natefinch/atomic"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {

	max := flag.Int("max", 0, "")

	flag.Parse()

	ctx := context.Background()

	dec, err := decode.NewDecoder(ctx, "image://")

	if err != nil {
		log.Fatal(err)
	}

	paths := flag.Args()

	for _, path := range paths {

		// START OF something like go-image-reader
		
		fh, err := os.Open(path)

		if err != nil {
			log.Fatal(err)
		}

		defer fh.Close()

		im, format, err := dec.Decode(ctx, fh)

		if err != nil {
			log.Fatal(err)
		}

		// END OF something like go-image-reader
		
		im, err = resize.ResizeImageMax(ctx, im, *max)

		if err != nil {
			log.Fatal(err)
		}

		root := filepath.Dir(path)

		fname := filepath.Base(path)
		ext := filepath.Ext(path)

		short_name := strings.Replace(fname, ext, "", 1)
		new_name := fmt.Sprintf("%s-%d%s", short_name, *max, ext)

		new_path := filepath.Join(root, new_name)

		// START OF something like go-image-writer
		
		enc_uri := fmt.Sprintf("%s://", format)
		enc, err := encode.NewEncoder(ctx, enc_uri)

		if err != nil {
			log.Fatal(err)
		}

		var buf bytes.Buffer
		wr := bufio.NewWriter(&buf)

		err = enc.Encode(ctx, im, wr)

		if err != nil {
			log.Fatal(err)
		}

		wr.Flush()

		br := bytes.NewReader(buf.Bytes())

		err = atomic.WriteFile(new_path, br)

		if err != nil {
			log.Fatal(err)
		}

		// END OF something like go-image-writer

		log.Println(new_path)
	}

}
