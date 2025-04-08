package Structs

import (
	"Photon/Math"
	"github.com/mdouchement/hdr"
	"image"
	"os"
)

// RGB texture

type TextureRGB struct {
	data   []Math.Vector3
	Width  int
	Height int
}

func ReadTextureRGB(img string) *TextureRGB {
	imgf, err := os.Open(img)
	if err != nil {
		panic(err)
	}
	defer imgf.Close()
	imageData, _, err := image.Decode(imgf)
	if err != nil {
		panic(err)
	}
	texture := &TextureRGB{
		Width:  imageData.Bounds().Max.X,
		Height: imageData.Bounds().Max.Y,
		data:   make([]Math.Vector3, imageData.Bounds().Max.X*imageData.Bounds().Max.Y),
	}
	for y := 0; y < texture.Width; y++ {
		for x := 0; x < texture.Height; x++ {
			arrIdx := y*texture.Width + x
			r, g, b, _ := imageData.At(x, y).RGBA()
			texture.data[arrIdx] = Math.Vector3{
				X: float64(r) / 255,
				Y: float64(g) / 255,
				Z: float64(b) / 255,
			}
		}
	}
	return texture
}

func EmptyTextureRGB(width, height int) *TextureRGB {
	tex := &TextureRGB{
		data:   make([]Math.Vector3, width*height),
		Width:  width,
		Height: height,
	}
	return tex
}

func TextureRGBFromHDR(hdrIm hdr.Image) *TextureRGB {
	tex := &TextureRGB{
		Width:  hdrIm.Bounds().Max.X,
		Height: hdrIm.Bounds().Max.Y,
	}
	tex.data = make([]Math.Vector3, tex.Width*tex.Height)

	for y := 0; y < tex.Height; y++ {
		for x := 0; x < tex.Width; x++ {
			r, g, b, _ := hdrIm.HDRAt(x, y).HDRRGBA()
			tex.data[y*tex.Height+x] = Math.Vector3{r, g, b}
		}
	}
	return tex
}

func (texture *TextureRGB) At(uv Math.Vector2) Math.Vector3 {
	if uv.U > 1 || uv.V > 1 || uv.U < 0 || uv.V < 0 {
		panic("invalid UV coordinates")
	}
	x := int(uv.U * float64(texture.Width))
	y := int(uv.V * float64(texture.Height))
	return texture.data[y*texture.Height+x]
}

// Grayscale texture

type TextureGrayscale struct {
	data   []float64
	Width  int
	Height int
}

func ReadTextureGrayscale(img string) *TextureGrayscale {
	imgf, err := os.Open(img)
	if err != nil {
		panic(err)
	}
	defer imgf.Close()
	imageData, _, err := image.Decode(imgf)
	if err != nil {
		panic(err)
	}
	texture := &TextureGrayscale{
		Width:  imageData.Bounds().Max.X,
		Height: imageData.Bounds().Max.Y,
		data:   make([]float64, imageData.Bounds().Max.X*imageData.Bounds().Max.Y),
	}
	for y := 0; y < texture.Width; y++ {
		for x := 0; x < texture.Height; x++ {
			arrIdx := y*texture.Width + x
			r, g, b, _ := imageData.At(x, y).RGBA()
			texture.data[arrIdx] = (float64(r) + float64(g) + float64(b)) / 765
		}
	}
	return texture
}

func (texture *TextureGrayscale) At(uv Math.Vector2) float64 {
	if uv.U > 1 || uv.V > 1 || uv.U < 0 || uv.V < 0 {
		panic("invalid UV coordinates")
	}
	x := int(uv.U * float64(texture.Width))
	y := int(uv.V * float64(texture.Height))
	return texture.data[y*texture.Height+x]
}
