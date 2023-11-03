package response

import (
	"fmt"
	"time"
)

type TimeJson time.Time

func (t TimeJson) MarshalJSON() ([]byte, error) {
	stmp := fmt.Sprintf("\"%s\"", time.Time(t).Format("2006-01-01"))
	return []byte(stmp), nil
}

type UserResponse struct {
	Id       int32    `json:"id,omitempty"`
	Mobile   string   `json:"mobile,omitempty"`
	Password string   `json:"password,omitempty"`
	NickName string   `json:"nickName,omitempty"`
	Birthday TimeJson `json:"birthday,omitempty"`
	Gender   string   `json:"gender,omitempty"`
	Role     int      `json:"role,omitempty"`
}
