package main

import (
	"testing"
	"time"
)

func getCurrentTime() (time.Time , string , time.Time) {
    currentTime := time.Now().Format("02-01-2006 15:04:05")
    parsedTime, err := time.Parse("02-01-2006 15:04:05", currentTime)
    currentTest := time.Now()
    if err != nil {
        panic(err)
    }

    return parsedTime , currentTime , currentTest
}

func getLocalTimeIstanbul() (time.Time ){
    location, _ := time.LoadLocation("Europe/Istanbul")
    currentTime := time.Now().In(location).Format("02-01-2006 15:04:05")
    parsedTime, err := time.Parse("02-01-2006 15:04:05", currentTime)
    if err != nil {
        panic(err)
    }

    return parsedTime

}


func TestGetCurrentTime(t *testing.T) {
    currentTime, currentTimeString , currentTest := getCurrentTime()
    localIstanbul := getLocalTimeIstanbul() 

    if _, ok := interface{}(currentTime).(time.Time); ok {
        t.Errorf("Beklenen tür date : %T",currentTime )
        t.Log("Current Time : ", currentTime)

    }


    t.Log("Time : ", currentTime)
    t.Log("Time String : ", currentTimeString)
    t.Log("Current Test : ", currentTest)
    t.Log("Local Time : ", localIstanbul)

}

