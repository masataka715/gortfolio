package chat

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

func UploaderHandler(w http.ResponseWriter, req *http.Request) {
	userId := req.FormValue("userid")
	file, header, err := req.FormFile("avatarFile")
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	filename := filepath.Join("chat/avatars", userId+filepath.Ext(header.Filename))
	err = ioutil.WriteFile(filename, data, 0777)
	log.Println(filename)
	if err != nil {
		io.WriteString(w, err.Error())
		return
	}
	// io.WriteString(w, "成功")
	w.Header()["Location"] = []string{"/chat"}
	w.WriteHeader(http.StatusMovedPermanently)
}
