package main

import (
	"image"
	_ "image/jpeg"
	"os"
)

type Texture struct {
	Width  int
	Height int
	Data   []byte
}

func LoadTexture(path string) Texture {

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	bounds := img.Bounds()
	width := bounds.Max.X
	height := bounds.Max.Y

	data := make([]byte, width*height*4)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			coord := ToArrayCoords(x, y, width, 4)

			data[coord] = byte(r >> 8)
			data[coord+1] = byte(g >> 8)
			data[coord+2] = byte(b >> 8)
			data[coord+3] = byte(a >> 8)
		}
	}

	return Texture{Width: width, Height: height, Data: data}
}

func SampleTexture(tex Texture, u, v float32) (byte, byte, byte, byte) {
	u = clamp(u, 0, 1)
	v = clamp(v, 0, 1)

	x := int(u * float32(tex.Width-1))
	y := int(v * float32(tex.Height-1))

	idx := (y*tex.Width + x) * 4
	return tex.Data[idx], tex.Data[idx+1], tex.Data[idx+2], tex.Data[idx+3]

}
