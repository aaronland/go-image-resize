package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aaronland/go-image-resize"
	"github.com/aaronland/go-image-cli"	
	"log"
	"path/filepath"
	"strings"
	"image"
)

func main() {

	max := flag.Int("max", 0, "")

	flag.Parse()

	ctx := context.Background()

	cb := func(ctx context.Context, im image.Image, path string) (image.Image, string, error) {

		new_im, err := resize.ResizeImageMax(ctx, im, *max)

		if err != nil {
			return nil, "", err
		}

		root := filepath.Dir(path)

		fname := filepath.Base(path)
		ext := filepath.Ext(path)

		short_name := strings.Replace(fname, ext, "", 1)
		new_name := fmt.Sprintf("%s-%d%s", short_name, *max, ext)

		new_path := filepath.Join(root, new_name)
		
		return new_im, new_path, nil
	}

	paths := flag.Args()
	
	err := cli.Process(ctx, cb, paths...)

	if err != nil {
		log.Fatal(err)
	}

}
