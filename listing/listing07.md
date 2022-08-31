Что выведет программа? Объяснить вывод программы.

```
package main
 
import (
    "fmt"
    "math/rand"
    "time"
)
 
func asChan(vs ...int) <-chan int {
   c := make(chan int)
 
   go func() {
       for _, v := range vs {
           c <- v
           time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
      }
 
      close(c)
  }()
  return c
}
 
func merge(a, b <-chan int) <-chan int {
   c := make(chan int)
   go func() {
       for {
           select {
               case v := <-a:
                   c <- v
              case v := <-b:
                   c <- v
           }
      }
   }()
 return c
}
 
func main() {
 
   a := asChan(1, 3, 5, 7)
   b := asChan(2, 4 ,6, 8)
   c := merge(a, b )
   for v := range c {
       fmt.Println(v)
   }
}
```

Вывод:

```
1
2
3
4
5
6
7
8
0
0
0...
```

Так как данные закончились функция asChan закрыла оба канала, а в функции по объединению каналов в один в select нет проверки на закрытость канала, и поэтому она читает из закрытых каналов, они в свою очередь по дефолту отдают 0