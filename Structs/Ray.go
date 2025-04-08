package Structs

import (
	"Photon/Math"
	"math"
)

const epsilon = 0.00001

func IntersectRayTriangle(rDirection, rOrigin Math.Vector3, tri *Triangle) (hit bool, intersectionPoint Math.Vector3, barycentricIntersection Math.Vector2) {
	var h, s, q Math.Vector3
	e1 := tri.Edge12()
	e2 := tri.Edge13()
	var a, f, u, v float64
	h = rDirection.Cross(e2)
	a = e1.Dot(h)

	if a > -epsilon && a < epsilon {
		return false, Math.ZeroVector3(), Math.ZeroVector2()
	}

	f = 1 / a
	s = rOrigin.Sub(tri.V1Pos)
	u = f * s.Dot(h)

	if u < 0 || u > 1 { // U + V must be less than 1
		return false, Math.ZeroVector3(), Math.ZeroVector2()
	}

	q = s.Cross(e1)
	v = f * rDirection.Dot(q)

	if v < 0 || u+v > 1 {
		return false, Math.ZeroVector3(), Math.ZeroVector2()
	}

	t := f * e2.Dot(q)

	if t <= epsilon {
		return false, Math.ZeroVector3(), Math.ZeroVector2()
	}
	p := rOrigin.Add(rDirection.FMul(t))
	return true, p, Math.Vector2{u, v}
}

func IntersectRayAABB(rDirection, rOrigin Math.Vector3, aabb *AABoundingBox) bool {
	dirFrac := Math.Vector3{
		X: 1.0 / rDirection.X,
		Y: 1.0 / rDirection.Y,
		Z: 1.0 / rDirection.Z,
	}
	t1 := (aabb.Point1.X - rOrigin.X) * dirFrac.X
	t2 := (aabb.Point2.X - rOrigin.X) * dirFrac.X
	t3 := (aabb.Point1.Y - rOrigin.Y) * dirFrac.Y
	t4 := (aabb.Point2.Y - rOrigin.Y) * dirFrac.Y
	t5 := (aabb.Point1.Z - rOrigin.Z) * dirFrac.Z
	t6 := (aabb.Point2.Z - rOrigin.Z) * dirFrac.Z

	tmin := math.Max(math.Max(math.Min(t1, t2), math.Min(t3, t4)), math.Min(t5, t6))
	tmax := math.Min(math.Min(math.Max(t1, t2), math.Max(t3, t4)), math.Max(t5, t6))

	if tmax < 0 || tmax < tmin {
		return false
	}

	return true
}
