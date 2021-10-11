package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	config "main/config"
	"main/internal/pkg/e/auths"

	"golang.org/x/crypto/bcrypt"
	// mysql
	uuid "github.com/satori/go.uuid"
	mysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// User 用户表
type User struct {
	gorm.Model
	UID      string `gorm:"type:varchar(255)"`
	Username string `gorm:"type:varchar(255)"`
	Password string `gorm:"type:varchar(255)"`
	RoleID   uint   // belong to Role
	Role     Role   `gorm:"foreignKey:RoleID"`
}

// Role 角色表
type Role struct {
	gorm.Model
	Name     string   `gorm:"type:varchar(255)"` // 角色名
	Desc     string   `gorm:"type:varchar(255)"` // 角色描述
	AuthInfo AuthInfo `gorm:"type:json"`         // 权限信息
}

// //////////////// Auth json //////////////////
// AuthInfo 权限信息 json
type AuthInfo struct {
	Auth []Auth `json:"auths"`
}

// Auth 权限
type Auth struct {
	Flag     string   `json:"flag"`     // 标志
	Desc     string   `json:"desc"`     // 描述
	Resource []string `json:"resource"` // 资源
}

// Scan authInfo反序列化
func (authInfo *AuthInfo) Scan(v interface{}) error {
	err := json.Unmarshal(v.([]byte), authInfo)
	return err
}

// Value authInfo序列化
func (authInfo AuthInfo) Value() (driver.Value, error) {
	b, err := json.Marshal(authInfo)
	fmt.Println(string(b))
	return string(b), err
}

// //////////////////////////////////

// Token token
type Token struct {
	gorm.Model
	Token    string    `gorm:"type:varchar(255)"`
	ExpireAt time.Time // 过期时间
	Enable   bool      // 是否生效
	UID      string    // belong to User
	User     User      `gorm:"foreignKey:UID;references:UID"`
}

// BFile .b
type BFile struct {
	gorm.Model
	Filename    string `gorm:"type:varchar(255)"`
	FileSize    uint
	Desc        string `gorm:"type:varchar(255)"`
	IsUse       bool
	File        File `gorm:"foreignKey:BFileID"`
	UploaderUID string
	User        User `gorm:"foreignKey:UploaderUID;references:UID"`
}

// File 保存文件
type File struct {
	gorm.Model
	Data    []byte `gorm:"type:bytes"`
	Hash    string `gorm:"type:varchar(255)"`
	BFileID uint
}

// Task 任务
type Task struct {
	gorm.Model
	Key         string `gorm:"type:varchar(255)"`
	Type        string `gorm:"type:varchar(255)"`
	Status      uint   // 0.发起 1.成功 2.失败
	OperatorUID string
	User        User `gorm:"foreignKey:OperatorUID;references:UID"`
}

// DB 数据库
var (
	db  *gorm.DB
	err error
)

func init() {
	dbName := config.Conf.DB.Name
	dbUsername := config.Conf.DB.Username
	dbPassword := config.Conf.DB.Password
	dbHost := config.Conf.DB.Host
	dbPort := config.Conf.DB.Port

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUsername, dbPassword,
		dbHost, dbPort, dbName)

	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		panic("fail to open db")
	}

	if err = db.AutoMigrate(&User{}); err != nil {
		panic(err.Error())
	}

	if err = db.AutoMigrate(&Role{}); err != nil {
		panic(err.Error())
	}

	if err = db.AutoMigrate(&Token{}); err != nil {
		panic(err.Error())
	}
	if err = db.AutoMigrate(&BFile{}); err != nil {
		panic(err.Error())
	}
	if err = db.AutoMigrate(&File{}); err != nil {
		panic(err.Error())
	}

	if err = db.AutoMigrate(&Task{}); err != nil {
		panic(err.Error())
	}
}

// FirstInitData 首次初始化数据
func FirstInitData() {
	err := db.Transaction(func(tx *gorm.DB) error {
		var user User
		result := tx.Where(&User{Username: "admin"}).First(&user)
		if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
			return result.Error
		}

		if result.Error == nil {
			return errors.New("database already init")
		}

		roles := []Role{{
			Name: "admin",
			Desc: "管理员",
			AuthInfo: AuthInfo{
				Auth: []Auth{
					{
						Flag:     auths.AuthAll,
						Resource: []string{"all"},
						Desc:     "所有权限",
					},
				},
			},
		}, {
			Name: "editor",
			Desc: "编辑者",
			AuthInfo: AuthInfo{
				Auth: []Auth{
					{
						Flag:     auths.AuthDownloadFile,
						Resource: []string{"all"},
						Desc:     "下载",
					},
					{
						Flag:     auths.AuthUploadFile,
						Resource: []string{"all"},
						Desc:     "上传",
					},
					{
						Flag:     auths.AuthViewFile,
						Resource: []string{"all"},
						Desc:     "查看",
					},
				},
			},
		}, {
			Name: "downloader",
			Desc: "下载者",
			AuthInfo: AuthInfo{
				Auth: []Auth{
					{
						Flag:     auths.AuthDownloadFile,
						Resource: []string{"all"},
						Desc:     "仅下载权限",
					},
				},
			},
		}}
		result = tx.Create(roles)
		if result.Error != nil {
			return result.Error
		}

		pwd := "admin" + config.Conf.Svr.AppKey
		hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		users := []User{
			{
				UID:      uuid.NewV1().String(),
				Username: "admin",
				Password: string(hash),
				RoleID:   roles[0].ID,
			},
			{
				UID:      uuid.NewV1().String(),
				Username: "editor",
				Password: string(hash),
				RoleID:   roles[1].ID,
			},
			{
				UID:      uuid.NewV1().String(),
				Username: "downloader",
				Password: string(hash),
				RoleID:   roles[2].ID,
			},
		}
		result = tx.Create(&users)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("database init succ.")
}
