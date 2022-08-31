Что выведет программа? Объяснить вывод программы. Объяснить внутреннее устройство интерфейсов и их отличие от пустых интерфейсов.

```
package main
 
import (
    "fmt"
    "os"
)
 
func Foo() error {
    var err *os.PathError = nil
    return err
}
 
func main() {
    err := Foo()
    fmt.Println(err)
    fmt.Println(err == nil)
}
```

Вывод:

```
<nil>
false
```

Структура интерфейса состоит из двух значений ссылка, по сути, на реализацию (itab), которая в свою очередь хранит тип и методы этого типа
и второе значение данные (data) которое фактически указывает на конкретную переменную с конкретным типом, пустому интерфейсу соответсвует любой тип
посколько у данного интерфейса нет никаких методово и хранить их не нужно, поэтому itab не просчитывается и указателя на него нет,
есть только мета данные в data
В нашем случае data будет равно nil, а ссылка itable не равна nil, поэтому err != [nil, nil]