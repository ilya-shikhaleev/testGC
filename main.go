package main

import (
	"log"
	"net/http"
	"runtime"
	"time"
)

import _ "net/http/pprof"

var People []*Person

const PeopleCount = 10000000

type Person struct {
	Name string
	Age  uint
}

func allocateInitial() {
	for i := 0; i < PeopleCount; i++ {
		People = append(People, &Person{"Ilya", 24})
	}
}

func allocateMore() {
	m := runtime.MemStats{}
	var prevTotalPause, currTotalPause uint64
	for {
		log.Println("Start allocating")
		for i := 0; i < PeopleCount; i++ {
			People = append(People, &Person{"Shikhaleev", 25})
		}
		People = People[0:PeopleCount]
		runtime.GC()
		runtime.ReadMemStats(&m)
		currTotalPause = m.PauseTotalNs / 1000000 //Nanoseconds to milliseconds
		log.Println("GC pause =", currTotalPause-prevTotalPause, "ms")
		prevTotalPause = currTotalPause
		time.Sleep(3 * time.Second)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	allocateInitial()
	go allocateMore()
	log.Println(http.ListenAndServe("127.0.0.1:8080", nil))
}
