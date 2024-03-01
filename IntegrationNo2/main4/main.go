package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func writeToFile(input string) {
	file, err := os.OpenFile("messages.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Ошибка при открытии файла:", err)
		return
	}
	defer file.Close()

	if input == "exit" {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	message := timestamp + " " + input + "\n"

	if _, err := file.WriteString(message); err != nil {
		fmt.Println("Ошибка при записи в файл:", err)
	}
}

func readFromFile() {
	data, err := ioutil.ReadFile("messages.txt")
	if err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return
	}

	if len(data) == 0 {
		fmt.Println("Файл пуст")
		return
	}

	fmt.Println(string(data))
}

func main() {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Введите сообщение: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		writeToFile(input)

		if input == "exit" {
			break
		}
	}

	readFromFile()
}
