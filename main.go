package main

import (
	"fmt"
	"github.com/toophy/chat_client/proto"
	"github.com/toophy/toogo"
)

// 主线程
type MasterThread struct {
	toogo.Thread
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
		name_fix = fmt.Sprintf("Conn[%d]", m.Id)
	}

	switch m.Msg {
	case "listen failed":
		this.LogFatal("%s : Listen failed[%s]", name_fix, m.Info)

	case "listen ok":
		this.LogInfo("%s : Listen(0.0.0.0:%d) ok.", name_fix, 8001)

	case "accept failed":
		this.LogFatal(m.Info)
		return false

	case "accept ok":
		this.LogDebug("%s : Accept ok", name_fix)

	case "connect failed":
		this.LogError("%s : Connect failed[%s]", name_fix, m.Info)

	case "connect ok":
		this.LogDebug("%s : Connect ok", name_fix)

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
func (this *MasterThread) On_packetError(m *toogo.Tmsg_packet) {
	toogo.CloseSession(this.Get_thread_id(), m.SessionId)
}

// 注册消息
func (this *MasterThread) On_registNetMsg() {
	this.RegistNetMsg(proto.C2M_login_Id, this.on_c2m_login)
}

func (this *MasterThread) on_c2m_login(pack *toogo.PacketReader, sessionId uint32) bool {
	msg := new(proto.C2M_login)
	msg.Read(pack)

	p := new(toogo.PacketWriter)
	d := make([]byte, 64)
	p.InitWriter(d)
	msgLoginRet := new(proto.M2C_login_ret)
	msgLoginRet.Ret = 0
	msgLoginRet.Msg = "ok"
	msgLoginRet.Write(p)

	p.PacketWriteOver()
	session := toogo.GetSessionById(sessionId)
	m := new(toogo.Tmsg_packet)
	m.Data = p.GetData()
	m.Len = uint32(p.GetPos())
	m.Count = uint32(p.Count)

	toogo.PostThreadMsg(session.MailId, m)
	return true
}

func main() {
	main_thread := new(MasterThread)
	main_thread.Init_thread(main_thread, toogo.Tid_master, "master", 1000, 100, 10000)
	toogo.Run(main_thread)
}
