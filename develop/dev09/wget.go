package main

import (
	"flag"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

/*
Реализовать утилиту wget с возможностью скачивать сайты целиком.
*/

func main() {
	var r bool // флаг по которому определяем загружать всю страницу ссылками или нет

	flag.BoolVar(&r, "r", false, "download full page")
	flag.Parse()

	argCheck := flag.Args()
	if len(argCheck) > 1 { // проверка ввода единственного агрумента
		logrus.Println("В качестве единственного аргумента введите URL-адрес")
		return
	}

	u := argCheck[0]

	addr, err := url.ParseRequestURI(u) // валидация URL адреса
	if err != nil {
		logrus.Println("Введите корректный URL-адрес", err)
		return
	}

	if !r { // если флага загрузки всей страницы нет
		err := pageDownload(addr.String())
		if err != nil {
			logrus.Fatalln("Не удалось загрузить страницу", err)
		}
	} else if r { // флаг есть
		err := fullSiteDownload(addr.String())
		if err != nil {
			logrus.Fatalln("Не удалось загрузить страницы", err)
		}
	}
}

// fullSiteDownload функция загрузки всей страницы
func fullSiteDownload(link string) error {
	resp, err := http.Get(link)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body) // для поиска ссылок использовал пакет от PuerkitoBio
	if err != nil {
		return err
	}

	var links []string

	doc.Find("body a").Each(func(index int, item *goquery.Selection) { // Ищем все ссылки на страницы
		linkTag := item
		link, _ := linkTag.Attr("href")
		links = append(links, link)
	})

	for _, link := range links { // итерируемся по слайсу ссылок
		addr, err := url.ParseRequestURI(link) // делаем валидацию ссылок, в случае провало пропускаем итерацию
		if err != nil {
			continue
		}

		err = pageDownload(addr.String()) // передаем по странично в функцию загрузки единчной страницы
		if err != nil {
			return err
		}
	}
	return nil
}

// pageDownload функция загрузки страницы
func pageDownload(link string) error {
	resp, err := http.Get(link)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	pathSlice := strings.Split(link, "/")
	downloadPath := "url/" + pathSlice[len(pathSlice)-2] // эта часть нужна чтоб просто сопоставить имена файлам

	file, err := os.Create(downloadPath) // создаем файл
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = io.Copy(file, resp.Body) // и просто записываем содержимое Body в него
	if err != nil {
		return err
	}

	return nil
}
