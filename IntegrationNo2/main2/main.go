package main

import (
	"fmt"
	"os"
)

func main() {
	fileName := "test.txt"

	fileInfo, err := os.Stat(fileName)
	if os.IsNotExist(err) || fileInfo.Size() == 0 {
		fmt.Println("Файл отсутствует или пуст")
		return
	}

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Ошибка при попытке открыть файл:", err)
		return
	}
	defer file.Close()

	buffer := make([]byte, fileInfo.Size())
	_, err = file.Read(buffer)
	if err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return
	}

	fmt.Println(string(buffer))
}
