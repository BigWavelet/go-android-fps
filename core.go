/*
* File: core.go
* Author : bigwavelet
* Description: core file
* Created: 2016-08-26
 */

package main

import (
    "os"
    "strings"
)



func strip(str string) (result string) {
    result = strings.Replace(str, " ", "", -1)
    result = strings.Replace(str, "\t", "", -1)
    return 
}


func splitLines(str string) (result []string) {
    tmp := strings.Replace(str, "\r\n", "\n", -1)
    tmp = strings.Replace(str, "\r", "", -1)
    result = strings.Split(tmp, "\n")
    return
}

func PathExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return false, err
}