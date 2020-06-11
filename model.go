package main

import "time"

var ballonSize uint8
var ballonSizeUplimit uint8
var chargeOrNot bool
var reliefOrNot bool
var chargeChannel = make(chan bool,1)
var reliefChannel = make(chan bool,1)
var accumulatedCount uint8
var receivedCount uint64
var responseTime uint8 // seconds that accumulate account
var lastReactTime time.Time
var signalReactLowerLimit uint8
var statusOnTime uint8

func initData(){
	ballonSize = 0
	ballonSizeUplimit = 255
	receivedCount = 0
	chargeOrNot = false
	reliefOrNot = false
	// in 10 seconds, we count the received signal number, then charge the ballon
	responseTime = 10
	lastReactTime = time.Now()
	accumulatedCount = 0
	signalReactLowerLimit = 2
	statusOnTime = 10
}


