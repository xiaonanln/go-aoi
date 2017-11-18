package aoi

type yAOIList struct {
	head *AOI
	tail *AOI
}

func newYAOIList() *yAOIList {
	return &yAOIList{}
}

func (sl *yAOIList) Insert(aoi *AOI) {
	insertCoord := aoi.y
	if sl.head != nil {
		p := sl.head
		for p != nil && p.y < insertCoord {
			p = p.yNext
		}
		// now, p == nil or p.coord >= insertCoord
		if p == nil { // if p == nil, insert AOI at the end of list
			tail := sl.tail
			tail.yNext = aoi
			aoi.yPrev = tail
			sl.tail = aoi
		} else { // otherwise, p >= AOI, insert AOI before p
			prev := p.yPrev
			aoi.yNext = p
			p.yPrev = aoi
			aoi.yPrev = prev

			if prev != nil {
				prev.yNext = aoi
			} else { // p is the head, so AOI should be the new head
				sl.head = aoi
			}
		}
	} else {
		sl.head = aoi
		sl.tail = aoi
	}
}

func (sl *yAOIList) Remove(aoi *AOI) {
	prev := aoi.yPrev
	next := aoi.yNext
	if prev != nil {
		prev.yNext = next
		aoi.yPrev = nil
	} else {
		sl.head = next
	}
	if next != nil {
		next.yPrev = prev
		aoi.yNext = nil
	} else {
		sl.tail = prev
	}
}

func (sl *yAOIList) Move(aoi *AOI, oldCoord Coord) {
	coord := aoi.y
	if coord > oldCoord {
		// moving to next ...
		next := aoi.yNext
		if next == nil || next.y >= coord {
			// no need to adjust in list
			return
		}
		prev := aoi.yPrev
		//fmt.Println(1, prev, next, prev == nil || prev.yNext == AOI)
		if prev != nil {
			prev.yNext = next // remove AOI from list
		} else {
			sl.head = next // AOI is the head, trim it
		}
		next.yPrev = prev

		//fmt.Println(2, prev, next, prev == nil || prev.yNext == next)
		prev, next = next, next.yNext
		for next != nil && next.y < coord {
			prev, next = next, next.yNext
			//fmt.Println(2, prev, next, prev == nil || prev.yNext == next)
		}
		//fmt.Println(3, prev, next)
		// no we have prev.X < coord && (next == nil || next.X >= coord), so insert between prev and next
		prev.yNext = aoi
		aoi.yPrev = prev
		if next != nil {
			next.yPrev = aoi
		} else {
			sl.tail = aoi
		}
		aoi.yNext = next

		//fmt.Println(4)
	} else {
		// moving to prev ...
		prev := aoi.yPrev
		if prev == nil || prev.y <= coord {
			// no need to adjust in list
			return
		}

		next := aoi.yNext
		if next != nil {
			next.yPrev = prev
		} else {
			sl.tail = prev // AOI is the head, trim it
		}
		prev.yNext = next // remove AOI from list

		next, prev = prev, prev.yPrev
		for prev != nil && prev.y > coord {
			next, prev = prev, prev.yPrev
		}
		// no we have next.X > coord && (prev == nil || prev.X <= coord), so insert between prev and next
		next.yPrev = aoi
		aoi.yNext = next
		if prev != nil {
			prev.yNext = aoi
		} else {
			sl.head = aoi
		}
		aoi.yPrev = prev
	}
}

func (sl *yAOIList) Mark(aoi *AOI) {
	prev := aoi.yPrev
	coord := aoi.y

	minCoord := coord - _DEFAULT_AOI_DISTANCE
	for prev != nil && prev.y >= minCoord {
		prev.markVal += 1
		prev = prev.yPrev
	}

	next := aoi.yNext
	maxCoord := coord + _DEFAULT_AOI_DISTANCE
	for next != nil && next.y <= maxCoord {
		next.markVal += 1
		next = next.yNext
	}
}

func (sl *yAOIList) ClearMark(aoi *AOI) {
	prev := aoi.yPrev
	coord := aoi.y

	minCoord := coord - _DEFAULT_AOI_DISTANCE
	for prev != nil && prev.y >= minCoord {
		prev.markVal = 0
		prev = prev.yPrev
	}

	next := aoi.yNext
	maxCoord := coord + _DEFAULT_AOI_DISTANCE
	for next != nil && next.y <= maxCoord {
		next.markVal = 0
		next = next.yNext
	}
}
