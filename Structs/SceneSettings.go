package Structs

type SceneSettings struct {
	MaxInitialRayDepth int
	MaxMapperRayDepth  int
	PhotonCount        int
	// Optimization structure settings
	KNearestPointRatio float64
	PhotonRadius       float64
	MaxPointsPerDomain int
	// Misc
	MinLightEnergy   float64
	AsyncThreads     int
	ViewerUpdateTime int
}

func NewSceneSettings(iRayDepth, mRayDepth, phCount int, phR float64) *SceneSettings {
	return &SceneSettings{
		MaxInitialRayDepth: iRayDepth,
		MaxMapperRayDepth:  mRayDepth,
		PhotonCount:        phCount,
		// For each 33-rd triangle there will be a cluster. That means, that average cluster size would be 33 triangles
		KNearestPointRatio: 0.03,
		PhotonRadius:       phR,
		MaxPointsPerDomain: 64,
		AsyncThreads:       16,
		MinLightEnergy:     0.01,
		ViewerUpdateTime:   1,
	}
}
