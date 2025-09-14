package handlers

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
	"log"
)

type Handler struct {
	Logger *log.Logger
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.Logger.Println("Ошибка парсинга формы:", err)
		http.Error(w, "Ошибка парсинга формы", http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.Logger.Println("Не удалось получить файл из формы:", err)
		http.Error(w, "Ошибка загрузки файла", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		h.Logger.Println("Ошибка чтения файла:", err)
		http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
		return
	}

	result, err := service.Convert(string(data))
	if err != nil {
		h.Logger.Println("Ошибка конвертации:", err)
		http.Error(w, "Ошибка конвертации данных", http.StatusInternalServerError)
		return
	}

	timeStr := time.Now().UTC().Format("20060102T150405Z")
	ext := filepath.Ext(header.Filename)
	filename := timeStr + ext

	newFile, err := os.Create(filename)
	if err != nil {
		h.Logger.Println("Ошибка создания файла:", err)
		http.Error(w, "Ошибка сервера при сохранении файла", http.StatusInternalServerError)
		return
	}
	defer newFile.Close()

	_, err = newFile.WriteString(result)
	if err != nil {
		h.Logger.Println("Ошибка записи в файл:", err)
		http.Error(w, "Ошибка сервера при записи файла", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(result))
}
