package Structs

import (
	"Photon/Math"
)

func RayCast(rDirection, rOrigin Math.Vector3, scene *Scene) (bool, Math.Vector3, Math.Vector2, *Triangle) {
	// First, we iterate through the BVH tree to find possible intersections (there might be more than one)
	leftoverNodes := []BVHNode{*scene.baseNode}
	var BVHIntersections []BVHNode
	for len(leftoverNodes) > 0 {
		node := leftoverNodes[0]
		leftoverNodes = leftoverNodes[1:]
		if IntersectRayAABB(rDirection, rOrigin, node.AABB) {
			// If the intersected node is a mesh
			if node.Mesh != nil {
				BVHIntersections = append(BVHIntersections, node)
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
	var nearestTriangle *Triangle = nil
	var closestIntersectionPoint Math.Vector3
	var closestIntersectionBarycentricPoint Math.Vector2
	for i := 0; i < len(BVHIntersections); i++ {
		node := BVHIntersections[i]
		// We iterate through all the clusters in a node and for ones that intersect the ray we intersect triangles
		// Then we compare the resulting intersected triangles to nearestTriangle
		for j := 0; j < len(node.TriangleClusters); j++ {
			if IntersectRayAABB(rDirection, rOrigin, node.TriangleClusters[j].AABB) {
				// The ray does intersect the cluster AABB, so we iterate through all the triangles and intersect them
				for k := 0; k < len(node.TriangleClusters[j].Triangles); k++ {
					tri := &node.TriangleClusters[j].Triangles[k]
					doesIntersect, intersectionPoint, intersectionBarycentricPoint := IntersectRayTriangle(rDirection, rOrigin, tri)
					if intersectionPoint.Sub(rOrigin).LenSq() <= 0.0001 { // The intersection point is too close and
						// is probably an intersection of the triangle the point lies on
						continue
					}
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
