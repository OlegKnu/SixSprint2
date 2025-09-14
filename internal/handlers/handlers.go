package handlers

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона: "+err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Ошибка отображения шаблона: "+err.Error(), http.StatusInternalServerError)
	}
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Ошибка парсинга формы: "+err.Error(), http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("uploadfile")
	if err != nil {
		http.Error(w, "Ошибка получения файла из формы: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Ошибка чтения файла: "+err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := service.Convert(string(data))
	if err != nil {
		http.Error(w, "Ошибка конвертации: "+err.Error(), http.StatusInternalServerError)
		return
	}

	ext := filepath.Ext(header.Filename)
	filename := time.Now().UTC().Format("20060102_150405") + ext

	outFile, err := os.Create(filename)
	if err != nil {
		http.Error(w, "Ошибка создания файла: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer outFile.Close()

	_, err = outFile.WriteString(result)
	if err != nil {
		http.Error(w, "Ошибка записи в файл: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(result))
}
