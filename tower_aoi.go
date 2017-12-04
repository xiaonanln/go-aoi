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

	aoiman.visitWatchedTowers(x, y, aoi.dist, func(tower *tower) {
		tower.addWatcher(obj)
	})

	t := aoiman.getTowerXY(x, y)
	t.addObj(obj, nil)
}

func (aoiman *TowerAOIManager) Leave(aoi *AOI) {
	obj := aoi.implData.(*aoiobj)
	obj.tower.removeObj(obj, true)

	aoiman.visitWatchedTowers(aoi.x, aoi.y, aoi.dist, func(tower *tower) {
		tower.removeWatcher(obj)
	})
}

func (aoiman *TowerAOIManager) Moved(aoi *AOI, x, y Coord) {
	oldx, oldy := aoi.x, aoi.y
	aoi.x, aoi.y = x, y
	obj := aoi.implData.(*aoiobj)
	t0 := obj.tower
	t1 := aoiman.getTowerXY(x, y)

	if t0 != t1 {
		t0.removeObj(obj, false)
		t1.addObj(obj, t0)
	}

	oximin, oximax, oyimin, oyimax := aoiman.getWatchedTowers(oldx, oldy, aoi.dist)
	ximin, ximax, yimin, yimax := aoiman.getWatchedTowers(x, y, aoi.dist)

	for xi := oximin; xi <= oximax; xi++ {
		for yi := oyimin; yi <= oyimax; yi++ {
			if xi >= ximin && xi <= ximax && yi >= yimin && yi <= yimax {
				continue
			}

			tower := &aoiman.towers[xi][yi]
			tower.removeWatcher(obj)
		}
	}

	for xi := ximin; xi <= ximax; xi++ {
		for yi := yimin; yi <= yimax; yi++ {
			if xi >= oximin && xi <= oximax && yi >= oyimin && yi <= oyimax {
				continue
			}

			tower := &aoiman.towers[xi][yi]
			tower.addWatcher(obj)
		}
	}
}

func (aoiman *TowerAOIManager) transXY(x, y Coord) (int, int) {
	xi := int((x - aoiman.minX) / aoiman.towerRange)
	yi := int((y - aoiman.minY) / aoiman.towerRange)
	return aoiman.normalizeXi(xi), aoiman.normalizeYi(yi)
}

func (aoiman *TowerAOIManager) normalizeXi(xi int) int {
	if xi < 0 {
		xi = 0
	} else if xi >= aoiman.xTowerNum {
		xi = aoiman.xTowerNum - 1
	}
	return xi
}

func (aoiman *TowerAOIManager) normalizeYi(yi int) int {
	if yi < 0 {
		yi = 0
	} else if yi >= aoiman.yTowerNum {
		yi = aoiman.yTowerNum - 1
	}
	return yi
}

func (aoiman *TowerAOIManager) getTowerXY(x, y Coord) *tower {
	xi, yi := aoiman.transXY(x, y)
	return &aoiman.towers[xi][yi]
}

func (aoiman *TowerAOIManager) getWatchedTowers(x, y Coord, aoiDistance Coord) (int, int, int, int) {
	ximin, yimin := aoiman.transXY(x-aoiDistance, y-aoiDistance)
	ximax, yimax := aoiman.transXY(x+aoiDistance, y+aoiDistance)
	//aoiTowerNum := int(aoiDistance/aoiman.towerRange) + 1
	//ximid, yimid := aoiman.transXY(x, y)
	//ximin, ximax := aoiman.normalizeXi(ximid-aoiTowerNum), aoiman.normalizeXi(ximid+aoiTowerNum)
	//yimin, yimax := aoiman.normalizeYi(yimid-aoiTowerNum), aoiman.normalizeYi(yimid+aoiTowerNum)
	return ximin, ximax, yimin, yimax
}

func (aoiman *TowerAOIManager) visitWatchedTowers(x, y Coord, aoiDistance Coord, f func(*tower)) {
	ximin, ximax, yimin, yimax := aoiman.getWatchedTowers(x, y, aoiDistance)
	for xi := ximin; xi <= ximax; xi++ {
		for yi := yimin; yi <= yimax; yi++ {
			tower := &aoiman.towers[xi][yi]
			f(tower)
		}
	}
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

func (t *tower) addObj(obj *aoiobj, fromOtherTower *tower) {
	obj.tower = t
	t.objs[obj] = struct{}{}
	if fromOtherTower == nil {
		for watcher := range t.watchers {
			if watcher == obj {
				continue
			}
			watcher.aoi.callback.OnEnterAOI(obj.aoi)
		}
	} else {
		// obj moved from other tower to this tower
		for watcher := range fromOtherTower.watchers {
			if watcher == obj {
				continue
			}
			if _, ok := t.watchers[watcher]; ok {
				continue
			}
			watcher.aoi.callback.OnLeaveAOI(obj.aoi)
		}
		for watcher := range t.watchers {
			if watcher == obj {
				continue
			}
			if _, ok := fromOtherTower.watchers[watcher]; ok {
				continue
			}
			watcher.aoi.callback.OnEnterAOI(obj.aoi)
		}
	}
}

func (t *tower) removeObj(obj *aoiobj, notifyWatchers bool) {
	obj.tower = nil
	delete(t.objs, obj)
	if notifyWatchers {
		for watcher := range t.watchers {
			if watcher == obj {
				continue
			}
			watcher.aoi.callback.OnLeaveAOI(obj.aoi)
		}
	}
}

func (t *tower) addWatcher(obj *aoiobj) {
	if _, ok := t.watchers[obj]; ok {
		log.Panicf("duplicate add watcher")
	}
	t.watchers[obj] = struct{}{}
	// now obj can see all objs under this tower
	for neighbor := range t.objs {
		if neighbor == obj {
			continue
		}
		obj.aoi.callback.OnEnterAOI(neighbor.aoi)
	}
}

func (t *tower) removeWatcher(obj *aoiobj) {
	if _, ok := t.watchers[obj]; !ok {
		log.Panicf("duplicate remove watcher")
	}

	delete(t.watchers, obj)
	for neighbor := range t.objs {
		if neighbor == obj {
			continue
		}
		obj.aoi.callback.OnLeaveAOI(neighbor.aoi)
	}
}

type aoiobj struct {
	aoi   *AOI
	tower *tower
}
