package aoi

const _DEFAULT_AOI_DISTANCE = 100

// XZListAOIManager is an implementation of AOICalculator using XZ lists
type XZListAOIManager struct {
	xSweepList *xAOIList
	zSweepList *yAOIList
}

func NewXZListAOICalculator() AOIManager {
	return &XZListAOIManager{
		xSweepList: newXAOIList(),
		zSweepList: newYAOIList(),
	}
}

// Enter is called when Entity enters Space
func (aoiman *XZListAOIManager) Enter(aoi *AOI, x, y Coord) {
	aoi.x, aoi.y = x, y
	aoiman.xSweepList.Insert(aoi)
	aoiman.zSweepList.Insert(aoi)
	aoiman.adjust(aoi)
}

// Leave is called when Entity leaves Space
func (aoiman *XZListAOIManager) Leave(aoi *AOI) {
	aoiman.xSweepList.Remove(aoi)
	aoiman.zSweepList.Remove(aoi)
}

// Moved is called when Entity moves in Space
func (aoiman *XZListAOIManager) Moved(aoi *AOI, x, y Coord) {
	oldX := aoi.x
	oldY := aoi.y
	aoi.x, aoi.y = x, y
	if oldX != x {
		aoiman.xSweepList.Move(aoi, oldX)
	}
	if oldY != y {
		aoiman.zSweepList.Move(aoi, oldY)
	}
	aoiman.adjust(aoi)
}

// adjust is called by Entity to adjust neighbors
func (aoiman *XZListAOIManager) adjust(aoi *AOI) {
	aoiman.xSweepList.Mark(aoi)
	aoiman.zSweepList.Mark(aoi)
	// AOI marked twice are neighbors
	for neighbor := range aoi.neighbors {
		if neighbor.markVal == 2 {
			// neighbors kept
			neighbor.markVal = -2 // mark this as neighbor
		} else { // markVal < 2
			// was neighbor, but not any more
			aoi.neighbors.Remove(neighbor)
			aoi.Callback.OnLeaveAOI(neighbor)
			neighbor.neighbors.Remove(aoi)
			neighbor.Callback.OnLeaveAOI(aoi)
		}
	}

	// travel in X list again to find all new neighbors, whose markVal == 2
	aoiman.xSweepList.GetClearMarkedNeighbors(aoi)
	// travel in Z list again to unmark all
	aoiman.zSweepList.ClearMark(aoi)
}
