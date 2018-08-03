package aoi

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
	"time"
)

const (
	MIN_X = -500
	MAX_X = 500
	MIN_Y = -500
	MAX_Y = 500

	NUM_OBJS              = 1000
	VERIFY_NEIGHBOR_COUNT = true
)

func TestXZListAOIManager(t *testing.T) {
	testAOI(t, "XZListAOI", NewXZListAOIManager(100), NUM_OBJS)
}

func TestTowerAOIManager(t *testing.T) {
	testAOI(t, "TowerAOI", NewTowerAOIManager(MIN_X, MAX_X, MIN_Y, MAX_Y, 10), NUM_OBJS)
}

type TestObj struct {
	aoi            AOI
	Id             int
	neighbors      map[*TestObj]struct{}
	totalNeighbors int64
	nCalc          int64
}

func (obj *TestObj) OnEnterAOI(otheraoi *AOI) {
	if VERIFY_NEIGHBOR_COUNT {
		other := obj.getObj(otheraoi)
		if obj == other {
			panic("should not enter self")
		}
		if _, ok := obj.neighbors[other]; ok {
			log.Panicf("duplicae enter aoi")
		}
		obj.neighbors[other] = struct{}{}
		obj.totalNeighbors += int64(len(obj.neighbors))
		obj.nCalc += 1
	}
}

func (obj *TestObj) OnLeaveAOI(otheraoi *AOI) {
	if VERIFY_NEIGHBOR_COUNT {
		other := obj.getObj(otheraoi)
		if obj == other {
			panic("should not leave self")
		}
		if _, ok := obj.neighbors[other]; !ok {
			log.Panicf("duplicate leave aoi")
		}
		delete(obj.neighbors, other)
		obj.totalNeighbors += int64(len(obj.neighbors))
		obj.nCalc += 1
	}
}

func (obj *TestObj) String() string {
	return fmt.Sprintf("TestObj<%d>", obj.Id)
}

func (obj *TestObj) getObj(aoi *AOI) *TestObj {
	return aoi.Data.(*TestObj)
}

func randCoord(min, max Coord) Coord {
	return min + Coord(rand.Intn(int(max)-int(min)))
}

func testAOI(t *testing.T, manname string, aoiman AOIManager, numAOI int) {
	objs := []*TestObj{}
	for i := 0; i < numAOI; i++ {
		obj := &TestObj{Id: i + 1, neighbors: map[*TestObj]struct{}{}}
		InitAOI(&obj.aoi, 100, obj, obj)
		objs = append(objs, obj)
		aoiman.Enter(&obj.aoi, randCoord(MIN_X, MAX_X), randCoord(MIN_Y, MAX_Y))
	}

	//proffd, _ := os.OpenFile(manname+".pprof", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	//defer proffd.Close()

	//pprof.StartCPUProfile(proffd)
	for i := 0; i < 10; i++ {
		t0 := time.Now()
		for _, obj := range objs {
			aoiman.Moved(&obj.aoi, obj.aoi.x+randCoord(-10, 10), obj.aoi.y+randCoord(-10, 10))
			aoiman.Leave(&obj.aoi)
			aoiman.Enter(&obj.aoi, obj.aoi.x+randCoord(-10, 10), obj.aoi.y+randCoord(-10, 10))
		}
		dt := time.Now().Sub(t0)
		t.Logf("%s tick %d objects takes %s", manname, numAOI, dt)
	}

	for _, obj := range objs {
		aoiman.Leave(&obj.aoi)
	}

	//pprof.StopCPUProfile()

	if VERIFY_NEIGHBOR_COUNT {
		totalCalc := int64(0)
		for _, obj := range objs {
			totalCalc += obj.nCalc
		}
		println("Average calculate count:", totalCalc/int64(len(objs)))
	}
}
