package main

import (
	"github.com/nsf/termbox-go"
)

type Editor struct {
	text [][]rune
	x int
	y int
}

func NewEditor() *Editor {
	return &Editor{
		[][]rune{{}},
		0,
		0,
	}
}

func (ed *Editor) AddLine() {
	ed.text = append(ed.text, []rune{})
	ed.x = 0
	ed.y++
}

func (ed *Editor) MoveLeft() {
	if ed.x == 0 {
		return
	}
	ed.x--
}

func (ed *Editor) MoveRight() {
	if len(ed.text[ed.y])-1 == ed.x {
		return
	}
	ed.x++
}

func (ed *Editor) AddRune(r rune) {
	ed.text[ed.y] = append(ed.text[ed.y], r)
	ed.x++
}

func (ed *Editor) Draw() {
	const color = termbox.ColorDefault
	for i, l := range ed.text {
		for j, r := range l {
			termbox.SetCell(j, i, r, color, color)
		}
	}
	termbox.SetCursor(ed.x, ed.y)
	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	const color = termbox.ColorDefault
	termbox.Clear(color, color)

	ed := NewEditor()

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
		case termbox.KeyArrowLeft:
			ed.MoveLeft()
		case termbox.KeyArrowRight:
			ed.MoveRight()
		default:
			ed.AddRune(ev.Ch)
		}
		ed.Draw()
	}
}
