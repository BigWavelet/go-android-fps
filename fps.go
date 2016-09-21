/*
* File: fps.go
* Author : bigwavelet
* Description: get android fps
* Created: 2016-08-26
 */

package main

import (
	"log"
    "errors"
    "time"
    //"regexp"
    "strings"
    "strconv"
)


const (
	nanosecondsPerSecond = 1e9
)

var (
    FpsData = make(chan int, 1)
)


// init frame data
func initFrameData(view string, serial string) (refreshPeriod float64, baseTimestamp float64, timestamps []float64, err error){
    d := AdbDevice{serial}
    out, _ := d.shell("dumpsys SurfaceFlinger --latency-clear " + view)
    output := strip(out)
    if output != "" {
        err = errors.New("not supported.")
        log.Println("not supported.")
        return    
    }
    
    time.Sleep(1e8)
    refreshPeriod, timestamps = frameData(view, serial)
    baseTimestamp = 0
    baseIndex := 0
    for _, timestamp := range timestamps {
        if timestamp != 0 {
            baseTimestamp = timestamp
            break
        }
        baseIndex++ 
    }
    if baseTimestamp == 0 {
        err = errors.New("Initial frame collect failed.")
        return
    }
    timestamps = timestamps[baseIndex:]
    return    
}


func frameData(view string, serial string) (refreshPeriod float64, timestamps []float64) {
    d := AdbDevice{serial}
    out, _ := d.shell("dumpsys SurfaceFlinger --latency " + view)
    lines := splitLines(out)
    line_0, _ := strconv.Atoi(lines[0])
    refreshPeriod = float64(line_0) / float64(nanosecondsPerSecond)
    //log.Println(line_0)
    //log.Println("refreshPeriod:", refreshPeriod)
    for _, line := range lines[1:] {
        fields := strings.Fields(line)
        if len(fields) != 3 {
            continue
        }
        //start, _ := strconv.Atoi(fields[0])
        submitting, _ := strconv.Atoi(fields[1])
        //submitted, _ := strconv.Atoi(fields[2])
        if submitting == 0 {
            continue
        }
        timestamp := float64(submitting) / float64(nanosecondsPerSecond)
        //log.Println("timestamp:", timestamp)
        timestamps = append(timestamps, timestamp)
    }
    return refreshPeriod, timestamps
}


func continueCollectFrameData(serial string) (err error){
    d := AdbDevice{serial}
    view, err := d.getTopView()
    log.Println("=========top view:" + view)
    if err != nil {
        err = errors.New("Fail to get current SurfaceFlinger view.")
        log.Println(("Fail to get current SurfaceFlinger view."))
        return 
    }
    _, baseTimestamp, timestamps, err := initFrameData(view, serial)
    if err != nil {
        return
    }
    for {
        _, tss := frameData(view, serial)
        lastIndex := 0
        length := len(timestamps)
        if length > 2 {
            recentTimestamp := timestamps[length - 2]
            for idx, val := range tss {
                if val == recentTimestamp {
                    lastIndex = idx
                    break
                }
            }
            timestamps = timestamps[:length-1]
            timestamps = append(timestamps, tss[lastIndex:]...)   
        }
        time.Sleep(time.Millisecond * 500)
        var ajustedTimestamps []float64
        for _, seconds := range(timestamps){
            seconds -= baseTimestamp
            if seconds > 1e6 {
                continue
            }
            ajustedTimestamps = append(ajustedTimestamps, seconds)
        }
        
        length = len(ajustedTimestamps)
        fromTime := ajustedTimestamps[length - 1] - 1.0
        fpsCount := 0 
        for _, seconds := range(ajustedTimestamps){
            if seconds > fromTime {
                fpsCount ++
            }
        }
        //log.Println(fpsCount)
        //write fps data to channel FpsData
        FpsData <- fpsCount
        //tt, ok := <- FpsData
        //log.Println(ok, ".......", tt)
    }
    return nil
}

