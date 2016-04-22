package master_thread

import (
	"github.com/toophy/chat_svr/proto"
)

// 注册消息
func (this *MasterThread) On_registNetMsg() {
	this.RegistNetMsg(proto.C2S_chat_Id, this.on_c2s_chat)
}
