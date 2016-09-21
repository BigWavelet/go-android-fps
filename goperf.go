/*
* File: fps.go
* Author : bigwavelet
* Description: get android fps
* Created: 2016-08-26
 */

package main

import (
    "io"
    "os"
    "fmt"
    "time"
    "runtime"
    //"regexp"
    "crypto/md5"
    "github.com/alecthomas/kingpin"
)


var (
    port = kingpin.Flag("port", "listen port, default:2333").Short('p').String()
    serial = kingpin.Flag("serial", "device serialno.").Short('s').String()
    version = kingpin.Flag("version", "Show version").Short('v').Bool()
)


var (
    VERSION = "0.0.1"
    BUILD_DATE = time.Now()
)



func showVersion() error {
    fmt.Printf("version: %v\n", VERSION)
    fmt.Printf("build: %v\n", BUILD_DATE)
    fmt.Printf("golang: %v\n", runtime.Version())
    fd, err := os.Open(os.Args[0])
    if err != nil {
        return err
    }
    md5h := md5.New()
    io.Copy(md5h, fd)
    fmt.Printf("md5sum: %x\n", md5h.Sum([]byte("")))
    return nil
}



func checkDeviceList() (isNoDevice bool, isMultiDev bool, err error) {
    d := AdbDevice{""}
    devs, err := d.getDeviceList()
    if err != nil {
        return 
    }
    isMultiDev = len(devs) > 1
    isNoDevice = len(devs) == 0
    return
}

func main() {

    kingpin.CommandLine.HelpFlag.Short('h')
    kingpin.Parse()

    if *version {
        showVersion()
        return
    }

    isNoDevice, isMultiDev, err := checkDeviceList()
    if err != nil {
        fmt.Println("Can not get device list, please check adb server and link your device.")
        return
    }
    if isNoDevice {
        fmt.Println("No device attached!")
        return
    }
    if len(*serial) == 0 && isMultiDev {
        fmt.Println("More than one device found! please use -s to specify the device")
        return
    }

    // start webserver
    if *port == "" {
        *port = "2333"
    }
    go startWebServer(*port)


    //get fps data
    continueCollectFrameData(*serial)
}
