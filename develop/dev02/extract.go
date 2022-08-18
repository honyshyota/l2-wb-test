package extract

import (
	"errors"
	"unicode"
)

/*
Создать Go-функцию, осуществляющую примитивную распаковку строки, содержащую повторяющиеся символы/руны, например:
"a4bc2d5e" => "aaaabccddddde"
"abcd" => "abcd"
"45" => "" (некорректная строка)
"" => ""

Дополнительно
Реализовать поддержку escape-последовательностей.
Например:
qwe\4\5 => qwe45 (*)
qwe\45 => qwe44444 (*)
qwe\\5 => qwe\\\\\ (*)

В случае если была передана некорректная строка, функция должна возвращать ошибку. Написать unit-тесты.
*/

var errIncorrectString = errors.New("некорректная строка")

// Функция распаковки
func Extract(s string) (string, error) {
	if len(s) == 0 { // Если длина строки равно 0 то возвращаем пустую строку без ошибки
		return "", nil
	} else if len(s) == 1 {
		return s, nil // Если длина строки равна 1 то возвращаем строку без ошибки
	}

	rString := []rune(s)  // Преобразуем строку в массив рун
	shield := false       // Флаг экрана
	prevSymIsDig := false // Флаг принадлежности предыдущего символа к цифре

	var result []rune // Переменная возращаемого функцией результата

	for i, val := range rString { // Итерируемся по нашему массиву рун
		if unicode.IsDigit(val) { // Если переменная цифра
			if i == 0 { // Если первый символ в массиве число возращаем пустую строку и ошибку
				return "", errIncorrectString
			} else if i != 0 {
				if !shield { // при отсутствии экранирования
					if !prevSymIsDig { // и предыдущий символ не число
						value := rString[i-1]             // берем предыдущий символ
						sRune := sliceLetters(value, val) // передаем его в функцию преобразователь, где на выход получаем массив рун
						result = append(result, sRune...) // аппендим в результат
						prevSymIsDig = true               // ставим флаг
					} else if prevSymIsDig { // если число
						return "", errIncorrectString // возвращаем пустую строку и ошибку
					}
				} else if shield { // с экранированием
					if !prevSymIsDig { // предыдущий символ не число
						prevSymIsDig = true // просто переключаем флаг
						if i == len(rString)-1 {
							result = append(result, val) // если в последней итерации цифра то просто записываем ее
						}
					} else if prevSymIsDig { // если число
						shield = false                    // выключаем флаг
						value := rString[i-1]             // берем предыдущий символ
						sRune := sliceLetters(value, val) // передаем в функцию преобразователь
						result = append(result, sRune...) // аппендим в результат
						prevSymIsDig = true               // переключаем флаг числа
					}
				}
			}
		} else if unicode.IsLetter(val) { // если символ буква
			if i != 0 { // Если перую букву не пропустим будет паника ибо берем предыдущий символ в массиве
				if i != len(rString)-1 { // Если индекс итерации не равен 						fmt.Println(5)длине нашего массива
					if !prevSymIsDig { // и предыдущий символ не число
						result = append(result, rString[i-1]) // просто аппендим предыдущий символ к результату
					} else if prevSymIsDig { // если предыдущий символ число просто переключаем флаг
						prevSymIsDig = false
					}
				} else if i == len(rString)-1 { // для последней итерации если там буква, в остальных случаях действий не требуется
					if !prevSymIsDig { // предыдущий символ не число
						result = append(result, rString[i-1], val) // аппендим два последних символа в результат
					} else if prevSymIsDig { // если число
						result = append(result, val) // аппендим только последний символ
					}
				}
			}
		} else if val == 92 { // если символ знак экранирования
			if prevSymIsDig { // предыдущий символ число
				prevSymIsDig = false // переключаем флаг числа
				if shield {          // если экран уже включен
					//shield = false // просто выключаем его
					result = append(result, rString[i-1])
				} else if !shield { // если выключен
					shield = true // просто включаем
				}
			} else if !prevSymIsDig { // предыдущий символ не число
				if shield { // экран включен
					shield = false // выключаем
				} else if !shield { // экран выключен
					result = append(result, rString[i-1]) // аппендим впредыдущий символ в результат
					shield = true                         // включаем экран
				}
			}
		}
	}

	return string(result), nil // в случае успеха возращаем результат и нулевую ошибку
}

// функция преобразователь
func sliceLetters(a rune, b rune) []rune {
	num := int(b - '0') // получаем инт из руны

	var result []rune

	for i := 0; i < num; i++ { // итерируемся до значения переданного инт
		result = append(result, a) // аппендим в результат переданное значение а в нужном колличестве
	}

	return result
}
