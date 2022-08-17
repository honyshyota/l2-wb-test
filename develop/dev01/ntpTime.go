package ntpTime

import (
	"io"
	"os"
	"time"

	"github.com/beevik/ntp"
	"github.com/sirupsen/logrus"
)

/*Создать программу печатающую точное время с использованием NTP -библиотеки.
Инициализировать как go module. Использовать библиотеку github.com/beevik/ntp.
Написать программу печатающую текущее время / точное время с использованием этой библиотеки.

Требования:
Программа должна быть оформлена как go module
Программа должна корректно обрабатывать ошибки библиотеки: выводить их в STDERR и возвращать ненулевой код выхода в OS
*/

// Требуемая функция инкапсулированы и имеет доступ из вне
func NtpTime() (time.Time, error) {
	time, err := ntp.Time("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		_, writeErr := io.WriteString(os.Stderr, err.Error()) // если возникает ошибка пишем ее в stderr
		if writeErr != nil {
			logrus.Println("Ошибка записи в stderr", writeErr)
			os.Exit(1)
		}
		os.Exit(1)
		return time, err // возвращаем ошибку
	}

	return time, nil // если все прошло нормально возвращаем нашу переменную time и нулевое значение вместо ошибки
}
