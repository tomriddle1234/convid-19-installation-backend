package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math"
	"net/http"
	"time"
)

func infoHandler(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{"info":"Convid-19 installation backend. v1.0"})
}

func accumulateSignalHandler(c *gin.Context){

	// watch out, if it's charging, then return directly not changing anything.

	if chargeStatus.val == true{
		c.JSON(http.StatusOK, gin.H{
			"status":"Accumulated Signal Received, Charging, so Not changing anything.",
			"accumulatedCount": accumulatedCount,
			"chargeOrNot": chargeStatus.val,
		})
	} else {
		receivedCount++
		currentTime := time.Now()
		difference := uint8(math.Round(currentTime.Sub(lastReactTime).Seconds()))

		// if exceed response time limit, then resets all
		if difference > responseTime{
			accumulatedCount = 1
			lastReactTime = currentTime

			chargeStatus.m.Lock()
			if chargeStatus.val != false{
				chargeChannel <- false
				chargeStatus.val = false
			}
			chargeStatus.m.Unlock()

		} else {
			accumulatedCount += 1
			// if exceed accumulate count within response time, switch status
			if accumulatedCount >= signalReactLowerLimit {
				lastReactTime = currentTime
				chargeStatus.m.Lock()
				if chargeStatus.val != true {
					chargeChannel <- true
					chargeStatus.val = true
				}
				chargeStatus.m.Unlock()
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"status":"Accumulated Signal Received.",
			"accumulatedCount": accumulatedCount,
			"chargeOrNot": chargeStatus.val,
		})
	}
}

func resetBallonHandler(c *gin.Context){
	//reset all numbers
	initData()
	c.JSON(http.StatusOK, gin.H{"status":"Reset Ballon Signal Received."})
}

func reliefSignalHandler(c *gin.Context){
	chargeStatus.m.Lock()
	chargeStatus.val = false
	chargeStatus.m.Unlock()

	reliefStatus.m.Lock()
	reliefStatus.val = true
	reliefStatus.m.Unlock()

	chargeChannel <- false
	reliefChannel <- true

	c.JSON(http.StatusOK, gin.H{"status":"Relief Signal Received."})
}

func statusCheckHandler(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{"charge":chargeStatus.val, "relief":reliefStatus.val})
}

func listenToChannel(){
	//chargeStatus.m.Lock()
	//reliefStatus.m.Lock()
	chargeMsg := chargeStatus.val
	reliefMsg := reliefStatus.val
	//chargeStatus.m.Unlock()
	//reliefStatus.m.Unlock()

	for {
		// deal with charge
		newChargeMsg := <- chargeChannel
		fmt.Println("listenToChannel() observer triggered")
		// only the first charge switch to true trigger it's turning to false with a timer
		if newChargeMsg != chargeMsg && newChargeMsg {
			chargeMsg = newChargeMsg
			chargeTimer := time.NewTimer(time.Duration(statusOnTime) * time.Second)
			fmt.Println("Start a charge recover timer.")
			log.Println("Start a charge recover timer.")
			go func() {
				<-chargeTimer.C
				fmt.Println("Now charge status recover to false.")
				log.Println("Now charge status recover to false.")
				chargeStatus.m.Lock()
				chargeStatus.val = false
				chargeStatus.m.Unlock()

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
				reliefStatus.m.Lock()
				reliefStatus.val = false
				reliefStatus.m.Unlock()

				reliefMsg = false
			}()
		}
	}

}