package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	numbers := make(chan int)

	go func() {
		i := 1
		for {
			select {
			case numbers <- i * i:
				i++
			}
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)

	go func() {
		<-sig
		fmt.Println("Выхожу из программы")
		os.Exit(0)
	}()

	for num := range numbers {
		fmt.Println(num)
	}
}
