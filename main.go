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
func (this *MasterThread) On_first_run() {
}

// 响应线程最先运行
func (this *MasterThread) On_pre_run() {
	// 处理各种最先处理的问题
}

// 响应线程运行
func (this *MasterThread) On_run() {
}

// 响应线程退出
func (this *MasterThread) On_end() {
}

// 响应网络事件
func (this *MasterThread) On_NetEvent(m *toogo.Tmsg_net) bool {

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

// 响应网络消息包
func (this *MasterThread) On_NetPacket(m *toogo.Tmsg_packet) bool {
	p := new(toogo.PacketReader)
	p.InitReader(m.Data, uint16(m.Count))

	fmt.Println(m.Count)

	for i := uint32(0); i < m.Count; i++ {
		msg_len := p.ReadUint16()
		msg_id := p.ReadUint16()
		fmt.Println(msg_len, msg_id)
		switch msg_id {
		case proto.C2M_login_Id:
			msg := new(proto.C2M_login)
			msg.Read(p)
			fmt.Printf("%d,%d,%-v\n", msg_len, msg_id, msg)
		}
	}

	return true
}

func main() {
	main_thread := new(MasterThread)
	main_thread.Init_thread(main_thread, toogo.Tid_master, "master", 100, 10000)
	toogo.Run(main_thread)
}
