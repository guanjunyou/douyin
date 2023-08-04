package models

import (
	"github.com/RaymondCode/simple-demo/utils"
	"gorm.io/gorm"
)

type User struct {
	utils.CommonEntity
	//Id            int64  `json:"id,omitempty"`
	Name            string `json:"name"`
	FollowCount     int64  `json:"follow_count"`
	FollowerCount   int64  `json:"follower_count"`
	Phone           string `json:"phone"`
	Password        string `json:"password"`
	Avatar          string `json:"avatar"`
	Gender          int    `json:"gender"`
	Age             int    `json:"age"`
	Nickname        string `json:"nickname"`
	Signature       string `json:"signature"`
	TotalFavorited  int64  `json:"total_favorited"`
	WorkCount       int64  `json:"work_count"`
	FavoriteCount   int64  `json:"favorite_count"`
	IsFollow        bool   `json:"is_follow"`
	BackgroundImage string `json:"background_image"`
}

func (table *User) TableName() string {
	return "user"
}

func GetUserById(Id int64) (User, error) {
	var user User
	// 传参禁止直接字符串拼接，防止SQL注入
	err := utils.GetMysqlDB().Where("id = ? AND is_deleted != ?", Id, 1).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func GetUserByName(name string) (User, error) {
	var user User
	// 传参禁止直接字符串拼接，防止SQL注入
	err := utils.GetMysqlDB().Where("name = ? AND is_deleted != ?", name, 1).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func SaveUser(user User) error {
	return utils.GetMysqlDB().Create(&user).Error
}

func UpdateUser(tx *gorm.DB, user User) error {
	err := tx.Save(&user).Error
	return err
}
