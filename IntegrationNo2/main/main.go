package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите сообщение: ")
	message, _ := reader.ReadString('\n')
	message = message[:len(message)-1]

	currentTime := time.Now()
	date := currentTime.Format("2006-01-02 15:04:05")

	file, err := os.OpenFile("output.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Ошибка открытия файла:", err)
		return
	}
	defer file.Close()

	line := fmt.Sprintf("%s %s\n", date, message)
	_, err = file.WriteString(line)
	if err != nil {
		fmt.Println("Ошибка записи в файл:", err)
		return
	}

	fmt.Println("Данные успешно записаны в файл.")
}
