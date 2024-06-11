package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

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

// Ниже напишите обработчики для каждого эндпоинта

// обработчик для получения всех задача
func getAllTasks(w http.ResponseWriter, req *http.Request) {
	response, err := json.Marshal(tasks)
	if err != nil {
		fmt.Printf("Ошибка сериализации: %v", err.Error())
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		fmt.Printf("Ошибка при записи ответа: %v", err)
		return
	}
}

// обработчик для отправки задачи на сервер
func postTask(w http.ResponseWriter, req *http.Request) {
	var task Task
	var buf bytes.Buffer
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		fmt.Printf("Ошибка при чтении тела: %v", err.Error())
		http.Error(w, "Ошибка сервера", http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		fmt.Printf("Ошибка десериализации: %v", err.Error())
		http.Error(w, "Ошибка сервера", http.StatusBadRequest)
		return
	}
	_, ok := tasks[task.ID]
	if ok {
		http.Error(w, "Задача уже существует", http.StatusBadRequest)
		return
	}
	tasks[task.ID] = task
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

// обработчик для получения задачи по ID
func getTaskByID(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	t, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}
	response, err := json.Marshal(t)
	if err != nil {
		fmt.Printf("Ошибка сериализации: %v", err.Error())
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(response)
	if err != nil {
		fmt.Printf("Ошибка при записи ответа: %v", err)
		return
	}
}

// обработчик удаления задачи по ID
func deleteTaskById(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	_, ok := tasks[id]
	if !ok {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}
	delete(tasks, id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", getAllTasks)
	r.Post("/tasks", postTask)
	r.Get("/tasks/{id}", getTaskByID)
	r.Delete("/tasks/{id}", deleteTaskById)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
