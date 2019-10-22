package vmwriter

import (
	"fmt"
	"os"
)

// VMWriter -- base struct
type VMWriter struct {
	output *os.File
}

// New -- create new vmwriter
func New(file *os.File) *VMWriter {
	return &VMWriter{output: file}
}

func (w *VMWriter) write(cmd string) {
	fmt.Fprintln(w.output, cmd)
}

// WritePush -- write push cmd
func (w *VMWriter) WritePush(segment string, index int) {
	w.write(fmt.Sprintf("push %s %d", segment, index))
}

// WritePop -- write pop cmd
func (w *VMWriter) WritePop(segment string, index int) {
	w.write(fmt.Sprintf("pop %s %d", segment, index))
}

// WriteArithmetic -- write arithmetic command
func (w *VMWriter) WriteArithmetic(cmd string) {
	w.write(cmd)
}

// WriteLabel -- write label command
func (w *VMWriter) WriteLabel(label string) {
	w.write(fmt.Sprintf("label %s", label))
}

// WriteGoto -- write go-to command
func (w *VMWriter) WriteGoto(label string) {
	w.write(fmt.Sprintf("goto %s", label))
}

// WriteIf -- write if-goto command
func (w *VMWriter) WriteIf(label string) {
	w.write(fmt.Sprintf("if-goto %s", label))
}

// WriteCall -- write call command
func (w *VMWriter) WriteCall(name string, nArgs int) {
	w.write(fmt.Sprintf("call %s %d", name, nArgs))
}

// WriteFunction -- write function command
func (w *VMWriter) WriteFunction(name string, nLocals int) {
	w.write(fmt.Sprintf("function %s %d", name, nLocals))
}

// WriteReturn -- write return command
func (w *VMWriter) WriteReturn() {
	w.write("return")
}
