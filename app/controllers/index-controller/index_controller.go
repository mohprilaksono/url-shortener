package index_controller

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/mohprilaksono/url-shortener/config"
	"github.com/mohprilaksono/url-shortener/utils"
	"github.com/mohprilaksono/url-shortener/utils/logs"
	"github.com/mohprilaksono/url-shortener/utils/str"
)

var templates *template.Template = template.Must(template.New("").ParseGlob(filepath.Join("views", "*")))

func Index(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	err := templates.ExecuteTemplate(buf, "index.html", nil)
	if err != nil {
		logs.Err(err.Error())
		return
	}

	w.Header().Add("Content-Type", "text/html")

	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}

func Store(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		logs.Err(err.Error())
		return
	}

	uri := r.Form.Get("url")

	file, err := utils.LoadFile()
	if err != nil && err != io.EOF {
		logs.Err(err.Error())
		return
	}

	defer utils.CloseFile(file)

	numberOfBytesWritten := atomic.LoadInt64(&config.NumberOfBytesWritten)
	key := str.Random()
	value := url.PathEscape(uri)
	numOfBytesWritten, err := file.WriteAt([]byte(fmt.Sprintf("%s|%s*", key, value)), numberOfBytesWritten)
	if err != nil {
		logs.Err(err.Error())
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}

	atomic.AddInt64(&config.NumberOfBytesWritten, int64(numOfBytesWritten))

	http.Redirect(w, r, "/show/" + key, http.StatusMovedPermanently)
}

func Show(w http.ResponseWriter, r *http.Request) {
	param := r.PathValue("url")

	file, err := utils.LoadFile()
	if err != nil && err != io.EOF {
		logs.Err(err.Error())
		return
	}

	defer utils.CloseFile(file)

	bufReader := bufio.NewReader(file)

	var result = map[string]string{
		"Key": "",
		"Value": "",
	}
	
	for {
		data, err := bufReader.ReadString('*')
		if err == io.EOF {
			break
		}

		key, value, _ := strings.Cut(data, "|")
		if key == param {
			value, err = url.PathUnescape(value[:len(value) - 1])
			if err != nil {
				logs.Err(err.Error())
				return
			} 

			result["Key"] = fmt.Sprintf("%s/%s", r.Host, key)
			result["Value"] = value
			break
		}
	}

	buf := new(bytes.Buffer)
	err = templates.ExecuteTemplate(buf, "show.html", result)
	if err != nil {
		logs.Err(err.Error())
		return
	}

	w.Header().Add("Content-Type", "text/html")

	w.WriteHeader(http.StatusOK)
	buf.WriteTo(w)
}

func Go(w http.ResponseWriter, r *http.Request) {
	uri := r.PathValue("url")

	file, err := utils.LoadFile()
	if err != nil && err != io.EOF {
		logs.Err(err.Error())
		return
	}

	defer utils.CloseFile(file)

	var result string
	
	buf := bufio.NewReader(file)

	for {
		data, err := buf.ReadString('*')
		if err == io.EOF {
			break
		}

		key, value, _ := strings.Cut(data, "|")
		if key == uri {
			value, err = url.PathUnescape(value[:len(value) - 1])
			if err != nil {
				logs.Err(err.Error())
				return
			}

			result = value
			break
		}
	}

	http.Redirect(w, r, result, http.StatusMovedPermanently)
}