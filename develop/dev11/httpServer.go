package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

/*
Реализовать HTTP-сервер для работы с календарем. В рамках задания необходимо работать строго со стандартной HTTP-библиотекой.

В рамках задания необходимо:
Реализовать вспомогательные функции для сериализации объектов доменной области в JSON.
Реализовать вспомогательные функции для парсинга и валидации параметров методов /create_event и /update_event.
Реализовать HTTP обработчики для каждого из методов API, используя вспомогательные функции и объекты доменной области.
Реализовать middleware для логирования запросов


Методы API:
POST /create_event
POST /update_event
POST /delete_event
GET /events_for_day
GET /events_for_week
GET /events_for_month

Параметры передаются в виде www-url-form-encoded (т.е. обычные user_id=3&date=2019-09-09).
В GET методах параметры передаются через queryString, в POST через тело запроса.
В результате каждого запроса должен возвращаться JSON-документ содержащий либо {"result": "..."}
в случае успешного выполнения метода, либо {"error": "..."} в случае ошибки бизнес-логики.

В рамках задачи необходимо:
Реализовать все методы.
Бизнес логика НЕ должна зависеть от кода HTTP сервера.
В случае ошибки бизнес-логики сервер должен возвращать HTTP 503.
В случае ошибки входных данных (невалидный int например) сервер должен возвращать HTTP 400.
В случае остальных ошибок сервер должен возвращать HTTP 500.
Web-сервер должен запускаться на порту указанном в конфиге и выводить в лог каждый обработанный запрос.
*/

func main() {
	config, err := newConfig() // инициализируем конфиг переменную
	if err != nil {
		logrus.Errorln("Не удалось загрузить конфиг файл: ", err)
		return
	}

	store := newStore()                 // инициализируем хранилище
	router := newRouter(store, config)  // инициализируем роутер и передаем в него хранилище и конфиг
	server := newServer(config, router) // инициализируем сервер и передаем в него конфиг и роуте

	server.ListenAndServe() // слушаем порт
}

// newServer конструктор сервера
func newServer(conf *Config, router http.Handler) *http.Server {
	srv := &http.Server{
		Addr:         conf.Port, // параметры из файла окружения
		ReadTimeout:  conf.ReadTimeOut,
		WriteTimeout: conf.WriteTimeOut,
		Handler:      router, // в качестве хэндлера роутер
	}

	return srv
}

// newRouter конструктор роутера
func newRouter(Store *Store, conf *Config) http.Handler {
	mux := http.NewServeMux()

	handler := newHandler(Store, conf) // инициализируем основной хэндлер
	// POST запросы
	mux.HandleFunc("/create_event", handler.CreateEventHandler)
	mux.HandleFunc("/update_event", handler.UpdateEventHandler)
	mux.HandleFunc("/delete_event", handler.DeleteEventHandler)
	// GET запросы
	mux.HandleFunc("/events_for_day", handler.GetDayEvents)
	mux.HandleFunc("/events_for_week", handler.GetWeekEvents)
	mux.HandleFunc("/events_for_month", handler.GetMonthEvents)
	// Реализация middleware как логера запросов
	mw := newLogger(mux)

	// возвращаем промежуточный хэндлер
	// Logger реализует интерфейс http.Handler
	return mw
}

// Handler структура для реализации ручек
type Handler struct {
	Store *Store
	Conf  *Config
}

// newHandler конструктор ручек
func newHandler(store *Store, conf *Config) *Handler {
	return &Handler{
		Store: store, // передаем в хэндлеры хранилище и конфиг
		Conf:  conf,  // ибо так удобней
	}
}

// CreateEventHandler хэндлер для обработки создания события
func (h *Handler) CreateEventHandler(w http.ResponseWriter, r *http.Request) {
	var e Event

	err := e.Decode(r.Body) // декодим тело запроса
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = e.Validate() // делаем валидацию значений
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	e.Date = time.Now() // присваивает нынешнюю дату в запись

	err = h.Store.Create(&e) // записываем в хранилище
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	resultResponse(w, "Событие успешно создано", []*Event{&e}, http.StatusCreated)
}

// UpdateEventHandler хэндлер обработки запросов обновления событий
func (h *Handler) UpdateEventHandler(w http.ResponseWriter, r *http.Request) {
	// работает по аналогии с предыдущим
	var e Event

	err := e.Decode(r.Body)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = e.Validate()
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	e.Date = time.Now()

	err = h.Store.Update(&e)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	resultResponse(w, "Событие успешно обновлено", []*Event{&e}, http.StatusOK)
}

// DeleteEventHandler хэндлер для обработки запросов удаления событий
func (h *Handler) DeleteEventHandler(w http.ResponseWriter, r *http.Request) {
	// работает по аналогии с предыдущим за исключением отсутствия валидации
	var e Event

	err := e.Decode(r.Body)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	event, err := h.Store.Delete(&e)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	resultResponse(w, "Событие успешно удалено", []*Event{event}, http.StatusOK)
}

// GetDayEvents хэндлер для обработки запроса на все события за день
func (h *Handler) GetDayEvents(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id")) // берем id из query string
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	date, err := time.Parse(h.Conf.DateLayout, r.URL.Query().Get("date")) // оттуда же берем дату
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := h.Store.GetDayEvents(userID, date) // ищем события в хранилище
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	resultResponse(w, "Результаты за день", events, http.StatusOK)
}

// GetWeekEvents хэндлер для обработки запроса на все событий за неделю
func (h *Handler) GetWeekEvents(w http.ResponseWriter, r *http.Request) {
	// работает по аналогии с предыдущим
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	date, err := time.Parse(h.Conf.DateLayout, r.URL.Query().Get("date"))
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := h.Store.GetWeekEvents(userID, date)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	resultResponse(w, "Результаты за неделю", events, http.StatusOK)
}

// GetMonthEvents хэндлер для обработки запроса на все событий за месяц
func (h *Handler) GetMonthEvents(w http.ResponseWriter, r *http.Request) {
	// работает по аналогии с предыдущим
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	date, err := time.Parse(h.Conf.DateLayout, r.URL.Query().Get("date"))
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	events, err := h.Store.GetMonthEvents(userID, date)
	if err != nil {
		errorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	resultResponse(w, "Результаты за месяц", events, http.StatusOK)
}

// Logger структура для реализации логера
type Logger struct {
	Handler http.Handler
}

// newLogger конструктор логгера
func newLogger(mux http.Handler) *Logger {
	return &Logger{Handler: mux}
}

// ServeHTTP реализаций интерфейса http.Handler
func (l *Logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now() // нужно для подсчета занятого операцией времени

	l.Handler.ServeHTTP(w, r) // передаем дальше в следующий хэнлер

	logrus.Println(r.Method, r.URL.Path, time.Since(start)) // в стандартный вывод выводим логирования запросов
}

// Config структура для реализации конфига
type Config struct {
	Host         string
	Port         string
	ReadTimeOut  time.Duration
	WriteTimeOut time.Duration
	DateLayout   string
}

// newConfig конструктор
func newConfig() (*Config, error) {
	err := godotenv.Load("conf.env")
	if err != nil {
		return nil, err
	}

	readTemp := os.Getenv("READ_TO")
	readTO, err := time.ParseDuration(readTemp)
	if err != nil {
		return nil, err
	}

	writeTemp := os.Getenv("WRITE_TO")
	writeTO, err := time.ParseDuration(writeTemp)
	if err != nil {
		return nil, err
	}

	conf := &Config{
		Host:         os.Getenv("HOST"),
		Port:         os.Getenv("PORT"),
		ReadTimeOut:  readTO,
		WriteTimeOut: writeTO,
		DateLayout:   os.Getenv("DATE_LAYOUT"),
	}

	return conf, nil
}

// Event основная модель события
type Event struct {
	UserID      int       `json:"user_id"`
	EventName   string    `json:"event_name"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

// Decode функция декодирования из json-данных
func (e *Event) Decode(r io.Reader) error {
	err := json.NewDecoder(r).Decode(&e)
	if err != nil {
		return err
	}

	return nil
}

// Validate валидатор
func (e *Event) Validate() error {
	if e.UserID <= 0 {
		return errors.New("введите корректный User ID")
	}

	if e.EventName == "" {
		return errors.New("введите корректное имя события")
	}

	if e.Description == "" {
		return errors.New("введите корректное описание")
	}

	return nil
}

// Store структура для реализации хранилища
type Store struct {
	*sync.Mutex
	events map[int][]*Event
}

// newStore конструктор
func newStore() *Store {
	return &Store{
		Mutex:  &sync.Mutex{},
		events: make(map[int][]*Event),
	}
}

// Create метод создания новой записи события
func (s *Store) Create(e *Event) error {
	s.Lock()         // блокируем запись в переменную
	defer s.Unlock() // при выходе разблокируем

	events, ok := s.events[e.UserID] // проверяем наличие данных в мапе по юзер ид
	if ok {
		for _, event := range events {
			if event.EventName == e.EventName { // если название события уже есть, то возвращаем ошибку
				return errors.New("cобытие с таким названием уже существует")
			}
		}
	}

	// если же проверка прошла аппендим переданное аргументом событие в хранилище
	s.events[e.UserID] = append(s.events[e.UserID], e)

	return nil
}

// Update метод обновления событий
func (s *Store) Update(e *Event) error {
	s.Lock()
	defer s.Unlock()

	idx := -1

	events, ok := s.events[e.UserID] // проверяем наличие событий у пользователя
	if !ok {
		return errors.New("пользователь с таким id отсутствует") // если их нет возращаем ошибку
	}

	for i, event := range events { // итерируемся по событиям
		if event.EventName == e.EventName { // если находим событие с таким же именем как и у переданного события
			idx = i // присваиваем переменной idx номер события в слайсе
			break
		}
	}

	if idx == -1 { // если не нашли событие возращаем ошибку
		return errors.New("у данного пользователя отсутсвует событие с таким именем")
	}

	s.events[e.UserID][idx] = e // если пройдены все проверки присваиваем обновленное событие

	return nil
}

// Delete метод удаления событий
func (s *Store) Delete(e *Event) (*Event, error) {
	// работает аналогично предыдущему за исключением присвоения
	s.Lock()
	defer s.Unlock()

	idx := -1

	events, ok := s.events[e.UserID]
	if !ok {
		return nil, errors.New("такого пользователя не существует")
	}

	for i, event := range events {
		if event.EventName == e.EventName {
			idx = i
			break
		}
	}

	if idx == -1 {
		return nil, errors.New("у данного пользователя отсутствует событие с таким именем")
	}

	eventsLength := len(s.events[e.UserID])                      // нам нужна длина слайса с событиями
	deletedEvent := s.events[e.UserID][idx]                      // событие которое будем возращать для ответа пользователю
	s.events[e.UserID][idx] = s.events[e.UserID][eventsLength-1] // ну и здесь просто удаляем меняя значения с последним
	s.events[e.UserID] = s.events[e.UserID][:eventsLength-1]     // и обрезаем слайс на единицу

	return deletedEvent, nil
}

// GetDayEvents метод по поиску событий за день
func (s *Store) GetDayEvents(userID int, date time.Time) ([]*Event, error) {
	s.Lock()
	defer s.Unlock()

	var result []*Event // возвращаемый результат

	events, ok := s.events[userID] // ищем события по ид
	if !ok {                       // если их нет возращаем ошибку
		return nil, errors.New("пользователя с таким id не существует")
	}

	for _, event := range events { // итерируемся по событиям и сравниваем значения даты с значением в хранилище
		if event.Date.Year() == date.Year() && event.Date.Month() == date.Month() && event.Date.Day() == date.Day() {
			result = append(result, event) // если проверка пройдена аппендим событие в результат
		}
	}

	return result, nil
}

// GetWeekEvents метод по поиску событий за неделю
func (s *Store) GetWeekEvents(userID int, date time.Time) ([]*Event, error) {
	// работает по аналогии с предыдущим
	s.Lock()
	defer s.Unlock()

	var result []*Event

	events, ok := s.events[userID]
	if !ok {
		return nil, errors.New("пользователя с таким id не существует")
	}

	for _, event := range events {
		y1, w1 := event.Date.ISOWeek()
		y2, w2 := date.ISOWeek()
		if y1 == y2 && w1 == w2 {
			result = append(result, event)
		}
	}

	return result, nil
}

// GetMonthEvents метод по поиску событий за месяц
func (s *Store) GetMonthEvents(userID int, date time.Time) ([]*Event, error) {
	// Работает по аналогии с предыдущим
	s.Lock()
	defer s.Unlock()

	var result []*Event

	events, ok := s.events[userID]
	if !ok {
		return nil, errors.New("пользователя с таким id не существует")
	}

	for _, event := range events {
		if event.Date.Year() == date.Year() && event.Date.Month() == date.Month() {
			result = append(result, event)
		}
	}

	return result, nil
}

// errorRespose печать ошибки в ответе
func errorResponse(w http.ResponseWriter, er string, status int) {
	type Resp struct {
		Err string `json:"error"`
	}

	response := &Resp{Err: er}

	byteResp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(byteResp)
}

// resultResponse печать результата в ответе
func resultResponse(w http.ResponseWriter, re string, e []*Event, status int) {
	type Resp struct {
		Result string   `json:"result"`
		Events []*Event `json:"events"`
	}

	response := &Resp{
		Result: re,
		Events: e,
	}

	byteResp, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Write(byteResp)
}
