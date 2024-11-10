package Structs

type SceneSettings struct {
	MaxInitialRayDepth int
	MaxMapperRayDepth  int
	PhotonCount        int
}

func NewSceneSettings(iRayDepth, mRayDepth, phCount int) *SceneSettings {
	return &SceneSettings{
		MaxInitialRayDepth: iRayDepth,
		MaxMapperRayDepth:  mRayDepth,
		PhotonCount:        phCount,
	}
}
