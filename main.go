package main

import (
	"os"
	"runtime"

	"git.mrcyjanek.net/mrcyjanek/bytestreamer/front"
	"git.mrcyjanek.net/mrcyjanek/bytestreamer/gui"
)

func main() {
	if runtime.GOOS == "android" {
		os.Chdir("/data/data/x.x.bytestreamer/files")
	}
	go gui.Start()
	front.Start()
	<-make(chan int)
}
