// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT. Created at 2022-01-17 21:12:32
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// User is the golang structure for table user.
type User struct {
	Id        uint        `json:"id"        description:"UID"`
	Passport  string      `json:"passport"  description:"账号"`
	Password  string      `json:"password"  description:"MD5密码"`
	Nickname  string      `json:"nickname"  description:"昵称"`
	Avatar    string      `json:"avatar"    description:"头像地址"`
	Status    int         `json:"status"    description:"状态 0:启用 1:禁用"`
	Gender    int         `json:"gender"    description:"性别 0: 未设置 1: 男 2: 女"`
	CreatedAt *gtime.Time `json:"createdAt" description:"注册时间"`
	UpdatedAt *gtime.Time `json:"updatedAt" description:"更新时间"`
}
