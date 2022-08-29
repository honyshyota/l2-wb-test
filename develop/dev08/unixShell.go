package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/mitchellh/go-ps"
)

/*
Необходимо реализовать свой собственный UNIX-шелл-утилиту с поддержкой ряда простейших команд:

- cd <args> - смена директории (в качестве аргумента могут быть то-то и то)
- pwd - показать путь до текущего каталога
- echo <args> - вывод аргумента в STDOUT
- kill <args> - "убить" процесс, переданный в качесте аргумента (пример: такой-то пример)
- ps - выводит общую информацию по запущенным процессам в формате *такой-то формат*


Так же требуется поддерживать функционал fork/exec-команд

Дополнительно необходимо поддерживать конвейер на пайпах (linux pipes, пример cmd1 | cmd2 | .... | cmdN).

*Шелл — это обычная консольная программа, которая будучи запущенной, в интерактивном сеансе выводит некое приглашение
в STDOUT и ожидает ввода пользователя через STDIN. Дождавшись ввода, обрабатывает команду согласно своей логике
и при необходимости выводит результат на экран. Интерактивный сеанс поддерживается до тех пор, пока не будет введена команда выхода (например \quit).
*/

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := os.Getwd() // берем нынешний путь к папке чтобы подставить в терминал
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}
		fmt.Fprint(os.Stdout, line, "$ ") // печать пути по аналогии со стандартным терминалом

		input, err := reader.ReadString('\n') // читаем из ввода
		if err != nil {
			fmt.Fprint(os.Stderr, err)
		}
		input = strings.TrimSuffix(input, "\n") // здесь отсекаем конечный перенос
		input = strings.TrimSpace(input) // так же отсекаем лишние пробелы

		if strings.Contains(input, "|") { // если присутстует пайп лайн
			sliceCmd := strings.Split(input, "|") // сплитим строку по символу пайп лайна

			for _, cmd := range sliceCmd { // итерируемся по срезу из команд
				cmd = strings.TrimSuffix(cmd, "\n") // по тому же принципу отрезаем все лишнее
				cmd = strings.TrimSpace(cmd)
				if strings.Contains(cmd, "&") { // проверка на fork
					err := fork(cmd)
					if err != nil {
						fmt.Fprint(os.Stderr, err, "\n")
					}
				} else { 
					err := changer(cmd)
					if err != nil {
						fmt.Fprint(os.Stderr, err, "\n")
					}
				}
			}
		} else { // в ином случае просто выполняем команду
			if strings.Contains(input, "&") {
				err := fork(input)
				if err != nil {
					fmt.Fprint(os.Stderr, err, "\n")
				}
			} else {
				err = changer(input)
				if err != nil {
					fmt.Fprint(os.Stderr, err, "\n")
				}
			}
		}
	}
}

// changer основная функция, по своей сути прослойка, для выполнения команд
func changer(in string) error {
	args := strings.Split(in, " ") // сплитим по пробелу

	switch args[0] {
	case "cd":
		err := cd(args[1])
		if err != nil {
			return err
		}
	case "pwd":
		err := pwd(args[0])
		if err != nil {
			return err
		}
	case "echo":
		echo(args[1])
	case "ps":
		err := gops()
		if err != nil {
			return err
		}
	case "kill":
		err := kill(args[1])
		if err != nil {
			return err
		}
	case "exec":
		err := goexec(args[1])
		if err != nil {
			return err
		}
	case `\quit`:
		os.Exit(0)
	default:
		return errors.New("недопустимая команда")
	}
	return nil
}

// cd функция для запроса смены директории
func cd(in string) error {
	err := os.Chdir(in)
	if err != nil {
		return err
	}
	return nil
}

// pwd функция для запроса вывода адреса актуального каталога
func pwd(in string) error {
	line, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Fprint(os.Stdout, line, "\n")
	return nil
}

// echo функция эхо
func echo(in string) {
	fmt.Fprint(os.Stdout, in, "\n")
}

// kill функция для азпроса убить процесс
func kill(in string) error {
	pid, err := strconv.Atoi(in)
	if err != nil {
		return err
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	err = proc.Kill()
	if err != nil {
		return err
	}
	return nil
}

// gops функция для запроса вывода всех процесов
func gops() error {
	// в стандартных инструментах не нашел средств достижения цели
	// поэтому подключил внешнюю
	procces, err := ps.Processes()
	if err != nil {
		return err
	}
	for _, proc := range procces {
		fmt.Fprint(os.Stdout, proc.Pid(), " ", proc.PPid(), " ", proc.Executable(), "\n")
	}
	return nil
}

// exec команды
func goexec(in string) error {
	cmd := exec.Command(in)

	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// fork команды
func fork(in string) error {
	id, _, err := syscall.Syscall(syscall.SYS_FORK, 0, 0, 0)
	if err != 0 {
		os.Exit(1)
	} else if id == 0 {
		in = strings.TrimSuffix(in, "&")
		_, err := fmt.Fprintln(os.Stdout, os.Getpid())
		if err != nil {
			return err
		}
		err = changer(in)
		if err != nil {
			return err
		}
		_, err = fmt.Fprintln(os.Stdout, "Завершен")
		if err != nil {
			return err
		}
		os.Exit(0)
	}
	return nil
}
