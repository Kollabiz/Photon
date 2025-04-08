package PhotonMapping

import (
	"Photon/Math"
	"Photon/Structs"
	"Photon/Utils"
	"strconv"
	"time"
)

const (
	MissPoint = -1
)

func PhotonMappingFirstPass(scene *Structs.Scene) *CameraPointCloud {
	camera := scene.GetCamera()
	settings := scene.GetSceneSettings()
	Utils.Log("creating camera point cloud (first pass)")
	cloud := &CameraPointCloud{
		Points:             make([]*CameraPoint, int(camera.GetResolution().U*camera.GetResolution().V)),
		Tree:               nil,
		MaxPointsPerDomain: settings.MaxPointsPerDomain,
	}

	// Starting from top-left (pixel #0), going to bottom right
	Utils.Log("iterating through camera pixels...")
	t := time.Now()
	for y := 0.0; y < camera.GetResolution().V; y++ {
		for x := 0.0; x < camera.GetResolution().U; x++ {
			o, d := camera.GetCameraGrid(Math.Vector2{x, y})
			var prevPoint *CameraPoint
			for i := 0; i < settings.MaxInitialRayDepth; i++ {
				doesIntersect, intersection, barycentric, triangle := Structs.RayCast(d.Normalized(), o, scene)
				if !doesIntersect {
					break
				}
				n := triangle.InterpolateNormals(barycentric)
				p := &CameraPoint{
					Position:  intersection,
					NextPoint: nil,
					I:         d.Normalized(),
					R:         d.Normalized().Reflect(n),
					Triangle:  triangle,
					Bary:      barycentric,
				}
				cloud.AddNonCameraPoint(p)
				o = intersection
				d = d.Normalized().Reflect(n).Normalized()
				// We are storing the photon path reversed, so that during image construction we don't have to create
				// arrays in order to reverse them
				// During construction we will traverse the path from last to first point
				p.NextPoint = prevPoint
				prevPoint = p
			}
			// The point will have an index of y*height+x
			if prevPoint == nil { // There were no intersections with the scene
				cloud.AddPoint(&CameraPoint{
					Position:           o,
					NextPoint:          nil,
					I:                  d.Normalized(),
					R:                  d.Normalized(),
					Triangle:           nil,
					Bary:               Math.Vector2{},
					Color:              Math.Vector3{},
					AccumulatedPhotons: MissPoint,
				}, int(y*camera.GetResolution().U+x))
			} else {
				cloud.AddPoint(prevPoint, int(y*camera.GetResolution().U+x))
			}
		}
	}
	t2 := time.Now()
	Utils.LogSuccess("Camera pass done in " + strconv.FormatFloat(t2.Sub(t).Seconds(), 'f', 2, 64) + " seconds")
	return cloud
}
