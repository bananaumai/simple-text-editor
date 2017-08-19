package main

import (
	"github.com/nsf/termbox-go"
)

type Editor struct {
	text [][]rune
	x    int
	y    int
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
	if len(ed.text[ed.y])-1 < ed.x {
		ed.x = len(ed.text[ed.y])
	}
}

func (ed *Editor) MoveDown() {
	if len(ed.text)-1 == ed.y {
		return
	}

	ed.y++
	if len(ed.text[ed.y])-1 < ed.x {
		ed.x = len(ed.text[ed.y])
	}
}

func (ed *Editor) AddLine() {
	currentLine := ed.text[ed.y]
	remainingLine := currentLine[:ed.x]
	newLine := currentLine[ed.x:]

	ed.text[ed.y] = remainingLine

	newText := make([][]rune, len(ed.text)+1)
	headLines := ed.text[:ed.y+1]
	tailLines := ed.text[ed.y+1:]
	copy(newText[:ed.y+1], headLines)
	newText[ed.y+1] = newLine
	copy(newText[ed.y+2:], tailLines)
	ed.text = newText

	ed.x = 0
	ed.y++
}

func (ed *Editor) RemoveLine() {
	if ed.y == 0 {
		return
	}

	newText := make([][]rune, len(ed.text)-1)
	head := ed.text[:ed.y]
	tail := ed.text[ed.y+1:]

	copy(newText[:ed.y], head)
	copy(newText[ed.y:], tail)

	ed.text = newText

	ed.y--
}

func (ed *Editor) AddRune(r rune) {
	if len(ed.text[ed.y]) == ed.x {
		ed.text[ed.y] = append(ed.text[ed.y], r)
		ed.x++
		return
	}

	newLine := make([]rune, len(ed.text[ed.y])+1)
	head := ed.text[ed.y][:ed.x]
	tail := ed.text[ed.y][ed.x:]

	copy(newLine[:ed.x], head)
	newLine[ed.x] = r
	copy(newLine[ed.x+1:], tail)

	ed.text[ed.y] = newLine

	ed.MoveRight()
}

func (ed *Editor) RemoveBackwardRune() {
	if ed.x == 0 && ed.y == 0 {
		return
	}

	if ed.x == 0 {
		currentLine := ed.text[ed.y]
		currentLineLen := len(currentLine)

		prevLine := ed.text[ed.y-1]
		prevLineLen := len(prevLine)

		newLine := make([]rune, prevLineLen+currentLineLen)
		copy(newLine[:prevLineLen], prevLine)
		copy(newLine[prevLineLen:], currentLine)
		ed.text[ed.y-1] = newLine

		ed.x = prevLineLen
		ed.RemoveLine()

		return
	}

	ed.MoveLeft()

	newLine := make([]rune, len(ed.text[ed.y])-1)
	head := ed.text[ed.y][:ed.x+1]
	tail := ed.text[ed.y][ed.x+1:]

	copy(newLine[:ed.x], head)
	copy(newLine[ed.x:], tail)

	ed.text[ed.y] = newLine
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
		case termbox.KeyBackspace, termbox.KeyBackspace2:
			ed.RemoveBackwardRune()
		default:
			if ev.Ch != 0 {
				ed.AddRune(ev.Ch)
			}
		}
		ed.Draw()
	}
}
