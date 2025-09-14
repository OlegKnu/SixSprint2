package handlers

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

type Handler struct {
	Logger *log.Logger
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	f, err := os.Open("index.html")
	if err != nil {
		h.Logger.Println("Ошибка открытия index.html:", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := io.Copy(w, f); err != nil {
		h.Logger.Println("Ошибка отправки index.html:", err)
	}
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.Logger.Println("Ошибка парсинга формы:", err)
		http.Error(w, "Ошибка парсинга формы", http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("myFile")
	if err != nil {
		h.Logger.Println("Ошибка получения файла из формы:", err)
		http.Error(w, "Ошибка получения файла из формы", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		h.Logger.Println("Ошибка чтения файла:", err)
		http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
		return
	}

	converted, err := service.Convert(string(data))
	if err != nil {
		h.Logger.Println("Ошибка конвертации:", err)
		http.Error(w, "Ошибка конвертации", http.StatusInternalServerError)
		return
	}

	ext := filepath.Ext(header.Filename)
	filename := time.Now().UTC().Format("20060102_150405") + ext

	f, err := os.Create(filename)
	if err != nil {
		h.Logger.Println("Ошибка создания файла:", err)
		http.Error(w, "Ошибка создания файла", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(converted); err != nil {
		h.Logger.Println("Ошибка записи в файл:", err)
		http.Error(w, "Ошибка записи в файл", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(converted))
}
