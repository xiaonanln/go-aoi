package aoi

// Coord is the type for coordinate axes values
type Coord float32

type AOI struct {
	x    Coord
	y    Coord
	dist Coord
	Data interface{}

	callback AOICallback
	implData interface{}
}

func InitAOI(aoi *AOI, dist Coord, data interface{}, callback AOICallback) {
	aoi.dist = dist
	aoi.Data = data
	aoi.callback = callback
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
