package filesin

import (
	"NAS/src/sqlmap"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
	"path"
)

var (
	redisCli = sqlmap.RedisClient
	myDb     = sqlmap.MyDb
	ctx      = sqlmap.Ctx
)

type (
	fileTom struct {
		UID      int64  `json:"UID"`               //鉴定权限
		UUID     string `json:"UUID" gorm:"index"` //唯一ID
		Name     string `json:"Name" gorm:"index"`
		Sha512   string `json:"Sha512"`           //验证文化完整性
		PID      string `json:"PID" gorm:"index"` //确定父系路径
		Size     int64  `json:"Size"`
		MIMEType string `json:"MIMEType" gorm:"index"`
		DiskName string `json:"DiskName"`
	}

	folderTom struct {
		UID  int64  `json:"UID"`
		UUID string `json:"UUID" gorm:"index"`
		PID  string `json:"PID"  gorm:"index"`
		Name string `json:"Name"`
	}
)

func (fileTom) TableName() string {
	return "Files"
}

func (folderTom) TableName() string {
	return "Folders"
}

func InitFile() {
	myDb.AutoMigrate(&fileTom{})
	myDb.AutoMigrate(&folderTom{})
	userDir("temp")
	go disk()
}

//AddFile 添加文件到mysql
func addFile(uid int64, name string, sha512 string, pid string, size int64, mime string, diskName string) error {
	var file fileTom
	file = fileTom{
		UID:      uid,
		UUID:     uuid.NewV4().String(),
		Name:     name,
		Sha512:   sha512,
		PID:      pid,
		Size:     size,
		MIMEType: mime,
		DiskName: diskName,
	}
	//myDb.Table("").Create(&file)
	return myDb.Table("files").Create(&file).Error
}

//AddFolder 添加虚拟文件夹
func addFolder(uid int64, name string, pid string) error {
	var fold folderTom
	fold = folderTom{
		UID:  uid,
		UUID: uuid.NewV4().String(),
		PID:  pid,
		Name: name,
	}
	return myDb.Table("folders").Create(&fold).Error
}

func delFiles(UUID string) {
	var dFile fileTom
	myDb.Model(&fileTom{}).Where("UUID = ?", UUID).Unscoped().Delete(&dFile)
}

func delFolder(UUID string) {
	var dFolder folderTom
	myDb.Model(&folderTom{}).Where("UUID = ?", UUID).Unscoped().Delete(&dFolder)
}

func upFiles(UUID string, newName string) {
	var todo fileTom
	todo = fileTom{Name: newName}
	myDb.Table("files").Model(&fileTom{}).Where("UUID = ?", UUID).Update(&todo)
}

func upFolder(UUID string, newName string) {
	var todo folderTom
	todo = folderTom{Name: newName}
	myDb.Table("folders").Model(&folderTom{}).Where("UUID = ?", UUID).Update(&todo)
}

func down(uid int64, uuid string) (nameDisk string, confirm bool) {
	var todo fileTom
	myDb.Table("files").First(&todo, "UID = ? AND UUID = ?", uid, uuid)
	if todo.UID == 0 {
		return "", false
	}
	return todo.DiskName, true
}

func getImage(uid int64, uuid string) (nameDisk string, confirm bool) {
	var todo fileTom
	myDb.Table("files").Where("UID = ? AND UUID = ?", uid, uuid).First(&todo)
	if todo.UID == 0 {
		return "", false
	}
	return todo.DiskName, true
}

func listSql(uid int64, pid string) ([]refFiles, []refFolders) {
	//var files []fileTom
	//var folders []folderTom
	var (
		refFile   []refFiles
		refFolder []refFolders
	)
	myDb.Table("files").Find(&refFile, "UID = ? AND p_id = ?", uid, pid)
	myDb.Table("folders").Find(&refFolder, "UID = ? AND p_id = ?", uid, pid)
	return refFile, refFolder
}

func getFiles(uid int64, sha512 string, size int64, name string, pid string) (str string, bol bool) {
	var todo fileTom
	myDb.Model(&fileTom{}).Where("Sha512 = ? AND Size = ?", sha512, size).First(&todo)
	if todo.UUID == "" {
		return "", false
	}
	if todo.PID == pid && todo.UID == uid {
		return "文件已存在", true
	}

	err := addFile(uid, name, todo.Sha512, pid, todo.Size, todo.MIMEType, todo.DiskName)
	if err != nil {
		log.Println(err.Error())
		return "奇怪的错误", true
	}
	if !getPid(uid, pid) {
		return "不合法的提交", true
	}
	return "上传完成", true
}

func getFolder(uid int64, pid string, name string) (str string, bol bool) {
	var fold folderTom
	myDb.Table("folders").First(&fold, "UID = ? AND p_id = ? AND Name = ?", uid, pid, name)
	if fold.UUID == "" {
		err := addFolder(uid, name, pid)
		if err != nil {
			log.Println("数据库写入失败", err.Error())
			return "数据库错误", false
		}
		return "新建成功", true
	}
	if fold.PID == pid && fold.Name == name {
		return "文件夹已存在", true
	}
	return "未知错误", false
}

func getFolderRename(uid int64, name string, pid string, new string) (str string, bol bool) {
	var fold folderTom
	myDb.Table("folders").First(&fold, "UID = ? AND Name = ? AND p_id = ?", uid, name, pid)
	if fold.UUID == "" {
		return "非法重命名提交", false
	}
	upFolder(fold.UUID, new)
	return "成功", true
}

func getFileRename(uid int64, name string, pid string, new string) (str string, bol bool) {
	var fol fileTom
	myDb.Table("files").First(&fol, "UID = ? AND Name = ? AND p_id = ?", uid, name, pid)
	if fol.UUID == "" {
		return "非法重命名提交", false
	}
	upFiles(fol.UUID, new)
	return "成功", true
}

func getDelFile(uid int64, name string, pid string) (str string, bol bool) {
	var fol fileTom
	myDb.Table("files").First(&fol, "UID = ? AND Name = ? AND p_id = ?", uid, name, pid)
	if fol.UUID == "" {
		return "文件不存在", false
	}
	delFiles(fol.UUID)
	go clear(fol.DiskName)
	return "已删除", true
}

func getDelFolder(uid int64, name string, pid string) (str string, bol bool) {
	var (
		fold folderTom
		file fileTom
	)
	myDb.Table("folders").First(&fold, "UID = ? AND Name = ? AND p_id = ?", uid, name, pid)
	if fold.UUID == "" {
		return "文件夹不存在", false
	}
	myDb.Table("files").First(&file, "UID = ? AND p_id = ?", uid, pid)
	if file.UUID == "" {
		delFolder(fold.UUID)
		return "已删除", true
	}
	return "文件夹下存在文件，删除失败", false
}

func clear(dn string) {
	var fold fileTom
	myDb.Table("files").First(&fold, "DiskName = ?", dn)
	if fold.UUID == "" {
		ext := path.Ext(dn)
		paths := dowDir(ext)
		remo(fmt.Sprintf("%s/%s", paths, dn))
	}
}

func getPid(uid int64, pid string) bool {
	var fold folderTom
	myDb.Table("folders").First(&fold, "UID = ? AND UUID = ?", uid, pid)
	if fold.UUID == "" {
		return false
	}
	return true
}
