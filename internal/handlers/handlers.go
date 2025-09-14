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

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.Logger.Println("Ошибка парсинга формы:", err)
		http.Error(w, "Ошибка парсинга формы", http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("uploadfile")
	if err != nil {
		h.Logger.Println("Ошибка получения файла из формы:", err)
		http.Error(w, "Ошибка получения файла из формы: "+err.Error(), http.StatusInternalServerError)
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
		http.Error(w, "Ошибка конвертации: "+err.Error(), http.StatusInternalServerError)
		return
	}

	now := time.Now().UTC().Format("20060102_150405")
	ext := filepath.Ext(header.Filename)
	filename := now + ext

	outFile, err := createFile(filename)
	if err != nil {
		h.Logger.Println("Ошибка создания файла:", err)
		http.Error(w, "Ошибка создания файла", http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = outFile.Write([]byte(converted))
	if err != nil {
		h.Logger.Println("Ошибка записи в файл:", err)
		http.Error(w, "Ошибка записи в файл", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte(converted))
	if err != nil {
		h.Logger.Println("Ошибка записи ответа:", err)
	}
}

func createFile(name string) (*os.File, error) {
	return os.Create(name)
}
