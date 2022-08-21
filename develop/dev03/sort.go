package main

import (
	"flag"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

/*
Отсортировать строки в файле по аналогии с консольной утилитой sort
(man sort — смотрим описание и основные параметры):
на входе подается файл из несортированными строками, на выходе — файл с отсортированными.

Реализовать поддержку утилитой следующих ключей:

-k — указание колонки для сортировки (слова в строке могут выступать в качестве колонок, по умолчанию разделитель — пробел)
-n — сортировать по числовому значению
-r — сортировать в обратном порядке
-u — не выводить повторяющиеся строки

Дополнительно

Реализовать поддержку утилитой следующих ключей:

-M — сортировать по названию месяца
-b — игнорировать хвостовые пробелы
-c — проверять отсортированы ли данные
-h — сортировать по числовому значению с учетом суффиксов

*/

type flags struct {
	k int
	n bool
	r bool
	u bool
}

// Флаги инициализировал без конструктора потому что задание выполняется одним файлом
// и запутаться в глобальных переменных будет сложно, ну и сам парсинг аргументов запускается
// через init()
var f flags

func init() {
	flag.IntVar(&f.k, "k", 1, "указание колонки для сортировки")
	flag.BoolVar(&f.n, "n", false, "сортировать по числовому значению")
	flag.BoolVar(&f.r, "r", false, "сортировать в обратном порядке")
	flag.BoolVar(&f.u, "u", false, "не выводить повторяющиеся строки")

}

func main() {
	flag.Parse() // Парсим аргументы

	var buf []byte

	fName := flag.Arg(0)
	if fName == "" {
		logrus.Println("Не указано имя файла")
		return
	} else {
		f, err := ioutil.ReadFile(fName)
		if err != nil {
			logrus.Println("Ошибка открытия файла")
			return
		}

		buf = f
	}

	sSlice := mainSort(buf, &f)

	file, err := os.Create("output.txt")
	if err != nil {
		logrus.Fatalln("Ошибка создания файла")
	}
	defer file.Close()

	for _, val := range sSlice {
		_, err := file.WriteString(val + "\n")
		if err != nil {
			logrus.Fatalln("Ошибка записи в файл")
		}
	}
}

func simpleSort(input []byte, flag *flags) []string {
	stringSlice := strings.Split(string(input), "\n")
	if flag.r {
		sort.Sort(sort.Reverse(sort.StringSlice(stringSlice)))
	} else {
		sort.Strings(stringSlice)
	}

	return stringSlice
}

func columnSort(input []byte, flag *flags) []string {
	stringSlice := strings.Split(string(input), "\n")

	keysStore := make(map[string][]string)
	keys := make([]string, 0, len(stringSlice))

	for i, val := range stringSlice {
		s := strings.Split(val, " ")

		var column string
		if flag.k <= len(s) {
			column = s[flag.k-1]
		} else {
			logrus.Println("Колонка в строке", i, "отсутствует")
			break
		}

		if _, ok := keysStore[column]; !ok {
			keys = append(keys, column)
		}

		keysStore[column] = append(keysStore[column], val)
	}

	if flag.r {
		sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	} else {
		sort.Strings(keys)
	}

	result := make([]string, 0, len(stringSlice))

	for _, key := range keys {
		result = append(result, keysStore[key]...)
	}

	return result
}

func numberSort(input []byte, flag *flags) []string {
	stringSlice := strings.Split(string(input), "\n")

	keysStore := make(map[float64][]string)
	keys := make([]float64, 0, len(stringSlice))

	for _, val := range stringSlice {
		s := strings.Split(val, " ")

		column, err := strconv.ParseFloat("-inf", 32)
		if err != nil {
			logrus.Fatal(err)
		}

		if flag.k <= len(s) {
			column, err = strconv.ParseFloat(s[flag.k-1], 32)
			if err != nil {
				column, err = strconv.ParseFloat("-inf", 32)
				if err != nil {
					logrus.Fatal(err)
				}
			}
		}

		if _, ok := keysStore[column]; !ok {
			keys = append(keys, column)
		}

		keysStore[column] = append(keysStore[column], val)
	}

	if flag.r {
		sort.Sort(sort.Reverse(sort.Float64Slice(keys)))
	} else {
		sort.Float64s(keys)
	}

	result := make([]string, 0, len(stringSlice))

	for _, key := range keys {
		result = append(result, keysStore[key]...)
	}

	return result
}

func deleteDuplicate(input []byte) []string {
	stringSlice := strings.Split(string(input), "\n")

	keysStore := make(map[string]struct{})

	for _, val := range stringSlice {
		if _, ok := keysStore[val]; !ok {
			keysStore[val] = struct{}{}
		}
	}

	result := make([]string, 0, len(stringSlice))

	for key := range keysStore {
		result = append(result, key)
	}

	return result
}

func mainSort(input []byte, flag *flags) []string {
	result := make([]string, 0, len(input))

	if flag.k != 1 {
		result = append(result, columnSort(input, flag)...)
	} else if flag.n {
		result = append(result, numberSort(input, flag)...)
	} else if flag.u {
		result = append(result, deleteDuplicate(input)...)
	} else {
		result = append(result, simpleSort(input, flag)...)
	}

	return result
}
