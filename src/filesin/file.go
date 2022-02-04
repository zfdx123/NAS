package filesin

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"
)

type (
	update struct {
		Name     string `json:"name"`
		Size     int64  `json:"size"`
		Sha512   string `json:"sha512"`
		PID      string `json:"PID"`
		MIMEType string `json:"MIMEType"`
	}

	mess struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Offset  int64  `json:"Offset,omitempty"`
	}

	//download struct {
	//	UUID string `json:"UUID"`
	//	Name string `json:"Name"`
	//}

	list struct {
		File   interface{}
		Folder interface{}
	}

	newFolder struct {
		Name string `json:"Name"`
		PID  string `json:"PID"`
	}
	reName struct {
		New  string `json:"New"`
		Name string `json:"Name"`
		PID  string `json:"PID"`
	}
	delFl struct {
		Name string `json:"Name"`
		Pid  string `json:"Pid"`
	}

	//tempImage struct {
	//	Name string `json:"Name"`
	//	Disk string `json:"Disk"`
	//}
)

// UpdateLoad 确认接口
func UpdateLoad(c *gin.Context) {
	uid := c.GetInt64("uid")
	var upd update
	err := c.BindJSON(&upd)
	if err != nil {
		c.JSON(http.StatusBadRequest, mess{
			Status:  -1,
			Message: "解析数据失败",
		})
		return
	}
	files, bol := getFiles(uid, upd.Sha512, upd.Size, upd.Name, upd.PID)
	if bol {
		c.JSON(http.StatusOK, mess{
			Status:  1,
			Message: files,
		})
		return
	} else {
		mar, _ := json.Marshal(upd)
		redisCli.Set(ctx, upd.Name+strconv.FormatInt(uid, 10), mar, time.Hour*24)
		val, err := redisCli.Get(ctx, upd.Name).Int64()
		if err != nil {
			val = 0
		}
		c.JSON(http.StatusOK, mess{
			Status:  0,
			Message: "处理上传",
			Offset:  val,
		})
	}
}

func Upload(c *gin.Context) {
	var (
		upd    update
		total  int64 = 0
		buff   int64 = 1024 * 8
		result int64
		read   int
	)
	uid := c.GetInt64("uid")
	formFile, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		log.Println("FormFile", err.Error())
		panic(err)
		//log.Println(err.Error())
	}
	redisCli.Set(ctx, fileHeader.Filename+"-Total", fileHeader.Size, time.Duration(0))
	ext := path.Ext(fileHeader.Filename)
	paths := userDir(ext)
	saveFile, openErr := os.OpenFile(fmt.Sprintf("%s/%s", paths, fileHeader.Filename), os.O_RDWR, os.ModePerm)
	if openErr != nil {
		log.Println("OpenFile", openErr.Error())
		saveFile, err = os.Create(fmt.Sprintf("%s/%s", paths, fileHeader.Filename))
		if err != nil {
			log.Println("Create", err)
		}
		//defer def(saveFile)
	} else {
		log.Println("续传", openErr.Error())
		result, err = redisCli.Get(ctx, fileHeader.Filename).Int64()
		if err == nil {
			total = result
			log.Println("获得文件上次存储位置", total)
			_, err = formFile.Seek(total, io.SeekStart)
			if err != nil {
				log.Println("Seek失败", err.Error())
			}
		}
	}

	buf := make([]byte, buff)

	for {
		read, err = formFile.Read(buf)
		if err != nil && err != io.EOF {
			log.Println("读取错误", err)
		}
		if read == 0 {
			break
		}
		if read < int(buff) {
			//buf = make([]byte, read)
			buf = buf[:read]
		}
		_, err = saveFile.WriteAt(buf, total)
		if err != nil {
			log.Println("WriteAt错误", err.Error())
		}
		total += int64(read)
		err = redisCli.Set(ctx, fileHeader.Filename, total, time.Duration(0)).Err()
		if err != nil {
			log.Println("Redis 保存进度失败", err.Error())
		}
	}

	err = saveFile.Close()
	if err != nil {
		log.Println("关闭数据流失败", err.Error())
	}

	val := redisCli.Get(ctx, fileHeader.Filename+strconv.FormatInt(uid, 10)).Val()
	sha512 := sha512File(fmt.Sprintf("%s/%s", paths, fileHeader.Filename))
	err = json.Unmarshal([]byte(val), &upd)
	if err != nil {
		log.Println("JSON", err.Error())
		go remo(fmt.Sprintf("%s/%s", paths, fileHeader.Filename))
		c.JSON(http.StatusOK, mess{
			Status:  1,
			Message: "数据校验失败，文件上传失败！JSON不合法",
		})
		return
	}

	if total == upd.Size {
		redisCli.Del(ctx, fileHeader.Filename)
		//redisCli.Del(fileHeader.Filename + "Ones")
		if upd.Sha512 != "" && upd.Sha512 == sha512 {
			redisCli.Del(ctx, fileHeader.Filename+strconv.FormatInt(uid, 10))
			c.JSON(http.StatusOK, mess{
				Status:  0,
				Message: "上传成功",
			})
			err = addFile(uid, fileHeader.Filename, sha512, upd.PID, fileHeader.Size, ext, fileHeader.Filename)
			if err != nil {
				log.Println("数据库写入", err.Error())
			}
			return
		} else {
			go remo(fmt.Sprintf("%s/%s", paths, fileHeader.Filename))
			c.JSON(http.StatusOK, mess{
				Status:  -1,
				Message: "数据校验失败，文件上传失败！",
			})
			return
		}
	} else {
		c.JSON(http.StatusOK, mess{
			Status:  2,
			Message: "文件上传未完成",
		})
		return
	}
}

// DownloadFileService Test资源文件下载
func DownloadFileService(c *gin.Context) {
	//var downs download
	//tok := c.Request.Header.Get("token")
	//uid := toke(tok)
	//err := c.BindJSON(&downs)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, mess{
	//		Status:  -1,
	//		Message: "数据解析失败",
	//	})
	//	return
	//}
	uuid := c.Param("uuid")
	uid := c.GetInt64("uid")
	myName, confirm := down(uid, uuid)
	if !confirm {
		c.JSON(http.StatusOK, mess{
			Status:  1,
			Message: "无文件",
		})
		return
	}
	ext := path.Ext(myName)
	//file, errOpenDown := os.Open(fmt.Sprintf("%s/%s", dowDir(ext), myName))
	//if errOpenDown != nil {
	//	log.Println("打开错误", errOpenDown.Error())
	//	c.JSON(http.StatusOK, mess{
	//		Status:  1,
	//		Message: "无文件",
	//	})
	//	return
	//}
	//file.Close()
	fileName := path.Base(myName)
	fileName = url.QueryEscape(fileName) // 防止中文乱码
	//c.Header("Content-Type", "application/octet-stream")
	c.Header("content-disposition", "attachment; filename=\""+fileName+"\"")
	c.File(fmt.Sprintf("%s/%s", dowDir(ext), myName))
}

func TempImages(c *gin.Context) {
	var imageDisk string
	uid := c.GetInt64("uid")
	uuid := c.Param("uuid")
	imageDisk = redisCli.Get(ctx, uuid+strconv.FormatInt(uid, 10)).Val()
	if imageDisk == "" {
		nameDisk, confirm := getImage(uid, uuid)
		if confirm {
			imageDisk = tempImageDisk(nameDisk, uuid, uid)
		}
	}
	c.File(imageDisk)
}

func ReNameFile(c *gin.Context) {
	var fol reName
	uid := c.GetInt64("uid")
	err := c.BindJSON(&fol)
	if err != nil {
		c.JSON(http.StatusOK, mess{
			Status:  -1,
			Message: "数据解析错误",
		})
		return
	}
	rename, bol := getFileRename(uid, fol.Name, fol.PID, fol.New)
	if bol {
		c.JSON(http.StatusOK, mess{
			Status:  0,
			Message: rename,
		})
	} else {
		c.JSON(http.StatusOK, mess{
			Status:  -1,
			Message: rename,
		})
	}
}

func DelFile(c *gin.Context) {
	var fol delFl
	uid := c.GetInt64("uid")
	err := c.BindJSON(&fol)
	if err != nil {
		c.JSON(http.StatusOK, mess{
			Status:  -1,
			Message: "数据解析错误",
		})
		return
	}
	rename, bol := getDelFile(uid, fol.Name, fol.Pid)
	if bol {
		c.JSON(http.StatusOK, mess{
			Status:  0,
			Message: rename,
		})
	} else {
		c.JSON(http.StatusOK, mess{
			Status:  -1,
			Message: rename,
		})
	}
}
