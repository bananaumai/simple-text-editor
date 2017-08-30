package editor

type EditorEvent int

const (
	evMoveUp = iota + 1
	evMoveDown
)

const (
	EditorEventMoveUp   EditorEvent = evMoveUp
	EditorEventMoveDown EditorEvent = evMoveDown
)

type EventListener func(ed *Editor)
