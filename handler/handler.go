package handler

import (
	"encoding/json"
	"fmt"
	"go-cloud/meta"
	"go-cloud/util"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func UploadHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			fmt.Println("err:", err)
			_, _ = io.WriteString(w, "internal server error")
			return
		}
		_, _ = io.WriteString(w, string(data))
	} else if r.Method == "POST" {
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get data,err:%s\n", err.Error())
			return
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "/Users/cyf/Pictures/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Println("failed to create file,err:", err)
		}
		defer newFile.Close()
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("io copy err:", err)
			return
		}
		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		fmt.Println("sha1", fileMeta.FileSha1)
		meta.UpdataFileMeta(fileMeta)

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "Upload finished")
}

func GetFIleMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form["filehash"][0]
	fMate := meta.GetFileMeta(filehash)
	data, err := json.Marshal(fMate)
	if err != nil {
		fmt.Println("json marshal err:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func FileQueryHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))
	fMate := meta.GetLastFileMetas(limitCnt)
	data, err := json.Marshal(fMate)
	if err != nil {
		fmt.Println("json marshal err:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)

	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("content-disposition", "attachment; filename=\""+fm.FileName+"\"")
	w.Write(data)

}

func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request) {
	_ = r.ParseForm()
	opType := r.Form.Get("op")
	fileSha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	meta.UpdataFileMeta(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func FileMetaDeleteHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fileSha1 := r.Form.Get("filehash")
	fileMeta := meta.GetFileMeta(fileSha1)
	os.Remove(fileMeta.Location)
	meta.RemoveFileMeta(fileSha1)
	w.WriteHeader(http.StatusOK)
}
