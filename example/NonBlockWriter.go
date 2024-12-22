package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/solsw/iohelper"
)

type sleepWriter struct {
	w io.Writer
}

func (sw *sleepWriter) Write(p []byte) (int, error) {
	time.Sleep(300 * time.Millisecond)
	return sw.w.Write(p)
}

type errSleepWriter struct {
	w io.Writer
	c int
}

func (esw *errSleepWriter) Write(p []byte) (int, error) {
	time.Sleep(300 * time.Millisecond)
	esw.c++
	if esw.c == 5 {
		return 0, errors.New("ERROR: errSleepWriter.c == 5")
	}
	return esw.w.Write(p)
}

func example() {
	fmt.Println("example start")
	// ctx, cancel := context.WithTimeout(context.Background(), 1000*time.Millisecond)
	// defer cancel()
	nbw := iohelper.NewNonBlockWriter(
		context.Background(),
		// ctx,

		// os.Stdout,
		// &sleepWriter{w: os.Stdout},
		&errSleepWriter{w: os.Stdout},

		0,
		func(error) bool { return true },
	)
	defer func() {
		if err := nbw.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	log.SetOutput(nbw)
	for i := 1; i <= 16; i++ {
		log.Println(i)
	}
	fmt.Println("example end")
}

func main() {
	example()
	fmt.Println("OK")
}
