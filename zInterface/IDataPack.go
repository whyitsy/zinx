package zInterface

type IDataPack interface {
	// GetHeadLen 获取包头长度方法
	GetHeadLen() uint32

	// Pack 封包方法(封装数据)
	Pack(msg IMessage) ([]byte, error)

	// UnPack 拆包方法(拆开数据)
	UnPack([]byte) (IMessage, error)
}
