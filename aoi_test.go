package aoi

import "testing"

func TestXZListAOIManager(t *testing.T) {
	testAOI(NewXZListAOICalculator())
}

func testAOI(aoiman AOIManager, numAOI int) {
	aois := []*AOI{}
	for i:=0;i<numAOI;i++ {
		aoi:=&AOI{}
		InitAOI(aoi)
		aois = append(aois, aoi)
	}
}