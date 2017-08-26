package main

import (
	"github.com/bananaumai/simple-text-editor/editor"
	"github.com/bananaumai/simple-text-editor/screen"
	"github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	const color = termbox.ColorDefault
	termbox.Clear(color, color)

	ed := editor.NewEditor()
	sc := screen.NewScreen()

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
			ed.AddLine()
		case termbox.KeyArrowLeft, termbox.KeyCtrlB:
			ed.MoveLeft()
		case termbox.KeyArrowRight, termbox.KeyCtrlF:
			ed.MoveRight()
		case termbox.KeyArrowUp, termbox.KeyCtrlP:
			ed.MoveUp()
		case termbox.KeyArrowDown, termbox.KeyCtrlN:
			ed.MoveDown()
		case termbox.KeyCtrlA:
			ed.GoToLineStart()
		case termbox.KeyCtrlE:
			ed.GoToLineEnd()
		case termbox.KeyBackspace, termbox.KeyBackspace2:
			ed.RemoveBackwardRune()
		case termbox.KeyDelete, termbox.KeyCtrlD:
			ed.RemoveForwardRune()
		case termbox.KeySpace:
			ed.AddRune(' ')
		case termbox.KeyTab:
			ed.AddRune('\t')
		default:
			if ev.Ch != 0 {
				ed.AddRune(ev.Ch)
			}
		}

		sc.Draw(ed)
	}
}
