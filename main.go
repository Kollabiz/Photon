package main

import (
	"Photon/App/PhotonMapping"
	"Photon/Math"
	"Photon/Utils"
	"flag"
	"strconv"
	"strings"
)

type ResolutionFlag struct {
	Width  int
	Height int
}

func (r *ResolutionFlag) String() string {
	return "{" + strconv.Itoa(r.Width) + ";" + strconv.Itoa(r.Height) + "}"
}

func (r *ResolutionFlag) Set(s string) error {
	comp := strings.Split(s, ";")
	if len(comp) < 2 {
		panic("invalid resolution. The resolution value must be specified in X;Y format")
	}
	w, err := strconv.Atoi(comp[0])
	h, err := strconv.Atoi(comp[1])
	r.Width = w
	r.Height = h
	return err
}

func main() {
	Utils.Log("Starting...")
	modelFile := flag.String("model", "", "model allows you to specify a path to an .obj file (all .mtl files must be in the same directory!)")
	envImage := flag.String("env", "", "env allows you to specify an .hdr image to use as environment texture")
	var resolution *ResolutionFlag = &ResolutionFlag{0, 0}
	flag.Var(resolution, "res", "res allows you to specify the image (and window) resolution")
	pitch := flag.Float64("pitch", 45, "pitch allows you to specify the pitch angle of the camera in degrees")
	yaw := flag.Float64("yaw", 180, "yaw allows you to specify the yaw angle of the camera in degrees")
	rad := flag.Float64("rad", 1, "rad allows you to specify the distance of the camera from the {0;0;0}")
	fov := flag.Float64("fov", 39.6, "fov allows you to specify the camera's FOV in degrees")
	phRad := flag.Float64("phrad", 0.01, "phrad allows you to specify the photon radius in units")
	flag.Parse()

	if resolution.Width == 0 || resolution.Height == 0 {
		panic("invalid resolution")
	}
	app := PhotonMapping.NewApp(resolution.Width, resolution.Height, *fov, *phRad)
	if *modelFile == "" {
		panic("no model file specified")
	}
	app.AddMeshesFromFile(*modelFile)
	cam := app.Scene.GetCamera()
	cam.Move(Math.Mat3ZRotation(Math.DegToRad(*yaw)).VecMul(
		Math.Mat3XRotation(Math.DegToRad(*pitch)).VecMul(Math.Vector3{
			Z: *rad,
		})))
	cam.SetRotation(Math.Vector3{
		X: Math.DegToRad(*pitch),
		Y: 0,
		Z: Math.DegToRad(*yaw),
	})
	if *envImage == "" {
		Utils.LogWarning("No environment image specified. Using plain environment")
		app.SetEnvironmentSimple(Math.Vector3{
			X: 0.1,
			Y: 0.1,
			Z: 0.1,
		})
	} else {
		app.SetEnvironmentImage(*envImage)
	}
	Utils.Log("Starting the app")
	app.Run()
}
