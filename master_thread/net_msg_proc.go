package master_thread

import (
	"github.com/toophy/chat_svr/proto"
	"github.com/toophy/toogo"
)

func (this *MasterThread) on_c2s_chat(pack *toogo.PacketReader, sessionId uint64) bool {
	msg := proto.C2S_chat{}
	msg.Read(pack)

	this.LogInfo("Say : %s", msg.Data)

	// 广播消息
	p := toogo.NewPacket(128, sessionId)

	msgChat := new(proto.S2C_chat)
	msgChat.Data = msg.Data
	msgChat.Write(p, pack.LinkTgid)

	p.PacketWriteOver()
	toogo.SendPacket(p)

	return true
}
