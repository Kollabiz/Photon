package PhotonMapping

import (
	"Photon/Math"
	"Photon/Structs"
	"Photon/Utils"
	"math"
	"math/rand"
	"strconv"
	"time"
)

const (
	Epsilon           = 0.999
	LightSourcePhoton = 0
)

func addPhotonToAPoint(photonColor Math.Vector3, rayDir Math.Vector3, point *CameraPoint) {
	weight := point.Triangle.Material.SampleLight(point.Bary, point.I, rayDir,
		point.Triangle.InterpolateNormals(point.Bary), 1, photonColor)
	point.Color = point.Color.Add(weight)
}

func rayAbsorptionDice(rayColor Math.Vector3, randGen *rand.Rand) bool {
	// This function models the ray absorption. The dimmer the ray's color (the less energy it has), the higher the
	// chance it gets absorbed
	return randGen.Float64() > rayColor.ColorGrayscale()
}

// That function should be called in a separate thread

func AsyncPhotonCast(scene *Structs.Scene, env *Environment, pointCloud *CameraPointCloud, firstLightSource int,
	handler *PhotonThreadHandler) {
	randGen := rand.New(rand.NewSource(time.Now().UnixMilli()))
	if len(pointCloud.NonCameraPoints) == 0 {
		handler.wg.Done()
		return
	}
	i := 0
	l := firstLightSource
	lights := scene.GetLightSources()
	settings := scene.GetSceneSettings()
	envWindowSize := float64(settings.MaxPointsPerDomain) * 8
	envWindow := 0.0
	var rayOrigin Math.Vector3
	var rayDirection Math.Vector3
	var rayColor Math.Vector3
	var tri *Structs.Triangle
	var bary Math.Vector2

	for handler.busy {
		if i == LightSourcePhoton { // The current photon is cast from a light source
			if len(lights) == 0 {
				i = (i + 1) % 2
				continue
			}
			light := lights[l]
			l = (l + 1) % len(lights)
			rayOrigin = light.GetPosition()
			rayDirection = light.GetRandomPoint(randGen)
			rayColor = light.GetLightColor().FMul(light.GetLightIntensityInDirection(rayDirection))
		} else { // The current photon is cast from the environment

			// Environment photon casting is based on the fact that light paths are symmetrical
			// That means that instead of casting random photons from the environment itself, we can cast a random
			// ray from the point cloud, thus simplifying the algorithm and saving some processing power

			// Selecting a random point in the cloud
			randNum := randGen.NormFloat64()
			rnd := int(math.Abs(randNum)*envWindowSize+envWindow) % len(pointCloud.NonCameraPoints)
			randGen.Seed(time.Now().UnixMilli())
			point := pointCloud.NonCameraPoints[rnd]
			envWindow = float64(rnd % len(pointCloud.NonCameraPoints))
			tri = point.Triangle
			bary = point.Bary

			// Point's normal
			normal := point.Triangle.InterpolateNormals(bary)
			// Selecting a random direction on a hemisphere (with its pole parallel to the point's normal)
			roughness := 1 - math.Pow(tri.Material.GetRoughness(tri.InterpolateTexcoords(bary)), 2)
			refl := point.I.Reflect(normal)
			rayDirection = Utils.RandomPointOnHemisphere(randGen).FromSingleVectorBasis(normal.Inverse())
			rayDirection = Math.InterpolateVector3(rayDirection, refl, roughness)
			rayOrigin = point.Position
			// We only need to know whether we hit anything or not, all other data is irrelevant
			hit, _, _, _ := Structs.RayCast(rayDirection, rayOrigin, scene)
			n := pointCloud.Tree.LocateNeighborPoints(point.Position, settings.PhotonRadius)
			rayColor = env.SampleEnvironment(rayDirection)
			// Adding the environment photon to the neighboring points
			pointCloud.Mu.Lock()
			for j := 0; j < len(n.Points); j++ {
				if n.Points[j].Position.Sub(rayOrigin).Len() > settings.PhotonRadius {
					continue
				}
				if !hit {
					addPhotonToAPoint(rayColor, rayDirection.Inverse(), n.Points[j])
				}
				n.Points[j].AccumulatedPhotons += 1
			}
			pointCloud.Mu.Unlock()
			rayRefl := rayDirection.Inverse().Reflect(normal)
			rayColor = point.Triangle.Material.SampleLight(point.Bary, rayRefl, rayDirection,
				point.Triangle.InterpolateNormals(point.Bary), 1, rayColor)
			// Ray direction for the next photon is a reflection of the current direction's inverse over the point's
			// normal
			// And the ray origin stays the same
			rayDirection = rayRefl
			handler.phCount.Add(1)
		}

		// Now we can proceed with normal photon casting
		// Environment photons are practically the same, except for the first 'bounce', since we have to sample the
		// environment texture first

		for n := settings.MaxMapperRayDepth; n >= 0 && !rayAbsorptionDice(rayColor, randGen); n-- {
			// Cast a ray from the previously selected point
			hit, pos, nBary, nTri := Structs.RayCast(rayDirection, rayOrigin, scene)
			if !hit {
				break
			} else {
				normal := nTri.InterpolateNormals(nBary)
				// In case the ray did hit something, locate all the nearest points to the hit position
				// And add this photon to them
				neighbors := pointCloud.Tree.LocateNeighborPoints(pos, settings.PhotonRadius)
				pointCloud.Mu.Lock()
				for p := 0; p < len(neighbors.Points); p++ {
					if neighbors.Points[p].Position.Sub(pos).Len() > settings.PhotonRadius {
						continue
					}
					addPhotonToAPoint(rayColor, rayDirection, neighbors.Points[p])
					neighbors.Points[p].AccumulatedPhotons += 1
				}
				pointCloud.Mu.Unlock()
				rayRefl := rayDirection.Reflect(normal)
				rayColor = tri.Material.SampleLight(bary, rayRefl.Inverse(), rayDirection, normal, 1, rayColor)
				tri = nTri
				bary = nBary
				rayOrigin = pos
				rayDirection = rayRefl
				handler.phCount.Add(1)
			}
		}
		i = (i + 1) % 2
	}
	Utils.LogSuccess("Exiting thread #" + strconv.Itoa(firstLightSource))
	handler.wg.Done()
}
