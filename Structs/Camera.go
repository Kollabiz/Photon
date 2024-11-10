package Structs

import (
	"Photon/Math"
	"math"
)

type Camera struct {
	Transform *Math.Transform
	focalLength float64
	lensSize Math.Vector2
	resolution Math.Vector2
}

func lensSizeFromResolution(resolution Math.Vector2) Math.Vector2 {
	if resolution.U > resolution.V {
		aspectRatio := resolution.U / resolution.V
		return Math.Vector2{
			U: 1,
			V: aspectRatio,
		}
	} else {
		aspectRatio := resolution.V / resolution.U
		return Math.Vector2{
			U: aspectRatio,
			V: 1,
		}
	}
}

func focalLengthFromFOV(fov float64) float64 {
	return math.Tan(fov/180*math.Pi)
}

func NewCamera(position, rotation Math.Vector3, resolution Math.Vector2, fov float64) *Camera {
	c := &Camera{}
	c.Transform = Math.NewTransform(position, rotation, Math.Vector3{1,1,1})
	c.lensSize = lensSizeFromResolution(resolution)
	c.focalLength = focalLengthFromFOV(fov)
	c.resolution = resolution
	return c
}

func (c *Camera) GetCameraGrid(uv Math.Vector2) (pointPos Math.Vector3, direction Math.Vector3) {
	dX := c.lensSize.U / c.resolution.U
	dY := c.lensSize.V / c.resolution.V
	point := Math.Vector3{
		X: uv.U*dX,
		Y: uv.V*dY,
		Z: 0,
	}
	focalP := Math.Vector3{
		X: 0,
		Y: 0,
		Z: -c.focalLength,
	}
	point = c.Transform.GetRotationMatrix().VecMul(point)
	return point.Add(c.Transform.GetPosition()), c.Transform.GetRotationMatrix().VecMul(focalP).Sub(point)
}
