package models

import "github.com/RaymondCode/simple-demo/utils"

type User struct {
	CommonEntity
	//Id            int64  `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	FollowCount    int64  `json:"follow_count,omitempty"`
	FollowerCount  int64  `json:"follower_count,omitempty"`
	Phone          string `json:"phone,omitempty"`
	Password       string `json:"password,omitempty"`
	Icon           string `json:"icon,omitempty"`
	Gender         int    `json:"gender,omitempty"`
	Age            int    `json:"age,omitempty"`
	Nickname       string `json:"nickname,omitempty"`
	Siganture      string `json:"signature,omitempty"`
	TotalFavorited string `json:"total_favorited,omitempty"`
	WorkCount      string `json:"work_count,omitempty"`
	FavoriteCount  string `json:"favorite_count,omitempty"`
}

func (table *User) TableName() string {
	return "user"
}

func GetUserById(Id int64) (User, error) {
	var user User
	// 传参禁止直接字符串拼接，防止SQL注入
	err := utils.DB.Where("id = ? AND is_deleted != ?", Id, 1).First(&user).Error
	if err != nil {
		return user, err
	}
	return user, nil
}

func SaveUser(user User) error {
	return utils.DB.Create(&user).Error
}
