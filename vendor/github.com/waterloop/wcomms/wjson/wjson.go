package wjson

import (
    "time"
    "encoding/json"

	wbin "wcomms/wbinary"
)

type CommPacketJson struct {
	Time int64     `json:"time"`
	Type string    `json:"type"`
	Id   uint8     `json:"name"`
	Data []float32 `json:"data"`
}

// Helper to get time in milliseconds
func currentTimeMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// Converts Binary Communication Packet to JSON string/bytes
func PacketEncodeJson(packet *wbin.CommPacket) ([]byte, error) {
	packetJson := &CommPacketJson{
		Time: currentTimeMs(),
		Type: wbin.TypeToString(packet.PacketType),
		Id:   packet.PacketId,
		Data: []float32{packet.Data1, packet.Data2, packet.Data3},
	}
	return json.Marshal(packetJson)
}

// Converts JSON string/bytes to Binary Communication Packet
func PacketDecodeJson(encoded []byte) (*wbin.CommPacket, error) {
	packetJson := &CommPacketJson{}
	err := json.Unmarshal(encoded, packetJson)
	packet := &wbin.CommPacket{
		PacketType: wbin.StringToType(packetJson.Type),
		PacketId:   packetJson.Id,
		Data1:      packetJson.Data[0],
		Data2:      packetJson.Data[1],
		Data3:      packetJson.Data[2],
	}
	return packet, err
}
