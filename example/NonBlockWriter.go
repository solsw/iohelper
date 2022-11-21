package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/solsw/iohelper"
)

type errWriter struct {
	w io.Writer
	c int
}

func (ew *errWriter) Write(p []byte) (int, error) {
	time.Sleep(500 * time.Millisecond)
	ew.c++
	if ew.c == 5 {
		return 0, errors.New("ERROR: errWriter.c == 5")
	}
	return ew.w.Write(p)
}

func example() {
	fmt.Println("example start")
	nbw := iohelper.NewNonBlockWriter(
		// os.Stdout,
		&errWriter{w: os.Stdout},
		0,
		func(error) bool {
			return true
		},
	)
	defer func() {
		if err := nbw.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	log.SetOutput(nbw)
	for i := 1; i <= 8; i++ {
		log.Println(i)
	}
	fmt.Println("example end")
}

func main() {
	example()
	fmt.Println("OK")
}
