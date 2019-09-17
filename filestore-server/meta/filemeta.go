package meta

import (
	"fmt"
	"sort"

	mydb "github.com/bigpengry/netdisc/filestore-server/db"
)

// FileMeta : 文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta : 新增/更新文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
	fmt.Println(fileMetas)
}

// UpdateFileMetaDB : 新增/更新文件元信息到数据库
func UpdateFileMetaDB(fmeta FileMeta) bool {
	return mydb.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.Location, fmeta.FileSize)
}

// GetFileMeta : 通过sha1值获取文件元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// GetFileMetaDB : 从MySQL获取文件元信息
func GetFileMetaDB(filesha1 string) (FileMeta, error) {
	fmeta := FileMeta{}
	tfile, err := mydb.GetFileMeta(filesha1)
	if err != nil {
		return fmeta, err
	}
	fmeta = FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		Location: tfile.FileAddr.String,
		FileSize: tfile.FileSize.Int64,
	}
	return fmeta, nil
}

//RemoveFileMeta : 删除文件元信息操作
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}

// GetLastFileMetas : 获取多个文件的元信息
func GetLastFileMetas(count int) []FileMeta {
	fMetaArray := make([]FileMeta, len(fileMetas))
	i := 0
	for _, v := range fileMetas {
		fMetaArray[i] = v
		i++
	}
	sort.Sort(ByUploadTime(fMetaArray))
	return fMetaArray[0:count]
}
