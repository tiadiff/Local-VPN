package tunnel

// Command types
const (
	CmdConnect = 0x01
	CmdData    = 0x02
)

// A simplistic custom protocol packet for multiplexing
// [Type:1][ID:4][Length:4][Payload...]
// Type 0x01: Connect Request (Payload: target host string)
// Type 0x02: Data (Payload: data)
