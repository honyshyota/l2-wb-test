Что выведет программа? Объяснить вывод программы. 
Рассказать про внутреннее устройство слайсов и что происходит при передаче их в качестве аргументов функции.

```
package main
 
import (
  "fmt"
)
 
func main() {
  var s = []string{"1", "2", "3"}
  modifySlice(s)
  fmt.Println(s)
}
 
func modifySlice(i []string) {
  i[0] = "3"
  i = append(i, "4")
  i[1] = "5"
  i = append(i, "6")
}
```

Вывод:

```
[3 2 3]
```

Так как слайсы передаются по значению в первой строке modifySlice() присваиваем 0 позиции значение 3, функция аппенд создает новый слайс, и дальше мы работаем с ним внутри функции никак не влияя на значения слайса из мэйн функции
Структура слайса это ссылка на изначальный массив, длина слайса и вместимость