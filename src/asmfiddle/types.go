package asmfiddle

// Project captures one "fiddle". It will be stored in database.
// It is like the main of the whole thing.
type Project interface {
	FS() FileSystem
	ID() string
}

type Folder interface {
	Files() map[string][]byte
	Folders() map[string]Folder
}

// FileSystem is for storing all files that user is editing
type FileSystem interface {
	Root() Folder
	ReadFile(string) ([]byte, error)
	SaveFile(name string, content []byte) error
	DeleteFile(string) error
	DeleteFolder(string) error
}

type LCD interface {
	Write([]byte)
}

type KeyCode int

const (
	KeyCtrl KeyCode = iota
	KeyAlt
	KeyCommand
	KeyA
	KeyB
	KeyC
	Key1
	Key2
	KeyEsc
)

type KeyEvent interface {
	Code() KeyCode
	Down() bool
}

type KeyboardHandler func(KeyEvent)
type Keyboard interface {
	OnKey(KeyboardHandler)
}

type MouseCode int

const (
	MouseLeftDown MouseCode = iota
	MouseLeftUp
	MouseMove
)

type MouseEvent interface {
	Pos() (int, int)
	Code() MouseCode
}

type MouseHandler func(MouseEvent)
type Mouse interface {
	OnMouse(MouseHandler)
}

type Registers interface {
}

type Machine interface {
	Registers() Registers
	RAM() []byte
	Stack() (stack []int, pos int)
	Run()
}
