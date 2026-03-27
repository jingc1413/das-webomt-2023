package main

import (
	"fmt"

	"gomt/dasutils/tgui"
)

func main() {
	guiSrv, err := tgui.NewGuiServer()
	if err != nil {
		fmt.Println(err)
		return
	}
	guiSrv.Run()
}
