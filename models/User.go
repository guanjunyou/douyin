package models

type User struct {
	CommonEntity
	//Id            int64  `json:"id,omitempty"`
	Name          string `json:"name,omitempty"`
	UserName      string `json:"username,omitempty"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	Phone         string `json:"phone,omitempty"`
	Password      string `json:"password_count,omitempty"`
	Icon          string `json:"icon,omitempty"`
	Gender        int    `json:"gender,omitempty"`
	Age           int    `json:"age,omitempty"`
	NickName      string `json:"nickname,omitempty"`
}

func (table *User) TableName() string {
	return "user"
}

func GetUserById(Id int64) (User, error) {
	var user User
	DB.Where("id = ?", Id).Find(&user)
	return user, nil
}
