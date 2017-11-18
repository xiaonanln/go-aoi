package aoi

type Coord float32

type AOI struct {
	x Coord
	y Coord
	Callback AOICallback

	// Fields for XZListAOIManager
	neighbors map[*AOI] struct{}
	xPrev, xNext *AOI
	yPrev, yNext *AOI
	markVal int
}

func InitAOI(aoi *AOI, callback AOICallback) {
	aoi.neighbors = make(map[*AOI] struct{})
	aoi.Callback = callback
}

type AOICallback interface {
	OnEnterAOI(aoi *AOI)
	OnLeaveAOI(aoi *AOI)
}

type AOIManager interface {
	Enter(aoi *AOI, x, y Coord )
	Leave(aoi *AOI)
	Moved(aoi *AOI, x, y Coord)
}
