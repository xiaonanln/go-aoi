package aoi

type Coord float32

type AOI struct {
	x        Coord
	y        Coord
	Callback AOICallback

	// Fields for XZListAOIManager
	neighbors    AOISet
	xPrev, xNext *AOI
	yPrev, yNext *AOI
	markVal      int
}

func InitAOI(aoi *AOI, callback AOICallback) {
	aoi.Callback = callback
	aoi.neighbors = make(AOISet)
}

type AOICallback interface {
	OnEnterAOI(other *AOI)
	OnLeaveAOI(other *AOI)
}

type AOIManager interface {
	Enter(aoi *AOI, x, y Coord)
	Leave(aoi *AOI)
	Moved(aoi *AOI, x, y Coord)
}
