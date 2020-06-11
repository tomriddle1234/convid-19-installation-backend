package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"time"
)

func infoHandler(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{"info":"Convid-19 installation backend. v1.0"})
}

func accumulateSignalHandler(c *gin.Context){
	receivedCount++
	currentTime := time.Now()
	difference := uint8(math.Round(currentTime.Sub(lastReactTime).Seconds()))

	// if exceed response time limit, then resets all
	if difference > responseTime{
		accumulatedCount = 1
		lastReactTime = currentTime
		if chargeOrNot != false{
			chargeOrNot = false
			chargeChannel <- false
		}
	} else {
		accumulatedCount += 1
		// if exceed accumulate count within response time, switch status
		if accumulatedCount >= signalReactLowerLimit {
			lastReactTime = currentTime
			if chargeOrNot != true {
				chargeOrNot = true
				chargeChannel <- true
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"status":"Accumulated Signal Received.",
		"accumulatedCount": accumulatedCount,
		"chargeOrNot": chargeOrNot,
	})
}

func resetBallonHandler(c *gin.Context){
	//reset all numbers
	initData()
	c.JSON(http.StatusOK, gin.H{"status":"Reset Ballon Signal Received."})
}

func reliefSignalHandler(c *gin.Context){
	chargeOrNot = false
	reliefOrNot = true
	chargeChannel <- false
	reliefChannel <- true
	c.JSON(http.StatusOK, gin.H{"status":"Relief Signal Received."})
}

func statusCheckHandler(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{"charge":chargeOrNot, "relief":reliefOrNot})
}

func listenToChannel(){
	chargeMsg := chargeOrNot
	reliefMsg := reliefOrNot
	for {
		// deal with charge
		newChargeMsg := <- chargeChannel
		fmt.Println("listenToChannel() observer triggered")
		// only the first charge switch to true trigger it's turning to false with a timer
		if newChargeMsg != chargeMsg && newChargeMsg {
			chargeMsg = newChargeMsg
			chargeTimer := time.NewTimer(time.Duration(statusOnTime) * time.Second)
			fmt.Println("Start a charge recover timer.")
			go func() {
				<-chargeTimer.C
				fmt.Println("Now charge status recover to false.")
				chargeOrNot = false
				chargeMsg = false
			}()
		}

		// deal with relief
		newReliefMsg := <- reliefChannel
		fmt.Println(newReliefMsg)
		if newReliefMsg != reliefMsg && newReliefMsg {
			reliefMsg = newReliefMsg
			reliefTimer := time.NewTimer(time.Duration(statusOnTime) * time.Second)
			go func(){
				<- reliefTimer.C
				fmt.Println("Now relief status recover to false.")
				reliefOrNot = false
				reliefMsg = false
			}()
		}
	}

}