package main

import (
	"fmt"
	"github.com/toophy/chat_client/proto"
	"github.com/toophy/toogo"
	"sync"
)

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

// 主线程
type MasterThread struct {
	toogo.Thread
	AccountsMutex sync.RWMutex            // 帐号读写锁
	AccountsId    map[uint64]*GameAccount // 帐号ID索引
	AccountsName  map[string]*GameAccount // 帐号Name索引
	RolesMutex    sync.RWMutex            // 帐号读写锁
	RolesId       map[uint64]*GameRole    // 角色ID索引
	RolesName     map[string]*GameRole    // 角色Name索引
}

// 首次运行
func (this *MasterThread) On_firstRun() {
}

// 响应线程最先运行
func (this *MasterThread) On_preRun() {
	// 处理各种最先处理的问题
}

// 响应线程运行
func (this *MasterThread) On_run() {
}

// 响应线程退出
func (this *MasterThread) On_end() {
}

// 响应网络事件
func (this *MasterThread) On_netEvent(m *toogo.Tmsg_net) bool {

	name_fix := m.Name
	if len(name_fix) == 0 {
		name_fix = fmt.Sprintf("Conn[%d]", m.SessionId)
	}

	switch m.Msg {
	case "listen failed":
		this.LogFatal("%s : Listen failed[%s]", name_fix, m.Info)

	case "listen ok":
		this.LogInfo("%s : Listen(%s) ok.", name_fix, toogo.GetSessionById(m.SessionId).GetIPAddress())

	case "accept failed":
		this.LogFatal(m.Info)
		return false

	case "accept ok":
		this.LogDebug("%s : Accept ok", name_fix)

	case "connect failed":
		this.LogError("%s : Connect failed[%s]", name_fix, m.Info)

	case "connect ok":
		this.LogDebug("%s : Connect ok", name_fix)

		p := toogo.NewPacket(128, m.SessionId)

		msgLogin := new(proto.S2G_registe)
		msgLogin.Sid = 1
		msgLogin.Write(p)

		toogo.SendPacket(p)

	case "read failed":
		this.LogError("%s : Connect read[%s]", name_fix, m.Info)

	case "pre close":
		this.LogDebug("%s : Connect pre close", name_fix)

	case "close failed":
		this.LogError("%s : Connect close failed[%s]", name_fix, m.Info)

	case "close ok":
		this.LogDebug("%s : Connect close ok.", name_fix)
	}

	return true
}

// -- 当网络消息包解析出现问题, 如何处理?
func (this *MasterThread) On_packetError(sessionId uint64) {
	toogo.CloseSession(this.Get_thread_id(), sessionId)
}

// 注册消息
func (this *MasterThread) On_registNetMsg() {
	this.RegistNetMsg(proto.C2S_chat_Id, this.on_c2s_chat)
}

func (this *MasterThread) on_c2s_chat(pack *toogo.PacketReader, sessionId uint64) bool {
	msg := proto.C2S_chat{}
	msg.Read(pack)

	this.LogInfo("Say : %s", msg.Data)

	return true
}

func (this *MasterThread) AddAccount(a *GameAccount) bool {
	// a.Id范围检查
	// a.Name长度检查[1,64]

	this.AccountsMutex.Lock()
	defer this.AccountsMutex.Unlock()

	if _, ok := this.AccountsId[a.Id]; ok {
		return false
	}
	if _, ok := this.AccountsName[a.Name]; ok {
		return false
	}

	this.AccountsId[a.Id] = a
	this.AccountsName[a.Name] = a
	return true
}

func (this *MasterThread) AddRole(p *GameRole, a *GameAccount) bool {

	this.AccountsMutex.Lock()
	defer this.AccountsMutex.Unlock()
	this.RolesMutex.Lock()
	defer this.RolesMutex.Unlock()

	if _, ok := a.RolesId[p.Id]; ok {
		return false
	}
	if _, ok := a.RolesName[p.Name]; ok {
		return false
	}
	if _, ok := this.RolesId[p.Id]; ok {
		return false
	}
	if _, ok := this.RolesName[p.Name]; ok {
		return false
	}

	this.RolesId[p.Id] = p
	this.RolesName[p.Name] = p
	this.RolesId[p.Id] = p
	this.RolesId[p.Id] = p

	return true
}

func (this *MasterThread) GetRoleById(id uint64) *GameRole {
	this.RolesMutex.RLock()
	defer this.RolesMutex.RUnlock()
	if v, ok := this.RolesId[id]; ok {
		return v
	}
	return nil
}

func (this *MasterThread) GetRoleByName(name string) *GameRole {
	this.RolesMutex.RLock()
	defer this.RolesMutex.RUnlock()
	if v, ok := this.RolesName[name]; ok {
		return v
	}
	return nil
}

func (this *MasterThread) GetAccountById(id uint64) *GameAccount {
	this.AccountsMutex.RLock()
	defer this.AccountsMutex.RUnlock()
	if v, ok := this.AccountsId[id]; ok {
		return v
	}
	return nil
}

func (this *MasterThread) GetAccountByName(name string) *GameAccount {
	this.AccountsMutex.RLock()
	defer this.AccountsMutex.RUnlock()
	if v, ok := this.AccountsName[name]; ok {
		return v
	}
	return nil
}

func main() {
	main_thread := new(MasterThread)
	main_thread.Init_thread(main_thread, toogo.Tid_master, "master", 1000, 100, 10000)
	toogo.Run(main_thread)
}
