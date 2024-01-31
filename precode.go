package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Task ...
type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

// Обработчик для получения списка задач
func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	// Преобразование карты задач в список
	var taskList []Task
	for _, task := range tasks {
		taskList = append(taskList, task)
	}
	// Преобразование списка задач в формат JSON
	jsonTasks, err := json.Marshal(tasks)
if err != nil {
    log.Println("error marshaling tasks:", err)
    http.Error(w, "failed to marshal tasks", http.StatusInternalServerError)
    return
}

w.Header().Set("Content-Type", "application/json")
_, err = w.Write(jsonTasks)
if err != nil {
    log.Println("error writing response:", err)
    return
}
// Обработчик для получения информации о задаче по ее ID
func getTaskByIDHandler (w http.ResponseWriter, r *http.Request) {
	// Извлечение ID из URL-параметров запроса
	taskID := chi.URLParam(r, "id")
	// Поиск задачи по ID в карте задач
	task, ok := tasks[taskID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Преобразование задачи в формат JSON
	jsonTasks, err := json.Marshal(tasks)
if err != nil {
    log.Println("error marshaling tasks:", err)
    http.Error(w, "failed to marshal tasks", http.StatusInternalServerError)
    return
}

	w.Header().Set("Content-Type", "application/json")
_, err = w.Write(jsonTasks)
if err != nil {
    log.Println("error writing response:", err)
    return
}
}
// Обработчик для создания новой задачи
func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Чтение тела запроса в структуру Task
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	// Генерация нового ID для задачи
	newID := strconv.Itoa(len(tasks) + 1)

	// Добавление задачи в карту задач
	newTask.ID = newID
	tasks[newID] = newTask

	// Отправка ответа с созданной задачей в формате JSON
	jsonNewTask, err := json.Marshal(newTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(jsonNewTask)
}
// Обработчик для обновления существующей задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Извлечение ID из URL-параметров запроса
	taskID := chi.URLParam(r, "id")
	// Поиск задачи по ID в карте задач
	task, ok := tasks[taskID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Чтение тела запроса в структуру Task
	var updatedTask Task
	err := json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Обновление задачи в карте задач
	tasks[taskID] = updatedTask

	w.WriteHeader(http.StatusOK)
}

// Обработчик для удаления задачи по ее ID
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Извлечение ID из URL-параметров запроса
	taskID := chi.URLParam(r, "id")
	// Поиск задачи по ID в карте задач
	_, ok := tasks[taskID]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Удаление задачи из карты задач
	delete(tasks, taskID)

	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// Регистрация обработчиков
	r.Get("/tasks", getTasksHandler)
	r.Get("/tasks/{id}", getTaskByIDHandler)
	r.Post("/tasks", createTaskHandler)
	r.Put("/tasks/{id}", updateTaskHandler)
	r.Delete("/tasks/{id}", deleteTaskHandler)

	// запускаем сервер
	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}

