/*
* File: device.go
* Author : bigwavelet
* Description: android device interface
* Created: 2016-08-26
 */

package main

import (
    //"log"
    "strconv"
    "os/exec"
    "regexp"
    "strings"
)


const (
    ADB_PATH = "adb"
)




type AdbDevice struct {
    Serial string
}

//adb shell
func (d *AdbDevice) shell(cmds ...string) (out string, err error) {
    args := []string{}
    if len(d.Serial) == 0 {
        args = append(args, "shell")    
    }else{
        args = append(args, "-s", d.Serial, "shell")
    }
    args = append(args, cmds...)
    output, err := exec.Command(ADB_PATH, args...).Output()
    out = string(output)
    return 
}

//adb command

func (d *AdbDevice) run(cmds ...string) (out string, err error) {
    args := []string{}
    if len(d.Serial) != 0 {
        args = append(args, "-s", d.Serial)
    }
    args = append(args, cmds...)
    output, err := exec.Command(ADB_PATH, args...).Output()
    out = string(output)
    return
}

// get device list
func (d *AdbDevice) getDeviceList() (devs []string, err error) {
    cmd := exec.Command(ADB_PATH, "devices")
    out, err := cmd.Output()
    if err != nil {
        return nil, err
    }
    patten := regexp.MustCompile(`^[^\s]+\t(device)`)
    lines := strings.Split(string(out), "\n")
    
    for i := 0; i < len(lines); i++ {
        line := lines[i]
        if patten.MatchString(line) {
            fields := strings.Split(line, "\t")
            if len(fields) != 2 {
                continue
            }
            sn := fields[0]
            devs = append(devs, sn)
        }
    }
    return
}



//get device info
// func getDeviceInfo(serialno string) {
// }




func (d* AdbDevice) getTopView() (topView string, err error) {
    out, err := d.shell("dumpsys SurfaceFlinger")
    if err != nil {
        return "", err
    }
    output := strings.Replace(out, "\r\n", "\n", -1)
    lines := strings.Split(output, "\n")
    maxArea := 0
    topView = ""
    for index, line := range lines {
        if !strings.HasPrefix(line, "+ Layer "){
            continue
        } 
        patten := regexp.MustCompile(`\((.+)\)`)
        if !patten.MatchString(line){
            continue
        }
        viewName := patten.FindString(line)
        viewName = strings.Replace(viewName, "(", "", -1)
        viewName = strings.Replace(viewName, ")", "", -1)
        patten2 := regexp.MustCompile(`\d+`)
        tmp := patten2.FindAllString(lines[index+4], -1)
        if len(tmp) != 4 {
            continue
        }
        x0, _ := strconv.Atoi(tmp[0])
        y0, _ := strconv.Atoi(tmp[1])
        x1, _ := strconv.Atoi(tmp[2])
        y1, _ := strconv.Atoi(tmp[3])
        curArea := (x1 - x0) * (y1 - y0)
        if curArea > maxArea {
            maxArea = curArea
            topView = viewName 
        }
    }
    return topView, nil
}




/*
func main() {
    if err := exec.Command("adb", "start-server").Run(); err != nil {
        log.Fatal(err)
    }
    
    d := AdbDevice{"S8GMTO7TT44P599S"}
    result, err := d.getTopView()
    log.Println(result)
    log.Println(err)
}

*/