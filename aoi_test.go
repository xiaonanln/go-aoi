package aoi

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
	"unsafe"
)

const (
	MIN_X = -500
	MAX_X = 500
	MIN_Y = -500
	MAX_Y = 500

	NUM_OBJS = 4000
)

func TestXZListAOIManager(t *testing.T) {
	testAOI(NewXZListAOICalculator(), NUM_OBJS)
}

type TestObj struct {
	aoi AOI
	Id  int
}

func (obj *TestObj) OnEnterAOI(other *AOI) {
	//log.Printf("%s: enter aoi %s", obj, obj.getObj(other))
	//print(len(obj.aoi.neighbors), " ")
}

func (obj *TestObj) OnLeaveAOI(other *AOI) {
	//log.Printf("%s: leave aoi %s", obj, obj.getObj(other))
	//print(len(obj.aoi.neighbors), " ")
}

func (obj *TestObj) String() string {
	return fmt.Sprintf("TestObj<%d>", obj.Id)
}

func (obj *TestObj) getObj(aoi *AOI) *TestObj {
	return (*TestObj)(unsafe.Pointer(aoi))
}

func randCoord(min, max Coord) Coord {
	return min + Coord(rand.Intn(int(max)-int(min)))
}

func testAOI(aoiman AOIManager, numAOI int) {
	objs := []*TestObj{}
	for i := 0; i < numAOI; i++ {
		obj := &TestObj{Id: i + 1}
		InitAOI(&obj.aoi, obj)
		objs = append(objs, obj)
		aoiman.Enter(&obj.aoi, randCoord(MIN_X, MAX_X), randCoord(MIN_Y, MAX_Y))
	}

	for i := 0; i < 10; i++ {
		t0 := time.Now()
		for _, obj := range objs {
			aoiman.Moved(&obj.aoi, obj.aoi.x+randCoord(-10, 10), obj.aoi.y+randCoord(-10, 10))
		}
		dt := time.Now().Sub(t0)
		log.Printf("%T tick %d objects takes %s", aoiman, numAOI, dt)
	}

}
