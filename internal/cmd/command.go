package cmd

import "io"

type Command interface {
	Run() error
	Start() error
	CombinedOutput() ([]byte, error)
	Output() ([]byte, error)
	Wait() error
	StdinPipe() (io.WriteCloser, error)
	StderrPipe() (io.ReadCloser, error)
}
