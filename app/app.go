package app

import (
	"os"
	"syscall"
	"time"
)

type List struct {
	Description   string
	CreatedAt     time.Time
	FormattedTime string
	Done          bool
}

type Todo struct {
	List []List
}

func NewTodo() *Todo {
	return &Todo{}
}

func closeFile(f *os.File) {
	// Unlock the file
	syscall.Flock(int(f.Fd()), syscall.LOCK_UN)

	// Close the file
	_ = f.Close()
}
