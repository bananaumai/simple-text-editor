package main

import (
	"github.com/mattn/go-runewidth"
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

func (ed *Editor) GoToLineStart() {
	ed.x = 0
}

func (ed *Editor) GoToLineEnd() {
	ed.x = len(ed.text[ed.y])
}

func (ed *Editor) AddLine() {
	currentLine := ed.text[ed.y]

	remainingLine := make([]rune, len(currentLine[:ed.x]))
	copy(remainingLine, currentLine[:ed.x])
	newLine := make([]rune, len(currentLine[ed.x:]))
	copy(newLine, currentLine[ed.x:])

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
		ed.removeLine(ed.y)
		ed.y--

		return
	}

	ed.x--

	newLine := make([]rune, len(ed.text[ed.y])-1)
	head := ed.text[ed.y][:ed.x+1]
	tail := ed.text[ed.y][ed.x+1:]

	copy(newLine[:ed.x], head)
	copy(newLine[ed.x:], tail)

	ed.text[ed.y] = newLine
}

func (ed *Editor) RemoveForwardRune() {
	if ed.x == len(ed.text[ed.y]) && ed.y == len(ed.text)-1 {
		return
	}

	if ed.x == len(ed.text[ed.y]) {
		currentLine := ed.text[ed.y]
		currentLineLen := len(currentLine)

		nextLine := ed.text[ed.y+1]
		nextLineLen := len(nextLine)

		newLine := make([]rune, nextLineLen+currentLineLen)
		copy(newLine[:currentLineLen], currentLine)
		copy(newLine[currentLineLen:], nextLine)
		ed.text[ed.y] = newLine

		ed.x = currentLineLen
		ed.removeLine(ed.y + 1)

		return
	}

	newLine := make([]rune, len(ed.text[ed.y])-1)
	head := ed.text[ed.y][:ed.x]
	tail := ed.text[ed.y][ed.x+1:]

	copy(newLine[:ed.x], head)
	copy(newLine[ed.x:], tail)

	ed.text[ed.y] = newLine
}

func (ed *Editor) Draw() {
	const color = termbox.ColorDefault

	termbox.Clear(color, color)

	cursorPosX := 0
	for y, line := range ed.text {
		posX := 0
		for x, r := range line {
			termbox.SetCell(posX, y, r, color, color)
			width := runewidth.RuneWidth(r)
			posX += width
			if ed.y == y && ed.x > x {
				cursorPosX += width
			}
		}
	}

	termbox.SetCursor(cursorPosX, ed.y)

	termbox.Flush()
}

func (ed *Editor) removeLine(lineOffset int) {

	newText := make([][]rune, len(ed.text)-1)
	head := ed.text[:lineOffset]
	tail := ed.text[lineOffset+1:]

	copy(newText[:lineOffset], head)
	copy(newText[lineOffset:], tail)

	ed.text = newText
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
		ed.Draw()
	}
}
