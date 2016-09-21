/*
* File: web.go
* Author : bigwavelet
* Description: web handler
* Created: 2016-08-26
 */

package main

import (
    "os"
    "fmt"
    "log"
    "time"
    "strings"
    "io/ioutil"
    "net/http"
    "path/filepath"
    "github.com/gorilla/mux"
    "encoding/json"
    "github.com/gorilla/websocket"
)


type WebConfig struct {
    User    string
    Version string
}


type FpsConfig struct {
    closed bool
    fpsDataBuf []int
}

var (
    defaultDataDir = "data"
)
func (fpscfg *FpsConfig) renderHTML(w http.ResponseWriter, name string, data interface{}) {
    w.Header().Set("Content-Type", "text/html")
    wc := WebConfig{}
    wc.Version = VERSION
    if data == nil {
        data = wc
    }
    executeTemplate(w, name, data)
}


func (fpscfg *FpsConfig) hIndex(w http.ResponseWriter, r *http.Request) {
    fpscfg.renderHTML(w, "index", nil)
}



func (fpscfg *FpsConfig) startFps(w http.ResponseWriter, r *http.Request) {
    log.Println("start fps...")
    var data []byte
    data, _ = json.Marshal(map[string]interface{}{
        "status": 1,
        "data": "successfully started.",
        }) 
    fpscfg.closed = false
    fpscfg.fpsDataBuf = fpscfg.fpsDataBuf[:0]
    w.Write(data)
}


func (fpscfg *FpsConfig) stopFps(w http.ResponseWriter, r *http.Request) {
    var data []byte
    filename := mux.Vars(r)["filename"]
    if !strings.HasSuffix(filename, ".txt") {
        filename = filename + ".txt"
    }
    isDirExisted, _ := PathExists(defaultDataDir)
    if ! isDirExisted {
        err := os.MkdirAll(defaultDataDir, 0777)
        if err != nil {
            data, _ = json.Marshal(map[string]interface{}{
            "status": 0,
            "data": "can not create dir.",
            }) 
            log.Printf(err.Error())
            w.Write(data)
            return
        }
    }
    filename = filepath.Join(defaultDataDir, filename)
    log.Println("stop fps....", filename)
    fpscfg.closed = true
    fout, err := os.Create(filename)
    defer fout.Close()
    if err != nil {
        data, _ = json.Marshal(map[string]interface{}{
        "status": 0,
        "data": "filename exists.",
        }) 
        log.Printf(err.Error())
        w.Write(data)
        return    
    }
    for _, val := range fpscfg.fpsDataBuf {
        log.Printf("%v ", val)
        fout.WriteString(fmt.Sprintf("%d\r\n", val))
    }
    data, _ = json.Marshal(map[string]interface{}{
        "status": 1,
        "data": "successfully stopped.",
        }) 
    w.Write(data)

    
}

func (fpscfg *FpsConfig) fileList(w http.ResponseWriter, r *http.Request) {
    var files []string
    var data []byte
    file_list, err := ioutil.ReadDir(defaultDataDir)
    if err != nil {
        fmt.Println(err.Error())
        data, _ = json.Marshal(map[string]interface{}{
        "status": 0,
        "data": "can not found any file",
        }) 
    w.Write(data)
        return
    }
    for _, val := range file_list {
        log.Println(val.Name())
        fName := strings.Replace(val.Name(), ".txt", "", -1)
        files = append(files, fName)
    }
    data, _ = json.Marshal(map[string]interface{}{
        "status": 1,
        "data": files,
        }) 
    w.Write(data)

}


func (fpscfg *FpsConfig) getFps(w http.ResponseWriter, r *http.Request) {
    var data []byte
    filename := mux.Vars(r)["filename"]
    if !strings.HasSuffix(filename, ".txt") {
        filename = filename + ".txt"
    }
    filename = filepath.Join(defaultDataDir, filename)
    fin, err := os.Open(filename)
    if err != nil {
        data, _ = json.Marshal(map[string]interface{}{
            "status": 0,
            "data": "can not find file.",
        }) 
        w.Write(data)
        return
    }
    defer fin.Close()
    fd, err := ioutil.ReadAll(fin)
    if err != nil {
        data, _ = json.Marshal(map[string]interface{}{
            "status": 0,
            "data": "can not read file.",
        }) 
        w.Write(data)
        return
    }
    fpsData := splitLines(string(fd))
    for idx, val := range fpsData {
        if val == "" || val == "\n" {
            fpsData = append(fpsData[:idx], fpsData[idx+1:]...)
        }
    }
    data, _ = json.Marshal(map[string]interface{}{
        "status": 1,
        "data": fpsData,
        }) 
    w.Write(data)


}

var upgrader = websocket.Upgrader{}


//performance
func (fpscfg *FpsConfig) wsPerf(w http.ResponseWriter, r *http.Request) {
    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Print("upgrade:", err)
        return
    }
    defer c.Close()

    name := mux.Vars(r)["name"]
    switch name {
        case "fps":
            for {
                if !fpscfg.closed {
                    fps, _ := <-FpsData
                    data := make(map[string]int)
                    data["fps"] = fps
                    log.Println(data)
                    fpscfg.fpsDataBuf = append(fpscfg.fpsDataBuf, fps)
                    err := c.WriteJSON(data)
                    if err != nil {
                        break
                    }
                }
                time.Sleep(time.Millisecond * 10)
            }
        case "cpu":
            log.Println("coming soon...")
        default:
            log.Println("wrong performance type")
    }
}


func startWebServer(port string) {
    var fpsDataBuf []int
    fpscfg := FpsConfig{true, fpsDataBuf}
    r := mux.NewRouter()
    r.HandleFunc("/", fpscfg.hIndex)
    r.HandleFunc("/start_fps", fpscfg.startFps)
    r.HandleFunc("/stop_fps/{filename}", fpscfg.stopFps)
    r.HandleFunc("/api/fps_data/{filename}", fpscfg.getFps)
    r.HandleFunc("/api/file_list", fpscfg.fileList)
    r.HandleFunc("/ws/perfs/{name}", fpscfg.wsPerf)
    http.Handle("/", r)
    log.Println("start webserver here...")
    http.ListenAndServe(":" + port, nil)
}



