package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

const (
	ExitOk  = 0
	ExitErr = 1
)

type CLI struct {
	OutStream  io.Writer
	ErrStream  io.Writer
	CorrectNum int
}

func (c *CLI) Run(args []string) int {
	// オプションの処理
	var (
		filePath  string
		timeLimit int
	)

	flag.StringVar(&filePath, "f", "words.txt", "specify words file")
	flag.IntVar(&timeLimit, "t", 60, "specify time limit")
	flag.Parse()

	// ファイル読み込み
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("failed to open file: %v", err)
	}

	done := make(chan struct{})
	timer := time.NewTimer(time.Duration(timeLimit * int(time.Second)))

	go func(f *os.File) {
		// check input is correct
		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			target := scanner.Text()
			fmt.Fprint(c.OutStream, fmt.Sprintf("入力してください:  %v\n> ", target))
			stdin := bufio.NewScanner(os.Stdin)
			stdin.Scan()
			input := stdin.Text()

			if input == target {
				c.CorrectNum++
			}
		}
		done <- struct{}{}
	}(f)

	select {
	case <-timer.C:
		printResult := func() {
			fmt.Println("finish")
			fmt.Printf("limit time is over(%d seconds)\n", timeLimit)
			fmt.Fprintf(c.OutStream, "The number of your correct answers is %d\n", c.CorrectNum)
		}
		PrintWithDivLine(printResult)()
	case <-done:
		printResult := func() {
			fmt.Println("done all")
			fmt.Fprintf(c.OutStream, "The number of your correct answers is %d\n", c.CorrectNum)
		}
		PrintWithDivLine(printResult)()
	}

	return ExitOk
}

func PrintWithDivLine(f func()) func() {
	rf := func() {
		fmt.Println("------------")
		f()
		fmt.Println("------------")
	}
	return rf
}
