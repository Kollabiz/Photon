package PhotonMapping

import "sync"

type CameraPointCloud struct {
	Points             []*CameraPoint
	NonCameraPoints    []*CameraPoint
	Tree               *KDTreeSpace
	MaxPointsPerDomain int
	Mu                 sync.Mutex
}

func (cloud *CameraPointCloud) AddPoint(point *CameraPoint, i int) {
	cloud.Points[i] = point
}

func (cloud *CameraPointCloud) AddNonCameraPoint(point *CameraPoint) {
	cloud.NonCameraPoints = append(cloud.NonCameraPoints, point)
}

func (cloud *CameraPointCloud) ConstructTree() {
	cloud.Tree = ConstructKDTree(cloud.NonCameraPoints, cloud.MaxPointsPerDomain)
}
