package PhotonMapping

import (
	"Photon/FileFormats"
	"Photon/Math"
	"Photon/Structs"
	"Photon/Structs/BRDFS"
	"Photon/Utils"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"github.com/Kollabiz/GoColor"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"sync"
	"time"
)

const (
	intSeconds = 1_000_000_000
	// Light source types

	PointLight = 0
	ConeLight  = 1

	BgEnvironment = 0
	BgTransparent = 1
)

type App struct {
	Scene         *Structs.Scene
	CameraCloud   *CameraPointCloud
	threadHandler *PhotonThreadHandler
	env           *Environment
	fyneApp       fyne.App
	raster        *canvas.Raster
	width         int
	height        int
	mtlReader     *FileFormats.MTLParser
}

func NewApp(resolutionX, resolutionY int, fov, photonRadius float64) *App {
	Utils.Log("Creating App instance")
	nApp := &App{}
	sceneSettings := Structs.NewSceneSettings(4, 16, 65536, photonRadius)
	nApp.Scene = Structs.NewScene(resolutionX, resolutionY, fov, sceneSettings)
	nApp.threadHandler = &PhotonThreadHandler{
		maxThreads: nApp.Scene.GetSceneSettings().AsyncThreads,
		busy:       false,
		wg:         sync.WaitGroup{},
	}
	nApp.width = resolutionX
	nApp.height = resolutionY
	Utils.Log("Creating a Fyne app")
	nApp.fyneApp = app.New()

	Utils.Log("Creating a Fyne raster")

	//The raster update function. We want to keep it as simple as possible
	nApp.raster = canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
		point := nApp.CameraCloud.Points[fixScreenCoordinates(x, y, w, h, nApp.width, nApp.height)]
		// Since the photon path is stored in reverse order, we can just use the linked array as is
		if point.AccumulatedPhotons == MissPoint {
			return ldrToneMap(nApp.env.SampleEnvironment(point.I)).ToColor()
		}
		pixelColor := point.Color.FDiv(float64(point.AccumulatedPhotons))
		for point.NextPoint != nil {
			nPoint := point.NextPoint
			// Here we are treating the light 'reflected' from the other point as a light source
			// In fact, the Material's SampleLight function can be used for all kinds of light
			pixelColor = nPoint.Color.FDiv(float64(nPoint.AccumulatedPhotons))
			point = nPoint
		}
		return ldrToneMap(pixelColor).ToColor()
	})

	Utils.LogSuccess("Created Fyne raster")

	Utils.Log("Creating default environment")

	nApp.env = &Environment{
		plainColor: true,
		color:      Math.Vector3{0.33, 0.33, 0.33},
		image:      nil,
	}

	nApp.mtlReader = &FileFormats.MTLParser{Brdf: BRDFS.NewCookTorranceBRDF()}
	nApp.mtlReader.DropTables()

	return nApp
}

// Since Fyne's rasters are created with wrong size (IDK how, but it creates 606x606 rasters inside a 512x512 app.
// We have to do some interpolation, in order not to break the app (thanks Fyne)
func fixScreenCoordinates(x, y, w, h, realW, realH int) int {
	xScaled := int(float64(x) / float64(w) * float64(realW))
	yScaled := int(float64(y) / float64(h) * float64(realH))
	return yScaled*realW + xScaled
}

func ldrToneMap(color Math.Vector3) Math.Vector3 {
	return Math.Vector3{
		X: math.Min(math.Max(math.Sqrt(color.X), 0), 1),
		Y: math.Min(math.Max(math.Sqrt(color.Y), 0), 1),
		Z: math.Min(math.Max(math.Sqrt(color.Z), 0), 1),
	}
}

func (app *App) Run() {
	app.Scene.RebuildBVH()
	app.CameraCloud = PhotonMappingFirstPass(app.Scene)
	app.CameraCloud.ConstructTree()
	win := app.fyneApp.NewWindow("Photon renderer")
	win.Resize(fyne.NewSize(float32(app.width), float32(app.height)))
	win.SetFixedSize(true)
	Utils.Log("Starting async photon mapping")
	app.threadHandler.AllocThreads(app.Scene, app.CameraCloud, app.env)
	win.SetContent(app.raster)
	Utils.Log("Starting async window updater")
	go app.asyncAppUpdate()
	go app.asyncKeyboardListener()
	win.ShowAndRun()
}

func (app *App) asyncKeyboardListener() {
	var input string
	for input != "abort" {
		fmt.Scanln(&input)
		switch input {
		case "pause":
			app.threadHandler.Finish()
		case "export":
			if app.threadHandler.busy {
				GoColor.PrintlnFg256("Cannot export while rendering is in progress. Type 'pause' first", GoColor.LightRed)
				break
			}
			app.exportToImage("Export.png", BgEnvironment)
		case "resume":
			app.threadHandler.AllocThreads(app.Scene, app.CameraCloud, app.env)
			go app.asyncAppUpdate()
		case "abort":
			app.threadHandler.Finish()
			app.fyneApp.Quit()
		}
	}
}

func (app *App) asyncAppUpdate() {
	for app.threadHandler.busy {
		app.CameraCloud.Mu.Lock()
		app.raster.Refresh()
		app.CameraCloud.Mu.Unlock()
		time.Sleep(time.Duration(app.Scene.GetSceneSettings().ViewerUpdateTime * intSeconds))
	}
	// Adding a Refresh() call on exit in case the render finish didn't fit into the raster update intervals
	// (which it most certainly didn't)
	app.raster.Refresh()
}

func (app *App) exportToImage(filename string, backgroundMode int) {
	Utils.Log("exporting rendered image to " + filename)
	fl, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		Utils.LogError("could not open or create file " + filename)
		panic(err)
	}

	defer fl.Close()

	img := image.NewRGBA(image.Rect(0, 0, app.width, app.height))

	for y := 0; y < app.height; y++ {
		for x := 0; x < app.width; x++ {
			var pixelColor Math.Vector3
			point := app.CameraCloud.Points[y*app.width+x]
			// Since the photon path is stored in reverse order, we can just use the linked array as is
			if point.AccumulatedPhotons == MissPoint {
				pixelColor = app.env.SampleEnvironment(point.I)
			} else {
				pixelColor = point.Color.FDiv(float64(point.AccumulatedPhotons))
				for point.NextPoint != nil {
					nPoint := point.NextPoint
					pixelColor = Math.InterpolateVector3(nPoint.Triangle.Material.SampleLight(nPoint.Bary, nPoint.I, nPoint.R.Inverse(),
						nPoint.Triangle.InterpolateNormals(nPoint.Bary), 1, pixelColor), nPoint.Color.FDiv(
						float64(nPoint.AccumulatedPhotons)), 0.5)
					point = nPoint
				}
			}
			img.Set(x, y, ldrToneMap(pixelColor).ToColor())
		}
	}

	err = png.Encode(fl, img)
	if err != nil {
		Utils.LogError("Couldn't export to image")
	}
	Utils.LogSuccess("Finished exporting!")
}

func (app *App) AddMeshesFromFile(filename string) {
	meshes := FileFormats.ReadOBJFile(filename, app.mtlReader)
	for i := 0; i < len(meshes); i++ {
		app.Scene.AddObject(&meshes[i])
	}
}

func (app *App) AddLightSource(lightSourceType int, position Math.Vector3, direction Math.Vector3, color Math.Vector3,
	intensity float64, falloff float64) {
	switch lightSourceType {
	case PointLight:
		app.Scene.AddLightSource(&Structs.PointLight{
			Position:  position,
			Intensity: intensity,
			Color:     color,
		})
		break
	case ConeLight:
		app.Scene.AddLightSource(&Structs.ConeLight{
			Position:  position,
			Direction: direction,
			Intensity: intensity,
			Falloff:   falloff,
			Color:     color,
		})
	}
}

func (app *App) ChangeBRDF(brdf Structs.IBRDF) {
	app.mtlReader.Brdf = brdf
}

func (app *App) SetEnvironmentImage(img string) {
	app.env = NewHDREnvironment(img)
}

func (app *App) SetEnvironmentSimple(color Math.Vector3) {
	app.env = NewPlainEnvironment(color)
}
