package stdio

import (
	"bufio"
	"fmt"
	"os"
)

var Stdin = os.Stdin
var Stdout = os.Stdout

func SetStdin(f *os.File) {
	Stdin = f
}

func SetStdout(f *os.File) {
	Stdout = f
}

func Print(a ...any) {
	fmt.Fprint(Stdout, a...)
}

func Printf(format string, a ...any) {
	fmt.Fprintf(Stdout, format, a...)
}

func Scan(a ...any) {
	fmt.Fscan(Stdin, a...)
}

func Read(p []byte) (int, error) {
	return Stdin.Read(p)
}

func ReadString(delim byte) (string, error) {
	return bufio.NewReader(Stdin).ReadString(delim)
}
