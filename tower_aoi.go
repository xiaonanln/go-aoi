package aoi

import "log"

type TowerAOIManager struct {
	minX, maxX, minY, maxY Coord
	towerRange             Coord
	towers                 [][]tower
	xTowerNum, yTowerNum   int
}

func (aoiman *TowerAOIManager) Enter(aoi *AOI, x, y Coord) {
	aoi.x, aoi.y = x, y
	obj := &aoiobj{aoi: aoi}
	aoi.implData = obj
	t := aoiman.getTowerXY(x, y)
	t.addObj(obj)

}

func (aoiman *TowerAOIManager) transXY(x, y Coord) (int, int) {
	xi := int((x - aoiman.minX) / aoiman.towerRange)
	yi := int((y - aoiman.minY) / aoiman.towerRange)
	if xi < 0 {
		xi = 0
	} else if xi >= aoiman.xTowerNum {
		xi = aoiman.xTowerNum - 1
	}

	if yi < 0 {
		yi = 0
	} else if yi >= aoiman.yTowerNum {
		yi = aoiman.yTowerNum - 1
	}

	return xi, yi
}

func (aoiman *TowerAOIManager) getTowerXY(x, y Coord) *tower {
	xi, yi := aoiman.transXY(x, y)
	return &aoiman.towers[xi][yi]
}

func (aoiman *TowerAOIManager) Leave(aoi *AOI) {

}

func (aoiman *TowerAOIManager) Moved(aoi *AOI, x, y Coord) {

}

func (aoiman *TowerAOIManager) init() {
	numXSlots := int((aoiman.maxX-aoiman.minX)/aoiman.towerRange) + 1
	aoiman.xTowerNum = numXSlots
	numYSlots := int((aoiman.maxY-aoiman.minY)/aoiman.towerRange) + 1
	aoiman.yTowerNum = numYSlots
	aoiman.towers = make([][]tower, numXSlots)
	for i := 0; i < numXSlots; i++ {
		aoiman.towers[i] = make([]tower, numYSlots)
		for j := 0; j < numYSlots; j++ {
			aoiman.towers[i][j].init()
		}
	}
}

func NewTowerAOIManager(minX, maxX, minY, maxY Coord, towerRange Coord) AOIManager {
	aoiman := &TowerAOIManager{minX: minX, maxX: maxX, minY: minY, maxY: maxY, towerRange: towerRange}
	aoiman.init()

	return aoiman
}

type tower struct {
	objs     map[*aoiobj]struct{}
	watchers map[*aoiobj]struct{}
}

func (t *tower) init() {
	t.objs = map[*aoiobj]struct{}{}
	t.watchers = map[*aoiobj]struct{}{}
}

func (t *tower) addObj(obj *aoiobj) {
	t.objs[obj] = struct{}{}
}

func (t *tower) removeObj(obj *aoiobj) {
	delete(t.objs, obj)
}

func (t *tower) addWatcher(obj *aoiobj) {
	if _, ok := t.watchers[obj]; ok {
		log.Panicf("duplicate add watcher")
	}
	t.watchers[obj] = struct{}{}
	// now obj can see all objs under this tower
	for neighbor := range t.objs {
		obj.aoi.Callback.OnEnterAOI(neighbor.aoi)
	}
}

func (t *tower) removeWatcher(obj *aoiobj) {
	if _, ok := t.watchers[obj]; !ok {
		log.Panicf("duplicate remove watcher")
	}

	delete(t.watchers, obj)
	for neighbor := range t.objs {
		obj.aoi.Callback.OnLeaveAOI(neighbor.aoi)
	}
}

type aoiobj struct {
	aoi *AOI
}
