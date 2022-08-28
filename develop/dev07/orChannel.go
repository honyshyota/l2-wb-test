package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// Функция создает канал, в аргумент передается время через которое канал закроется
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	// or принимает аргументами любое колличество каналов, и возращает один, который закроется
	// при завершении любого из переданных каналов
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(5*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)

	// просто печать результата времени выполнения по которому можно отследить время работы функции
	fmt.Printf("fone after %v \n", time.Since(start))

}

func or(inChannels ...<-chan interface{}) <-chan interface{} {
	retChan := make(chan interface{}) // создаем результирующий канал

	var wg sync.WaitGroup // группа ожидания
	wg.Add(1)

	for _, ch := range inChannels { // итерация по входящим каналам
		go func(ch <-chan interface{}) { // анонимная функция слушает каждый канал на предмет завершения и передачи wg.Done()
			for range ch { // внутри функции запускаем цикл при завершении которого обнуляем счетчик группы ожидания
			}
			wg.Done()
		}(ch)
	}

	wg.Wait() // здесь просто ждем

	close(retChan) // закрываем результирующий канал
	return retChan // передаем его в качестве возращаемого значения функцией
}
