package mc_server

import (
	"Eulogist/core/fb_auth/py_rpc"
	"Eulogist/core/minecraft/protocol/packet"
	RaknetConnection "Eulogist/core/raknet"
	"encoding/json"
	"fmt"
)

// OnPyRpc 处理数据包 PyRpc。
//
// 如果必要，将使用 writePacketToClient
// 向 Minecraft 客户端发送新数据包
func (m *MinecraftServer) OnPyRpc(
	p *packet.PyRpc,
	writePacketToClient func(pk RaknetConnection.MinecraftPacket, useBytes bool) error,
) (shouldSendCopy bool, err error) {
	// 解码 PyRpc
	if p.Value == nil {
		return true, nil
	}
	content, err := py_rpc.Unmarshal(p.Value)
	if err != nil {
		return true, fmt.Errorf("OnPyRpc: %v", err)
	}
	// 根据内容类型处理不同的 PyRpc
	switch c := content.(type) {
	case *py_rpc.StartType:
		c.Content = m.fbClient.TransferData(c.Content)
		c.Type = py_rpc.StartTypeResponse
		err = m.WritePacket(
			RaknetConnection.MinecraftPacket{
				Packet: &packet.PyRpc{
					Value:         py_rpc.Marshal(c),
					OperationType: packet.PyRpcOperationTypeSend,
				},
			}, false,
		)
		if err != nil {
			return false, fmt.Errorf("OnPyRpc: %v", err)
		}
	case *py_rpc.GetMCPCheckNum:
		// 如果已完成零知识证明(挑战)，
		// 则不做任何操作
		if m.getCheckNumEverPassed {
			break
		}
		// 创建请求并发送到认证服务器并获取响应
		arg, _ := json.Marshal([]any{
			c.FirstArg,
			c.SecondArg.Arg,
			m.entityUniqueID,
		})
		ret := m.fbClient.TransferCheckNum(string(arg))
		// 解码响应并调整数据
		ret_p := []any{}
		json.Unmarshal([]byte(ret), &ret_p)
		if len(ret_p) > 7 {
			ret6, ok := ret_p[6].(float64)
			if ok {
				ret_p[6] = int64(ret6)
			}
		}
		// 完成零知识证明(挑战)
		err = m.WritePacket(
			RaknetConnection.MinecraftPacket{
				Packet: &packet.PyRpc{
					Value:         py_rpc.Marshal(&py_rpc.SetMCPCheckNum{ret_p}),
					OperationType: packet.PyRpcOperationTypeSend,
				},
			}, false,
		)
		if err != nil {
			return false, fmt.Errorf("OnPyRpc: %v", err)
		}
		// 标记零知识证明(挑战)已在当前会话下永久完成
		m.getCheckNumEverPassed = true
		// 返回值
		return false, nil
	default:
		// 对于其他种类的 PyRpc 数据包，
		// 返回 true 表示需要将数据包抄送至
		// Minecraft 客户端
		return true, nil
	}
	// 返回值
	return false, nil
}
