package main

import (
	"fmt"
	"strconv"
	"sync"
)

func squareWorker(input <-chan int, output chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range input {
		squared := num * num
		output <- squared
		fmt.Println("Квадрат: ", squared)
	}
	close(output)
}

func multiplyWorker(input <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for num := range input {
		result := num * 2
		fmt.Println("Произведение:", result)
	}
}

func main() {
	inputCh := make(chan int)
	squareCh := make(chan int)
	var wg sync.WaitGroup
	wg.Add(2)

	go squareWorker(inputCh, squareCh, &wg)
	go multiplyWorker(squareCh, &wg)

	for {
		var input string
		fmt.Scanln(&input)

		if input == "стоп" || input == "stop" {
			break
		}

		fmt.Println("Ввод: ", input)

		num, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Некорректный ввод. Введите число или 'стоп'/'stop")
			continue
		}

		inputCh <- num
	}

	close(inputCh)
	wg.Wait()
}
