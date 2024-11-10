package Structs

import (
	"Photon/Math"
	"Photon/Utils"
	"image"
	"image/color"
)

type RenderBuffer struct {
	Width  int
	Height int
	Image  *image.RGBA64
}

func MakeRenderBuffer(width, height int) *RenderBuffer {
	Utils.Log("creating a RenderBuffer...")
	buff := &RenderBuffer{}
	buff.Width = width
	buff.Height = height
	buff.Image = image.NewRGBA64(image.Rect(0, 0, width, height))
	Utils.LogSuccess("created a RenderBuffer!")
	return buff
}

func (buffer *RenderBuffer) Point(clr Math.Vector3, x, y int) {
	if x < 0 || x >= buffer.Width || y < 0 || y >= buffer.Height {
		Utils.LogError("trying to draw on RenderBuffer outside of its bounds")
		panic("pixel out of bounds when drawing onto RenderBuffer")
	}
	buffer.Image.SetRGBA64(x, y, color.RGBA64{
		R: uint16(clr.X * 65535),
		G: uint16(clr.Y * 65535),
		B: uint16(clr.Z * 65535),
		A: 0,
	})
}

func (buffer *RenderBuffer) At(x, y int) color.RGBA64 {
	if x < 0 || x >= buffer.Width || y < 0 || y >= buffer.Height {
		Utils.LogError("trying to index RenderBuffer outside of its bounds")
		panic("pixel out of bounds when reading RenderBuffer")
	}
	r, g, b, _ := buffer.Image.At(x, y).RGBA()
	return color.RGBA64{
		R: uint16(r / 65537),
		G: uint16(g / 65537),
		B: uint16(b / 65537),
		A: 0,
	}
}

func (buffer *RenderBuffer) Reset() {
	Utils.Log("resetting the RenderBuffer...")
	for y := 0; y < buffer.Height; y++ {
		for x := 0; x < buffer.Width; x++ {
			buffer.Image.SetRGBA64(x, y, color.RGBA64{0, 0, 0, 0})
		}
	}
	Utils.LogSuccess("RenderBuffer reset!")
}
