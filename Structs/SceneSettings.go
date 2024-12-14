package Structs

type SceneSettings struct {
	MaxInitialRayDepth int
	MaxMapperRayDepth  int
	PhotonCount        int
	// Optimization structure settings
	KNearestPointRatio float64
}

func NewSceneSettings(iRayDepth, mRayDepth, phCount int) *SceneSettings {
	return &SceneSettings{
		MaxInitialRayDepth: iRayDepth,
		MaxMapperRayDepth:  mRayDepth,
		PhotonCount:        phCount,
		// For each 33-rd triangle there will be a cluster. That means, that average cluster size would be 33 triangles
		KNearestPointRatio: 0.03,
	}
}
