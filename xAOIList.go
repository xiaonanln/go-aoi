package aoi

type xAOIList struct {
	head *AOI
	tail *AOI
}

func newXAOIList() *xAOIList {
	return &xAOIList{}
}

func (sl *xAOIList) Insert(aoi *AOI) {
	insertCoord := aoi.x
	if sl.head != nil {
		p := sl.head
		for p != nil && p.x < insertCoord {
			p = p.xNext
		}
		// now, p == nil or p.coord >= insertCoord
		if p == nil { // if p == nil, insert AOI at the end of list
			tail := sl.tail
			tail.xNext = aoi
			aoi.xPrev = tail
			sl.tail = aoi
		} else { // otherwise, p >= AOI, insert AOI before p
			prev := p.xPrev
			aoi.xNext = p
			p.xPrev = aoi
			aoi.xPrev = prev

			if prev != nil {
				prev.xNext = aoi
			} else { // p is the head, so AOI should be the new head
				sl.head = aoi
			}
		}
	} else {
		sl.head = aoi
		sl.tail = aoi
	}
}

func (sl *xAOIList) Remove(aoi *AOI) {
	prev := aoi.xPrev
	next := aoi.xNext
	if prev != nil {
		prev.xNext = next
		aoi.xPrev = nil
	} else {
		sl.head = next
	}
	if next != nil {
		next.xPrev = prev
		aoi.xNext = nil
	} else {
		sl.tail = prev
	}
}

func (sl *xAOIList) Move(aoi *AOI, oldCoord Coord) {
	coord := aoi.x
	if coord > oldCoord {
		// moving to next ...
		next := aoi.xNext
		if next == nil || next.x >= coord {
			// no need to adjust in list
			return
		}
		prev := aoi.xPrev
		//fmt.Println(1, prev, next, prev == nil || prev.xNext == AOI)
		if prev != nil {
			prev.xNext = next // remove AOI from list
		} else {
			sl.head = next // AOI is the head, trim it
		}
		next.xPrev = prev

		//fmt.Println(2, prev, next, prev == nil || prev.xNext == next)
		prev, next = next, next.xNext
		for next != nil && next.x < coord {
			prev, next = next, next.xNext
			//fmt.Println(2, prev, next, prev == nil || prev.xNext == next)
		}
		//fmt.Println(3, prev, next)
		// no we have prev.X < coord && (next == nil || next.X >= coord), so insert between prev and next
		prev.xNext = aoi
		aoi.xPrev = prev
		if next != nil {
			next.xPrev = aoi
		} else {
			sl.tail = aoi
		}
		aoi.xNext = next

		//fmt.Println(4)
	} else {
		// moving to prev ...
		prev := aoi.xPrev
		if prev == nil || prev.x <= coord {
			// no need to adjust in list
			return
		}

		next := aoi.xNext
		if next != nil {
			next.xPrev = prev
		} else {
			sl.tail = prev // AOI is the head, trim it
		}
		prev.xNext = next // remove AOI from list

		next, prev = prev, prev.xPrev
		for prev != nil && prev.x > coord {
			next, prev = prev, prev.xPrev
		}
		// no we have next.X > coord && (prev == nil || prev.X <= coord), so insert between prev and next
		next.xPrev = aoi
		aoi.xNext = next
		if prev != nil {
			prev.xNext = aoi
		} else {
			sl.head = aoi
		}
		aoi.xPrev = prev
	}
}

func (sl *xAOIList) Mark(aoi *AOI) {
	prev := aoi.xPrev
	coord := aoi.x

	minCoord := coord - _DEFAULT_AOI_DISTANCE
	for prev != nil && prev.x >= minCoord {
		prev.markVal += 1
		prev = prev.xPrev
	}

	next := aoi.xNext
	maxCoord := coord + _DEFAULT_AOI_DISTANCE
	for next != nil && next.x <= maxCoord {
		next.markVal += 1
		next = next.xNext
	}
}

func (sl *xAOIList) GetClearMarkedNeighbors(aoi *AOI) {
	prev := aoi.xPrev
	coord := aoi.x
	minCoord := coord - _DEFAULT_AOI_DISTANCE
	for prev != nil && prev.x >= minCoord {
		if prev.markVal == 2 {
			aoi.Callback.OnEnterAOI(prev)
			prev.Callback.OnEnterAOI(aoi)
		}
		prev.markVal = 0
		prev = prev.xPrev
	}

	next := aoi.xNext
	maxCoord := coord + _DEFAULT_AOI_DISTANCE
	for next != nil && next.x <= maxCoord {
		if next.markVal == 2 {
			aoi.Callback.OnEnterAOI(next)
			next.Callback.OnEnterAOI(aoi)
		}
		next.markVal = 0
		next = next.xNext
	}
	return
}
