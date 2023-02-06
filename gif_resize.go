package gif_resize

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/gif"

	"github.com/nfnt/resize"
)

type GifResize struct {
	Width  int
	Height int
}

func Init(width, height int) *GifResize {
	return &GifResize{
		Width:  width,
		Height: height,
	}
}

// func main() {
// 	dat, err := os.ReadFile("./gitt.gif")
// 	if err != nil {
// 		fmt.Print(err)
// 	}
// 	var data, errorss = Process(dat)
// 	if errorss != nil {
// 		fmt.Print(errorss)
// 	}
// 	f, err := os.Create("./loading.gif")
// 	check(err)
// 	defer f.Close()
// 	n2, err := f.Write(data)
// 	check(err)
// 	fmt.Printf("wrote %d bytes\n", n2)
// }

func (gr *GifResize) Process(data []byte) ([]byte, error) {
	var r2 = bytes.NewReader(data)

	// Decode the original gif.
	im, err := gif.DecodeAll(r2)
	if err != nil {
		return nil, err
	}

	// Create a new RGBA image to hold the incremental frames.
	firstFrame := im.Image[0].Bounds()
	b := image.Rect(0, 0, firstFrame.Dx(), firstFrame.Dy())
	img := image.NewRGBA(b)

	// Resize each frame.
	for index, frame := range im.Image {
		bounds := frame.Bounds()
		previous := img
		draw.Draw(img, bounds, frame, bounds.Min, draw.Over)
		im.Image[index] = imageToPaletted(resize.Resize(uint(gr.Width), uint(gr.Height), img, resize.NearestNeighbor), frame.Palette)

		switch im.Disposal[index] {
		case gif.DisposalBackground:
			img = image.NewRGBA(b)
		case gif.DisposalPrevious:
			img = previous
		}
	}

	// Set image.Config to new height and width
	im.Config.Width = gr.Width
	im.Config.Height = gr.Height

	var buff bytes.Buffer
	gif.EncodeAll(&buff, im)

	return buff.Bytes(), nil
}

func imageToPaletted(img image.Image, p color.Palette) *image.Paletted {
	b := img.Bounds()
	pm := image.NewPaletted(b, p)
	draw.FloydSteinberg.Draw(pm, b, img, image.Point{})
	return pm
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// gomobile bind -androidapi=26 -target android -o build/android/gifResize.aar
// gomobile bind -target=ios -iosversion=12 -o build/ios/GifResize.xcframework
