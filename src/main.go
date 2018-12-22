package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"ot"
)

func init() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	termbox.SetCursor(0, 0)
	termbox.HideCursor()
}

func main() {

	ot.Run(ot.QryServe(ot.QryApp()))

	pause()
}


func pause() {
	fmt.Println("\n\n\n按任意键退出...")
Loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			break Loop
		}
	}
}