package account

import ()

// 帐号
type GameAccount struct {
	Id        uint64               // 唯一ID
	Name      string               // 帐号名
	RolesId   map[uint64]*GameRole // 角色ID索引
	RolesName map[string]*GameRole // 角色Name索引
}

// 角色
type GameRole struct {
	Id          uint64 // 唯一ID
	AccountName uint64 // 帐号名
	Name        string // 角色名, 可以为空
	SessionId   uint64 // 网关中的会话ID
	IpAddress   string // 当前登录IP地址
}

func (this *GameRole) GetAccountId() uint64 {
	return this.Id & 0x3FFFFFFFFFF
}
