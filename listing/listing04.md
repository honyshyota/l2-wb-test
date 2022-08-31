Что выведет программа? Объяснить вывод программы.


```
package main
 
func main() {
    ch := make(chan int)
    go func() {
        for i := 0; i < 10; i++ {
            ch <- i
        }
    }()
 
    for n := range ch {
        println(n)
    }
}
```

Вывод:

```
0
1
2
3
4
5
6
7
8
9
fatal error: all goroutines are asleep - deadlock!
```

Дедлок потому что читаем канал в который никто не пишет, его надо было закрыть