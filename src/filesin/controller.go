package filesin

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"time"
)

func init() {
	exists, _ := pathExists("./file")
	if !exists {
		_ = os.Mkdir("./file", 600)
	}
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func userDir(ext string) string {
	//var path = "./file/" + strconv.FormatInt(uid, 10) + "a" + "/" + ext
	var pathDir = "./file/" + ext
	pathExistsBool, _ := pathExists(pathDir)
	if !pathExistsBool {
		errMkdir := os.Mkdir(pathDir, 600)
		if errMkdir != nil {
			log.Println(errMkdir)
		}
	}
	return pathDir
}

func dowDir(ext string) string {
	return "./file/" + ext
}

func def(file *os.File) {
	err := file.Close()
	if err != nil {
		log.Println("Clone Err", err.Error())
	}
}

func tempImageDisk(diskName string, uuid string, uid int64) string {
	ext := path.Ext(diskName)
	imgData, _ := ioutil.ReadFile(dowDir(ext) + "/" + diskName)
	buf := bytes.NewBuffer(imgData)
	image, err := imaging.Decode(buf)
	if err != nil {
		fmt.Println("图片处理", err.Error())
		return ""
	}
	//生成缩略图，尺寸宽100,高传0表示等比例放缩
	//最后缩略图尺寸为100*133
	sprintf := fmt.Sprintf("%s/%s", "./file/temp", strconv.FormatInt(time.Now().Unix(), 10)+ext)
	image = imaging.Resize(image, 100, 0, imaging.Lanczos)
	err = imaging.Save(image, sprintf)
	if err != nil {
		fmt.Println("缩略图保存失败", err.Error())
	}
	redisCli.Set(ctx, uuid+strconv.FormatInt(uid, 10), sprintf, time.Duration(0))
	return sprintf
}
