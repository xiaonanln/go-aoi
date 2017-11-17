package aoi

type AOI struct {
	x float32
	y float32
	Callback AOICallback
}

type AOICallback interface {
	OnEnterAOI(aoi *AOI)
	OnLeaveAOI(aoi *AOI)
}

type AOIManager interface {
	Enter(aoi *AOI, x, y float32 )
	Leave(aoi *AOI)
	Moved(aoi *AOI, x, y float32)
}

