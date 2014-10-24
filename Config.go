package gnet

type Config struct {
	Addr			string		`json:"listen_addr"`
	MaxClients		uint32		`json:"max_clients"`
	MaxBodyLen		uint32		`json:"max_body_len"`
	HBTimeout		uint32		`json:"heartbeat_timeout"`
	ReadTimeout		uint32		`json:"read_timeout"`
	WriteTimeout	uint32		`json:"write_timeout"`
}

var (
	Conf	Config
)