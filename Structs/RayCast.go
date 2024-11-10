package Structs

import (
	"Photon/Math"
	"Photon/Math/Mesh"
	"Photon/Math/Ray"
)

func RayCast(rDirection, rOrigin Math.Vector3, scene *Scene) (bool, Math.Vector3, Math.Vector2, *Mesh.Triangle) {
	// First, we have to iterate through the BVH tree
	currentNode := scene.baseNode
	// If the ray doesn't intersect scene's base node, we just cut it off early
	if !Ray.IntersectRayAABB(rDirection, rOrigin, currentNode.AABB) {
		return false, Math.ZeroVector3(), Math.ZeroVector2(), nil
	}
	for currentNode.Mesh == nil && len(currentNode.Children) > 0 {
		anyIntersected := false
		if Ray.IntersectRayAABB(rDirection, rOrigin, currentNode.Children[0].AABB) {
			currentNode = &currentNode.Children[0]
			anyIntersected = true
		} else if Ray.IntersectRayAABB(rDirection, rOrigin, currentNode.Children[1].AABB) {
			currentNode = &currentNode.Children[1]
			anyIntersected = true
		}
		if !anyIntersected {
			return false, Math.ZeroVector3(), Math.ZeroVector2(), nil
		}
	}
	// Now we can finally intersect the mesh
	var intersectedTriangle *Mesh.Triangle
	var intersectPos Math.Vector3
	var intersectUV Math.Vector2
	for i := 0; i < len(currentNode.Mesh.Triangles); i++ {
		tri := &currentNode.Mesh.Triangles[i]
		// Backface culling
		if tri.TriangleNormal.Dot(rDirection) < 0 {
			continue
		}
		intersects, position, uv := Ray.IntersectRayTriangle(rDirection, rOrigin, tri)
		// First iteration crutch
		if intersectedTriangle == nil {
			intersectedTriangle = tri
			intersectPos = position
			intersectUV = uv
			continue
		}
		if !intersects {
			continue
		}
		// We only want the closest intersection
		if rOrigin.Sub(intersectPos).LenSq() > rOrigin.Sub(position).LenSq() {
			intersectedTriangle = tri
			intersectPos = position
			intersectUV = uv
		}
	}
	if intersectedTriangle == nil {
		return false, Math.ZeroVector3(), Math.ZeroVector2(), nil
	}
	return true, intersectPos, intersectUV, intersectedTriangle
}
