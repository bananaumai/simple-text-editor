package main

import (
	"github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)

	var x, y int

	mainloop:
	for {
		ev := termbox.PollEvent()

		if ev.Type != termbox.EventKey {
			continue
		}


		switch ev.Key {
		case termbox.KeyEsc:
			break mainloop
		case termbox.KeyEnter:
			x = 0
			y++
		default:
			//log.Println(ev.Ch)
			termbox.SetCell(x, y, ev.Ch, coldef, coldef)
			termbox.Flush()
			x++
		}
	}
}
