package filesin

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

func remo(name string) bool {
	errRemove := os.Remove(name)
	if errRemove != nil {
		log.Println("删除", errRemove.Error())
		return false
	}
	return true
}

//func remName(oldName string, newName string) string {
//	errRename := os.Rename(oldName, newName)
//	if errRename != nil {
//		log.Println("重命名", errRename.Error())
//	}
//	return newName
//}
//
////验证文件是否相同优先使用sha512
//func compareByte(sFile *os.File, dFile *os.File) bool {
//	var sByte []byte = make([]byte, 512)
//	var dByte []byte = make([]byte, 512)
//	var sErr, dErr error
//	for {
//		_, sErr = sFile.Read(sByte)
//		_, dErr = dFile.Read(dByte)
//		if sErr != nil || dErr != nil {
//			if sErr != dErr {
//				return false
//			}
//			if sErr == io.EOF {
//				break
//			}
//		}
//		if bytes.Equal(sByte, dByte) {
//			continue
//		}
//		return false
//	}
//	return true
//}

func sha512File(filepath string) string {
	file, errOpen := os.Open(filepath)
	if errOpen != nil {
		return ""
	}
	defer def(file)
	m := sha512.New()
	_, errIo := io.Copy(m, file)
	if errIo != nil {
		return ""
	}
	return hex.EncodeToString(m.Sum(nil))
}

//处理长时间未上传完成的文件
func diskClear() {
	var cursor uint64
	val, _ := redisCli.Scan(ctx, cursor, "*", 10000).Val()
	for _, key := range val {
		if find := strings.Contains(key, "-Total"); find {
			i, err := redisCli.Get(ctx, key).Int64()
			if err != nil {
				log.Println("Clear Redis", err.Error())
			}
			strings.TrimRight(key, "-Total")
			i2, err := redisCli.Get(ctx, key).Int64()
			if err != nil {
				log.Println("Clear Redis", err.Error())
			}
			ext := path.Ext(key)
			sprintf := fmt.Sprintf("%s/%s/%s", "./file", ext, key)
			size := Size(sprintf)
			if i == i2 {
				if i == size {
					log.Println("文件大小与Redis相等")
					//fileSha := sha512File(sprintf)
					clear(key)
					redisCli.Del(ctx, key)
					redisCli.Del(ctx, key+"-Total")
				}
				if i2 == size {
					log.Println("删除长时间未处理的文件")
					boo := remo(sprintf)
					if boo {
						redisCli.Del(ctx, key)
						redisCli.Del(ctx, key+"-Total")
					} else {
						log.Println("删除失败")
					}
				}
			}
		}
	}
}

func Size(name string) int64 {
	var (
		n   int
		sum int64 = 0
	)
	file, err := os.Open(name)
	if err == nil {
		buf := make([]byte, 2014)
		for {
			n, err = file.Read(buf)
			sum += int64(n)
			if err == io.EOF {
				break
			}
		}
		return sum
	}
	return sum
}

func disk() {
	for {
		ticker := time.NewTicker(time.Hour * 24)
		<-ticker.C
		log.Println("计划任务")
		diskClear()
	}
}
