package anagramSearch

import (
	"sort"
	"strings"
)

func FindAnagrams(input *[]string) *map[string][]string {
	result := make(map[string][]string) // результат
	// мапа где ключем является сортированые по возрастанию буквы, а значеним слово
	sortedSlice := make(map[string]string)

	for _, val := range *input { // итерируемся по массиву переданному в аргументе
		if len(val) < 2 { // пропускаем слова из одной буквы
			continue
		}
		finder(val, &sortedSlice, &result) // передаем данные в нашу функцию поиска
	}

	for key, val := range result { // удаляем множества из одного слова
		if len(val) == 1 {
			delete(result, key)
		}
	}

	return &result // возвращаем указатель на результат
}

func finder(s string, sorted *map[string]string, result *map[string][]string) {
	toLower := strings.ToLower(s) // перевод слова в нижний регистр

	//toRune := []rune(toLower) // приводим получившуюся строку в тип рун

	sortedRune := []rune(toLower) // переменная с отсортированными рунами

	sort.Slice(sortedRune, func(i, j int) bool { return sortedRune[i] <= sortedRune[j] }) // сортируем руны в порядке возрастания

	// проверяем есть ли в нашей мапе, с сортированными буквами по возрастанию, ключ
	// если его нет добавляем такой ключ в мапу, в качестве значения добавляем искомое слово,
	// так же записываем данное слово одновременно и в ключ и в значение в мапе результирующей
	// а если такой сортированый порядок есть в мапе, то аппендим слово в результирующую мапу
	val, ok := (*sorted)[string(sortedRune)]
	if ok {
		(*result)[val] = append((*result)[val], toLower)
	} else {
		(*sorted)[string(sortedRune)] = toLower
		(*result)[toLower] = append((*result)[toLower], toLower)
	}
}
