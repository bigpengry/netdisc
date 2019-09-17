package handler

import (
	"encoding/json"
	"github.com/bigpengry/netdisc/filestore-server/meta"
	"github.com/bigpengry/netdisc/filestore-server/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

// UploadHandler : 处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		// 返回上传html页面
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internel server error")
			return
		}
		io.WriteString(w, string(data))
	}
	if r.Method == "POST" {
		// 接受文件流上传到本地目录
		file, header, err := r.FormFile("file")
		defer file.Close()
		if err != nil {
			fmt.Printf("Failed to get data,err:%s\n", err.Error())
			return
		}
		fileMeta := meta.FileMeta{
			FileName: header.Filename,
			Location: "./tmp/" + header.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		newFile, err := os.Create(fileMeta.Location)
		defer newFile.Close()
		if err != nil {
			fmt.Printf("Failed to create file,err:%s\n", err.Error())
			return
		}
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to save data,err:%s\n", err.Error())
			return
		}

		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		//meta.UpdateFileMeta(fileMeta)
		_ = meta.UpdateFileMetaDB(fileMeta)
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

// UploadSucHandler : 上传完成
func UploadSucHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Upload finished!")
}

// GetFileMetaHandler : 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileHash := r.Form.Get("filehash")
	fileMeta, err := meta.GetFileMetaDB(fileHash)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(fileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(data)
}

// QueryFileHandler : 文件查询接口
func QueryFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	limitCnt, _ := strconv.Atoi(r.Form.Get("limit"))
	fileMetas := meta.GetLastFileMetas(limitCnt)
	data, err := json.Marshal(fileMetas)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Write(data)
}

// DownloadHandler : 文件下载接口
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileHash := r.Form.Get("filehash")
	fileMeta := meta.GetFileMeta(fileHash)

	file, err := os.Open(fileMeta.Location)
	defer file.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("content-disposition", "attachment; filename=\""+fileMeta.FileName+"\"")
	w.Write(data)
}

// UpdateFileMetaHandler : 更新文件元信息接口
func UpdateFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	opType := r.Form.Get("op")
	fileHash := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	curFileMeta := meta.GetFileMeta(fileHash)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)

	w.WriteHeader(http.StatusOK)
}

// DeleteFileHandler : 删除文件元信息接口
func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	fileHash := r.Form.Get("fileHash")
	fm := meta.GetFileMeta(fileHash)
	os.Remove(fm.Location)

	meta.RemoveFileMeta(fileHash)

	w.WriteHeader(http.StatusOK)
}
