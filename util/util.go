package util

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "path/filepath"
)

type GoMonitor struct {
    RootDir    []string         // the dir to monitor that set at configuration file
    sourceDir  []string         // directory files in RootDir
    FileStatus map[string]int64 // file path and modtime
    WorkDir    string
    change     chan bool
    Interval   int
    cmd        *exec.Cmd
    BuildCmd   string
    RunCmd     string
}

var DefMonitor = NewGoMonitor()

func NewGoMonitor() (goMonitor *GoMonitor) {
    m := make(map[string]int64, 20)
    ch := make(chan bool, 1)
    goMonitor = &GoMonitor{FileStatus: m, change: ch}

    return
}

// add config dir to RootDir
func (w *GoMonitor) AddRootDir(path string) (err error) {
    dirinfo, err := os.Stat(path)
    if err != nil {
        log.Printf("%s\n", err)
        return
    }
    if dirinfo.IsDir() {
        w.RootDir = append(w.RootDir, path)
    }

    // walk the file tree at path
    filepath.Walk(path, w.walkFn)
    return
}

func (w *GoMonitor) walkFn(path string, f os.FileInfo, err error) error {
    if err != nil {
        log.Printf("%s\n", err)
        return err
    }
    if f.IsDir() {
        w.sourceDir = append(w.sourceDir, path)
    }
    // add filepath and modtime to map
    w.FileStatus[path] = f.ModTime().Unix()
    return nil
}

func (w *GoMonitor) PrintFile() {
    for file, time := range w.FileStatus {
        fmt.Printf("%s   %v \n", file, time)
    }
}
