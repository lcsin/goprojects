package domain

import "time"

type User struct {
	UID      int64  `json:"uid"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Passwd   string `json:"passwd"`
	Profile  string `json:"profile"`

	Birthday   time.Time `json:"birthday"`
	CreateTime time.Time `json:"createTime"`
}
