package PhotonMapping

import (
	"Photon/Math"
	"Photon/Structs"
	"Photon/Utils"
	"github.com/mdouchement/hdr"
	_ "github.com/mdouchement/hdr/codec/rgbe"
	"image"
	"math"
	"os"
)

func readHDRImage(path string) *Structs.TextureRGB {
	Utils.Log("reading HDR image \"" + path + "\"")
	fl, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer fl.Close()

	im, _, err := image.Decode(fl)
	if err != nil {
		panic(err)
	}
	hdrIm := im.(hdr.Image)

	tex := Structs.TextureRGBFromHDR(hdrIm)

	return tex
}

type Environment struct {
	plainColor bool
	color      Math.Vector3
	image      *Structs.TextureRGB
}

func NewHDREnvironment(hdrImage string) *Environment {
	env := &Environment{
		plainColor: false,
		color:      Math.Vector3{},
		image:      readHDRImage(hdrImage),
	}
	return env
}

func NewPlainEnvironment(color Math.Vector3) *Environment {
	return &Environment{
		plainColor: true,
		color:      color,
		image:      nil,
	}
}

func (env *Environment) SampleEnvironment(direction Math.Vector3) Math.Vector3 {
	if env.plainColor {
		return env.color
	}

	azimuth := math.Atan2(direction.Dot(Math.Vector3{Y: -1}), direction.Dot(Math.Vector3{X: 1}))/math.Pi + 1
	elevation := (direction.Dot(Math.Vector3{Z: -1}) + 1) / 2
	return env.image.At(Math.Vector2{azimuth / 2, elevation})
}
