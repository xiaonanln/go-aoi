package aoi

type AOISet map[*AOI]struct{}

func (s AOISet) Add(aoi *AOI) {
	s[aoi] = struct{}{}
}

func (s AOISet) Remove(aoi *AOI) {
	delete(s, aoi)
}

func (s AOISet) Contains(aoi *AOI) (ok bool) {
	_, ok = s[aoi]
	return
}
