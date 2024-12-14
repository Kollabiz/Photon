package Structs

import (
	"Photon/Math"
	"Photon/Math/BoundingVolumes"
	"Photon/Math/Mesh"
	"Photon/Math/Ray"
)

func RayCast(rDirection, rOrigin Math.Vector3, scene *Scene) (bool, Math.Vector3, Math.Vector2, *Mesh.Triangle) {
	// First, we iterate through the BVH tree to find possible intersections (there might be more than one)
	leftoverNodes := []BoundingVolumes.BVHNode{*scene.baseNode}
	var BVHIntersections []BoundingVolumes.BVHNode
	for len(leftoverNodes) > 0 {
		node := &leftoverNodes[0]
		leftoverNodes = leftoverNodes[1:]
		if Ray.IntersectRayAABB(rDirection, rOrigin, node.AABB) {
			// If the intersected node is a mesh
			if node.Mesh != nil {
				BVHIntersections = append(BVHIntersections, *node)
			} else {
				if node.Child1 != nil && node.Child2 != nil {
					leftoverNodes = append(leftoverNodes, *node.Child1)
					leftoverNodes = append(leftoverNodes, *node.Child2)
				}
			}
		}
	}

	if len(BVHIntersections) == 0 {
		// There are no BVH intersections, we can just bail early
		return false, Math.ZeroVector3(), Math.ZeroVector2(), nil
	}

	// Now we iterate through all the intersected nodes and intersect triangles
	// We can just sort the K-Nearest points by distance. Then finding the first intersection would be enough
	// It's already guaranteed to be the closest one
	var nearestTriangle *Mesh.Triangle = nil
	var closestIntersectionPoint Math.Vector3
	var closestIntersectionBarycentricPoint Math.Vector2
	for i := 0; i < len(BVHIntersections); i++ {
		node := BVHIntersections[i]
		// We iterate through all the clusters in a node and for ones that intersect the ray we intersect triangles
		// Then we compare the resulting intersected triangles to nearestTriangle
		for j := 0; j < len(node.TriangleClusters); j++ {
			if Ray.IntersectRayAABB(rDirection, rOrigin, node.TriangleClusters[j].AABB) {
				// The ray does intersect the cluster AABB, so we iterate through all the triangles and intersect them
				for k := 0; k < len(node.TriangleClusters[j].Triangles); k++ {
					tri := &node.TriangleClusters[j].Triangles[k]
					doesIntersect, intersectionPoint, intersectionBarycentricPoint := Ray.IntersectRayTriangle(rDirection, rOrigin, tri)
					if doesIntersect && (nearestTriangle == nil || closestIntersectionPoint.Sub(rOrigin).LenSq() > intersectionPoint.Sub(rOrigin).LenSq()) {
						nearestTriangle = tri
						closestIntersectionPoint = intersectionPoint
						closestIntersectionBarycentricPoint = intersectionBarycentricPoint
					}
				}
			}
		}
	}

	// Done! Now we can return the intersection point (or lack thereof)
	return nearestTriangle != nil, closestIntersectionPoint, closestIntersectionBarycentricPoint, nearestTriangle
}
