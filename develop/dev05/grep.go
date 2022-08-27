package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

/*Реализовать утилиту фильтрации по аналогии с консольной утилитой (man grep — смотрим описание и основные параметры).

Реализовать поддержку утилитой следующих ключей:
-A - "after" печатать +N строк после совпадения
-B - "before" печатать +N строк до совпадения
-C - "context" (A+B) печатать ±N строк вокруг совпадения
-c - "count" (количество строк)
-i - "ignore-case" (игнорировать регистр)
-v - "invert" (вместо совпадения, исключать)
-F - "fixed", точное совпадение со строкой, не паттерн
-n - "line num", напечатать номер строки
*/

func main() {
	flags := flagParse() // парсим аргументы и флаги командной строки
	args := argsParse()

	sliceFileText, err := readFile(args.fName) // считываем файл с адрессом указанным в аргументе
	if err != nil {
		logrus.Fatalln("файл не найден")
	}

	result, match := searchText(args, flags, sliceFileText) // обрабатываем слайс строк с учетом переданных аргументов

	if match != 0 { // если совпаденией 0 то выводим "совпадений не найдено", если совпадения есть то выводим их в stdout
		if flags.count {
			fmt.Println("Совпадений: ", match)
			fmt.Println(result)
		} else {
			fmt.Println(result)
		}
	} else {
		fmt.Println("Совпадений не найдено")
	}

}

// структура пользовательского типа для обработки флагов аргументов
type flags struct {
	after   int
	before  int
	context int
	count   bool
	ignore  bool
	invert  bool
	fixed   bool
	lineNum bool
}

// конструктор флагов
func newFlags() *flags {
	return &flags{}
}

// функция обработки флагов
func flagParse() *flags {
	flags := newFlags()
	flag.IntVar(&flags.after, "A", 0, "'after' печатать +N строк после совпадения")
	flag.IntVar(&flags.before, "B", 0, "'before' печатать +N строк до совпадения")
	flag.IntVar(&flags.context, "C", 0, "'context' (A+B) печатать ±N строк вокруг совпадения")
	flag.BoolVar(&flags.count, "c", false, "'count' (количество строк)")
	flag.BoolVar(&flags.ignore, "i", false, "'ignore-case' (игнорировать регистр)")
	flag.BoolVar(&flags.invert, "v", false, "'invert' (вместо совпадения, исключать)")
	flag.BoolVar(&flags.fixed, "F", false, "'fixed', точное совпадение со строкой, не паттерн")
	flag.BoolVar(&flags.lineNum, "n", false, "'line num', напечатать номер строки")
	flag.Parse()

	if flags.context > 0 { // чтобы отдельно не обрабатывать условие контекст можно передать нужные значения в after и before
		flags.after, flags.before = flags.context, flags.context
	}

	return flags
}

// пользовательский тип для обработки аргументов командной строки
type args struct {
	phrase string
	fName  string
}

// конструктор для типа args
func newArgs() *args {
	return &args{}
}

// функция для обработки аргументов
func argsParse() *args {
	args := newArgs() // Инициализуруем переменную пользовательского типа

	fArgs := flag.Args() // парсим аргументы

	if len(fArgs) < 2 { // если длина слайса аргументов меньше двух элементов, выводим хелп сообщение
		logrus.Println("Недостаточно аргументов, нужно: [flags] [phrase] [file name]")
	}

	slicePhrase := fArgs[:len(fArgs)-1]          // из всех аргументов объедниям в слайс все, кроме последнего элемента
	args.phrase = strings.Join(slicePhrase, " ") // объединяем в одну строку и присваиваем в переменную пользовательского типа
	args.fName = fArgs[len(fArgs)-1]             // последний элемент должен быть адресом файла, поэтому напрямую присваиваем его пользовательскому типу

	return args
}

// пользовательский тип который реализует хранилище для текста из файла
type storeText struct {
	store map[int]string
}

// конструктор
func newStoreText() *storeText {
	return &storeText{
		store: make(map[int]string),
	}
}

// функция чтения из файла
func readFile(fileName string) (*storeText, error) {
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	sliceString := strings.Split(string(file), "\n") // сплитим содержимое файла по символу переноса строки

	storeString := newStoreText() // инициализируем переменную хранилища

	for index, value := range sliceString { // итерируемся по слайсу
		storeString.store[index] = value // пишем в мапу где индексом будет номер строки, значением строка
	}

	return storeString, nil
}

// основная фукция реализующая базовый алгоритм grep
func searchText(arg *args, flag *flags, text *storeText) (string, int) {
	var condition bool // переключатель для выявления совпадений
	var result string  // результат
	var match int      // колличество совпадений, нужно для флага count

	// итерируемся по хранилищу текста из файла
	for index, value := range text.store {
		if flag.ignore { // если стоит флаг ignore приводим все к нижнему регистру
			value = strings.ToLower(value)
			arg.phrase = strings.ToLower(arg.phrase)
		}

		if flag.fixed { // флаг точного совпадения со строкой
			condition = arg.phrase == value
		} else {
			condition = strings.Contains(value, arg.phrase)
		}

		if condition { // если переключатель активен
			match++          // инкреминируем совпадения
			if flag.invert { // флаг исключения
				if flag.after > 0 && flag.before > 0 { // флаг контекст, либо до и после
					for i := 1; i < index-flag.before; i++ {
						if flag.lineNum {
							result += strconv.Itoa(i) + "." + text.store[i] + "\n"
						} else {
							result += text.store[i] + "\n"
						}
					}
					for i := index + flag.after + 1; i < len(text.store); i++ {
						if flag.lineNum {
							result += strconv.Itoa(i) + "." + text.store[i] + "\n"
						} else {
							result += text.store[i] + "\n"
						}
					}
				} else if flag.after > 0 { // флаг после
					for i := 1; i < index; i++ {
						if flag.lineNum {
							result += strconv.Itoa(i) + "." + text.store[i] + "\n"
						} else {
							result += text.store[i] + "\n"
						}
					}
					for i := index + flag.after + 1; i < len(text.store); i++ {
						if flag.lineNum {
							result += strconv.Itoa(i) + "." + text.store[i] + "\n"
						} else {
							result += text.store[i] + "\n"
						}
					}
				} else if flag.before > 0 { // флаг до
					for i := 1; i < index-flag.before; i++ {
						if flag.lineNum {
							result += strconv.Itoa(i) + "." + text.store[i] + "\n"
						} else {
							result += text.store[i] + "\n"
						}
					}
					for i := index + 1; i < len(text.store); i++ {
						if flag.lineNum {
							result += strconv.Itoa(i) + "." + text.store[i] + "\n"
						} else {
							result += text.store[i] + "\n"
						}
					}
				}
			} else { // если флага исключения нет делаем все аналогично, только наоборот
				if flag.after > 0 && flag.before > 0 {
					for i := index - flag.before; i < index+flag.after; i++ {
						if flag.lineNum {
							result += strconv.Itoa(i) + "." + text.store[i] + "\n"
						} else {
							result += text.store[i] + "\n"
						}
					}
				} else if flag.after > 0 {
					for i := index; i < index+flag.after; i++ {
						if flag.lineNum {
							result += strconv.Itoa(i) + "." + text.store[i] + "\n"
						} else {
							result += text.store[i] + "\n"
						}
					}
				} else if flag.before > 0 {
					for i := index - flag.before; i < index; i++ {
						if flag.lineNum {
							result += strconv.Itoa(i) + "." + text.store[i] + "\n"
						} else {
							result += text.store[i] + "\n"
						}
					}
				} else {
					if flag.lineNum {
						result += strconv.Itoa(index) + "." + value + "\n"
					} else {
						result += value + "\n"
					}
				}
			}
		}
	}

	return result, match
}
