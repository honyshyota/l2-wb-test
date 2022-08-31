Что выведет программа? Объяснить вывод программы.

```
package main
 
type customError struct {
     msg string
}
 
func (e *customError) Error() string {
    return e.msg
}
 
func test() *customError {
     {
         // do something
     }
     return nil
}
 
func main() {
    var err error
    err = test()
    if err != nil {
        println("error")
        return
    }
    println("ok")
}
```

Вывод:

```
error
```

Такой вывод потому что Test возращает ссылку на структуру реализующую интерфейс error, и чтобы 
быть равной nil надо чтоб и динамическая часть и статическая должны быть равны nil