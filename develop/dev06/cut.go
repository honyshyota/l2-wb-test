package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

/*
Реализовать утилиту аналог консольной команды cut (man cut).
Утилита должна принимать строки через STDIN,
разбивать по разделителю (TAB) на колонки и выводить запрошенные.

Реализовать поддержку утилитой следующих ключей:
-f - "fields" - выбрать поля (колонки)
-d - "delimiter" - использовать другой разделитель
-s - "separated" - только строки с разделителем
*/

func main() {
	flags := flagsParse() // парсим флаги

	err := readFromStdin(flags) // читаем из ввода
	if err != nil {
		logrus.Fatalln("невозможно прочитать пользовательский ввод, ", err)
	}
}

// пользовательский тип реализации флагов
type flags struct {
	fields    int
	delimiter string
	separated bool
}

// конструктор
func newFlags() *flags {
	return &flags{}
}

// функция парсинга флагов
func flagsParse() *flags {
	flags := newFlags() // инициализируем переменную флагов

	// реализуем запись в поля переменной нужных флагов
	flag.IntVar(&flags.fields, "f", 0, "'fields' - выбрать поля (колонки)")
	flag.StringVar(&flags.delimiter, "d", "\t", "'delimiter' - использовать другой разделитель")
	flag.BoolVar(&flags.separated, "s", false, "'separated' - только строки с разделителем")
	flag.Parse()

	return flags
}

// readFromStdin функция для считывания пользовательского ввода и обработки последующей
func readFromStdin(f *flags) error {
	for {
		reader := bufio.NewReader(os.Stdin)
		text, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		// если флаг -s и строка не содержит выбраный разделитель
		if f.separated && !strings.Contains(text, f.delimiter) {
			fmt.Println("")
			continue
		}

		// сплитим строку выбраным разделителем (по умолчания TAB)
		splitString := strings.Split(text, f.delimiter)
		if f.fields <= len(splitString) { // так же фильтруем по флагу -f
			fmt.Println("запрошенная колонка", f.fields, ": ", splitString[f.fields])
			continue
		}

		return nil
	}
}
