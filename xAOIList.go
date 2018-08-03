package aoi

type xAOIList struct {
	aoidist Coord
	head    *xzaoi
	tail    *xzaoi
}

func newXAOIList(aoidist Coord) *xAOIList {
	return &xAOIList{aoidist: aoidist}
}

func (sl *xAOIList) Insert(aoi *xzaoi) {
	insertCoord := aoi.aoi.x
	if sl.head != nil {
		p := sl.head
		for p != nil && p.aoi.x < insertCoord {
			p = p.xNext
		}
		// now, p == nil or p.coord >= insertCoord
		if p == nil { // if p == nil, insert xzaoi at the end of list
			tail := sl.tail
			tail.xNext = aoi
			aoi.xPrev = tail
			sl.tail = aoi
		} else { // otherwise, p >= xzaoi, insert xzaoi before p
			prev := p.xPrev
			aoi.xNext = p
			p.xPrev = aoi
			aoi.xPrev = prev

			if prev != nil {
				prev.xNext = aoi
			} else { // p is the head, so xzaoi should be the new head
				sl.head = aoi
			}
		}
	} else {
		sl.head = aoi
		sl.tail = aoi
	}
}

func (sl *xAOIList) Remove(aoi *xzaoi) {
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

func (sl *xAOIList) Move(aoi *xzaoi, oldCoord Coord) {
	coord := aoi.aoi.x
	if coord > oldCoord {
		// moving to next ...
		next := aoi.xNext
		if next == nil || next.aoi.x >= coord {
			// no need to adjust in list
			return
		}
		prev := aoi.xPrev
		//fmt.Println(1, prev, next, prev == nil || prev.xNext == xzaoi)
		if prev != nil {
			prev.xNext = next // remove xzaoi from list
		} else {
			sl.head = next // xzaoi is the head, trim it
		}
		next.xPrev = prev

		//fmt.Println(2, prev, next, prev == nil || prev.xNext == next)
		prev, next = next, next.xNext
		for next != nil && next.aoi.x < coord {
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
		if prev == nil || prev.aoi.x <= coord {
			// no need to adjust in list
			return
		}

		next := aoi.xNext
		if next != nil {
			next.xPrev = prev
		} else {
			sl.tail = prev // xzaoi is the head, trim it
		}
		prev.xNext = next // remove xzaoi from list

		next, prev = prev, prev.xPrev
		for prev != nil && prev.aoi.x > coord {
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

func (sl *xAOIList) Mark(aoi *xzaoi) {
	prev := aoi.xPrev
	coord := aoi.aoi.x

	minCoord := coord - sl.aoidist
	for prev != nil && prev.aoi.x >= minCoord {
		prev.markVal += 1
		prev = prev.xPrev
	}

	next := aoi.xNext
	maxCoord := coord + sl.aoidist
	for next != nil && next.aoi.x <= maxCoord {
		next.markVal += 1
		next = next.xNext
	}
}

func (sl *xAOIList) GetClearMarkedNeighbors(aoi *xzaoi) {
	prev := aoi.xPrev
	coord := aoi.aoi.x
	minCoord := coord - sl.aoidist
	for prev != nil && prev.aoi.x >= minCoord {
		if prev.markVal == 2 {
			aoi.neighbors[prev] = struct{}{}
			aoi.aoi.callback.OnEnterAOI(prev.aoi)
			prev.neighbors[aoi] = struct{}{}
			prev.aoi.callback.OnEnterAOI(aoi.aoi)
		}
		prev.markVal = 0
		prev = prev.xPrev
	}

	next := aoi.xNext
	maxCoord := coord + sl.aoidist
	for next != nil && next.aoi.x <= maxCoord {
		if next.markVal == 2 {
			aoi.neighbors[next] = struct{}{}
			aoi.aoi.callback.OnEnterAOI(next.aoi)
			next.neighbors[aoi] = struct{}{}
			next.aoi.callback.OnEnterAOI(aoi.aoi)
		}
		next.markVal = 0
		next = next.xNext
	}
	return
}
