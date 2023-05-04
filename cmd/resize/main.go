package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"image"
	
	"runtime"
        "github.com/mandykoh/prism"
        "github.com/mandykoh/prism/adobergb"
        "github.com/mandykoh/prism/displayp3"	
        "github.com/mandykoh/prism/srgb"
	"github.com/aaronland/go-image-resize"
	"github.com/aaronland/go-image-cli"		
)

// ToAdobeRGB converts all the coloura in 'im' to match the Adobe RGB colour profile.
func ToAdobeRGB(im image.Image) image.Image {

	input_im := prism.ConvertImageToNRGBA(im, runtime.NumCPU())
	new_im := image.NewNRGBA(input_im.Rect)

	for i := input_im.Rect.Min.Y; i < input_im.Rect.Max.Y; i++ {

		for j := input_im.Rect.Min.X; j < input_im.Rect.Max.X; j++ {

			inCol, alpha := adobergb.ColorFromNRGBA(input_im.NRGBAAt(j, i))
			outCol := srgb.ColorFromXYZ(inCol.ToXYZ())
			new_im.SetNRGBA(j, i, outCol.ToNRGBA(alpha))
		}
	}

	return new_im
}

// ToDisplayP3 converts all the coloura in 'im' to match the Apple Display P3 colour profile.
func ToDisplayP3(im image.Image) image.Image {

	input_im := prism.ConvertImageToNRGBA(im, runtime.NumCPU())
	new_im := image.NewNRGBA(input_im.Rect)

	for i := input_im.Rect.Min.Y; i < input_im.Rect.Max.Y; i++ {

		for j := input_im.Rect.Min.X; j < input_im.Rect.Max.X; j++ {

			inCol, alpha := displayp3.ColorFromNRGBA(input_im.NRGBAAt(j, i))
			outCol := srgb.ColorFromXYZ(inCol.ToXYZ())
			new_im.SetNRGBA(j, i, outCol.ToNRGBA(alpha))
		}
	}

	return new_im
}

func main() {

	max := flag.Int("max", 0, "")

	flag.Parse()

	ctx := context.Background()

	cb := func(ctx context.Context, im image.Image, path string) (image.Image, string, error) {

		new_im, err := resize.ResizeImageMax(ctx, im, *max)

		if err != nil {
			return nil, "", err
		}

		// new_im = ToDisplayP3(new_im)
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
