package editor

type Editor struct {
	Text [][]rune
	X    int
	Y    int
}

func NewEditor() *Editor {
	return &Editor{
		[][]rune{{}},
		0,
		0,
	}
}

func (ed *Editor) MoveLeft() {
	if ed.X == 0 {
		return
	}
	ed.X--
}

func (ed *Editor) MoveRight() {
	if len(ed.Text[ed.Y]) == ed.X {
		return
	}
	ed.X++
}

func (ed *Editor) MoveUp() {
	if ed.Y == 0 {
		return
	}

	ed.Y--
	if len(ed.Text[ed.Y])-1 < ed.X {
		ed.X = len(ed.Text[ed.Y])
	}
}

func (ed *Editor) MoveDown() {
	if len(ed.Text)-1 == ed.Y {
		return
	}

	ed.Y++
	if len(ed.Text[ed.Y])-1 < ed.X {
		ed.X = len(ed.Text[ed.Y])
	}
}

func (ed *Editor) GoToLineStart() {
	ed.X = 0
}

func (ed *Editor) GoToLineEnd() {
	ed.X = len(ed.Text[ed.Y])
}

func (ed *Editor) AddLine() {
	currentLine := ed.Text[ed.Y]

	remainingLine := make([]rune, len(currentLine[:ed.X]))
	copy(remainingLine, currentLine[:ed.X])
	newLine := make([]rune, len(currentLine[ed.X:]))
	copy(newLine, currentLine[ed.X:])

	ed.Text[ed.Y] = remainingLine

	newText := make([][]rune, len(ed.Text)+1)
	headLines := ed.Text[:ed.Y+1]
	tailLines := ed.Text[ed.Y+1:]
	copy(newText[:ed.Y+1], headLines)
	newText[ed.Y+1] = newLine
	copy(newText[ed.Y+2:], tailLines)
	ed.Text = newText

	ed.X = 0
	ed.Y++
}

func (ed *Editor) AddRune(r rune) {

	if len(ed.Text[ed.Y]) == ed.X {
		ed.Text[ed.Y] = append(ed.Text[ed.Y], r)
		ed.X++
		return
	}

	newLine := make([]rune, len(ed.Text[ed.Y])+1)
	head := ed.Text[ed.Y][:ed.X]
	tail := ed.Text[ed.Y][ed.X:]

	copy(newLine[:ed.X], head)
	newLine[ed.X] = r
	copy(newLine[ed.X+1:], tail)

	ed.Text[ed.Y] = newLine

	ed.MoveRight()
}

func (ed *Editor) RemoveBackwardRune() {
	if ed.X == 0 && ed.Y == 0 {
		return
	}

	if ed.X == 0 {
		currentLine := ed.Text[ed.Y]
		currentLineLen := len(currentLine)

		prevLine := ed.Text[ed.Y-1]
		prevLineLen := len(prevLine)

		newLine := make([]rune, prevLineLen+currentLineLen)
		copy(newLine[:prevLineLen], prevLine)
		copy(newLine[prevLineLen:], currentLine)
		ed.Text[ed.Y-1] = newLine

		ed.X = prevLineLen
		ed.removeLine(ed.Y)
		ed.Y--

		return
	}

	ed.X--

	newLine := make([]rune, len(ed.Text[ed.Y])-1)
	head := ed.Text[ed.Y][:ed.X+1]
	tail := ed.Text[ed.Y][ed.X+1:]

	copy(newLine[:ed.X], head)
	copy(newLine[ed.X:], tail)

	ed.Text[ed.Y] = newLine
}

func (ed *Editor) RemoveForwardRune() {
	if ed.X == len(ed.Text[ed.Y]) && ed.Y == len(ed.Text)-1 {
		return
	}

	if ed.X == len(ed.Text[ed.Y]) {
		currentLine := ed.Text[ed.Y]
		currentLineLen := len(currentLine)

		nextLine := ed.Text[ed.Y+1]
		nextLineLen := len(nextLine)

		newLine := make([]rune, nextLineLen+currentLineLen)
		copy(newLine[:currentLineLen], currentLine)
		copy(newLine[currentLineLen:], nextLine)
		ed.Text[ed.Y] = newLine

		ed.X = currentLineLen
		ed.removeLine(ed.Y + 1)

		return
	}

	newLine := make([]rune, len(ed.Text[ed.Y])-1)
	head := ed.Text[ed.Y][:ed.X]
	tail := ed.Text[ed.Y][ed.X+1:]

	copy(newLine[:ed.X], head)
	copy(newLine[ed.X:], tail)

	ed.Text[ed.Y] = newLine
}

func (ed *Editor) removeLine(lineOffset int) {

	newText := make([][]rune, len(ed.Text)-1)
	head := ed.Text[:lineOffset]
	tail := ed.Text[lineOffset+1:]

	copy(newText[:lineOffset], head)
	copy(newText[lineOffset:], tail)

	ed.Text = newText
}
