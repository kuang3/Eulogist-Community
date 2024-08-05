package mc_client

import (
	RaknetConnection "Eulogist/core/raknet"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"net"

	"github.com/sandertv/go-raknet"
)

// CreateListener 在 127.0.0.1 上以 Raknet 协议侦听 Minecraft 客户端的连接，
// 这意味着您成功创建了一个 Minecraft 数据包代理服务器。
// 稍后，您可以通过 m.GetServerAddress 来取得服务器地址
func (m *MinecraftClient) CreateListener() error {
	// 创建一个 Raknet 监听器
	listener, err := raknet.Listen("127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("CreateListener: %v", err)
	}
	// 获取监听器的地址
	address, ok := listener.Addr().(*net.UDPAddr)
	if !ok {
		return fmt.Errorf("CreateListener: Failed to get address for listener")
	}
	// 初始化变量
	m.listener = listener
	m.address = address
	m.connected = make(chan struct{}, 1)
	m.Raknet = RaknetConnection.NewRaknet()
	// 返回成功
	return nil
}

// WaitConnect 等待 Minecraft 客户端连接到服务器
func (m *MinecraftClient) WaitConnect() error {
	// 接受客户端连接
	conn, err := m.listener.Accept()
	if err != nil {
		return fmt.Errorf("WaitConnect: %v", err)
	}
	// 初始化变量
	serverKey, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	m.SetConnection(conn, serverKey)
	m.connected <- struct{}{}
	// 返回成功
	return nil
}

// GetServerIP 获取服务器的 IP 地址
func (m *MinecraftClient) GetServerIP() string {
	return m.address.IP.String()
}

// GetServerPort 获取服务器的端口号
func (m *MinecraftClient) GetServerPort() int {
	return m.address.Port
}

// ...
func (m *MinecraftClient) InitPlayerSkin() {
	m.playerSkin = &RaknetConnection.Skin{}
}

// ...
func (m *MinecraftClient) SetPlayerSkin(skin *RaknetConnection.Skin) {
	m.playerSkin = skin
}

// ...
func (m *MinecraftClient) GetPlayerSkin() *RaknetConnection.Skin {
	return m.playerSkin
}

// ...
func (m *MinecraftClient) GetEntityUniqueID() int64 {
	return m.entityUniqueID
}

// ...
func (m *MinecraftClient) SetEntityUniqueID(entityUniqueID int64) {
	m.entityUniqueID = entityUniqueID
}
