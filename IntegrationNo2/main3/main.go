package main

import (
	"fmt"
	"os"
)

func main() {
	file, err := os.Create("readonly.txt")
	if err != nil {
		fmt.Println("Ошибка создания файла:", err)
		return
	}
	defer file.Close()

	err = file.Chmod(0400)
	if err != nil {
		fmt.Println("Ошибка настройки прав доступа к файлу:", err)
		return
	}

	fmt.Println("Файл создан и установлены права только для чтения.")
	writeFile, err := os.OpenFile("readonly.txt", os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Ошибка открытия файла для записи:", err)
		return
	}
	defer writeFile.Close()

	_, err = writeFile.WriteString("Это тест.")
	if err != nil {
		fmt.Println("Ошибка записи файла", err)
		return
	}
}
