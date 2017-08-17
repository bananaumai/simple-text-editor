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

func (ed *Editor) MoveLeft() {
	if ed.x == 0 {
		return
	}
	ed.x--
}

func (ed *Editor) MoveRight() {
	if len(ed.text[ed.y]) == ed.x {
		return
	}
	ed.x++
}

func (ed *Editor) MoveUp() {
	if ed.y == 0 {
		return
	}

	ed.y--
	if len(ed.text[ed.y]) - 1 < ed.x {
		ed.x = len(ed.text[ed.y])
	}
}

func (ed *Editor) MoveDown() {
	if len(ed.text) - 1 == ed.y {
		return
	}

	ed.y++
	if len(ed.text[ed.y]) - 1 < ed.x {
		ed.x = len(ed.text[ed.y])
	}
}

func (ed *Editor) AddLine() {
	if len(ed.text) - 1 == ed.y {
		ed.text = append(ed.text, []rune{})
		ed.x = 0
		ed.y++
		return
	}

	newText := make([][]rune, len(ed.text) + 1)
	head := ed.text[:(ed.y + 1)]
	tail := ed.text[(ed.y + 1):]

	copy(newText[:(ed.y + 1)], head)
	newText[(ed.y + 1)] = []rune{}
	copy(newText[(ed.y + 2):], tail)

	ed.text = newText

	ed.MoveDown()
}

func (ed *Editor) AddRune(r rune) {
	if len(ed.text[ed.y]) == ed.x {
		ed.text[ed.y] = append(ed.text[ed.y], r)
		ed.x++
		return
	}

	newLine := make([]rune, len(ed.text[ed.y]) + 1)
	head := ed.text[ed.y][:ed.x]
	tail := ed.text[ed.y][ed.x:]

	copy(newLine[:ed.x], head)
	newLine[ed.x] = r
	copy(newLine[(ed.x + 1):], tail)

	ed.text[ed.y] = newLine

	ed.MoveRight()
}

func (ed *Editor) Draw() {
	const color = termbox.ColorDefault

	termbox.Clear(color, color)

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
		case termbox.KeyArrowUp:
			ed.MoveUp()
		case termbox.KeyArrowDown:
			ed.MoveDown()
		default:
			ed.AddRune(ev.Ch)
		}
		ed.Draw()
	}
}
