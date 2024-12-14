package Mesh

import "Photon/Math/BoundingVolumes"

type TriangleCluster struct {
	AABB      *BoundingVolumes.AABoundingBox
	Triangles []Triangle
}
