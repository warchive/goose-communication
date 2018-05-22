package wbinary

import (
	"math"
)

// Packet type
type PacketType uint8

// Packet type enums
const (
	Sensor  PacketType = iota // 00
	Command                   // 01
	State                     // 10
	Log                       // 11
)

func TypeToString(pType PacketType) string {
	switch pType {
	case Sensor:
		return "sensor"
	case Command:
		return "command"
	case State:
		return "state"
	case Log:
		return "log"
	default:
		return "unknown"
	}
}

func StringToType(pStr string) PacketType {
	switch pStr {
	case "sensor":
		return Sensor
	case "command":
		return Command
	case "state":
		return State
	case "log":
		return Log
	default:
		return Command
	}
}

// Packet specification:
// [type:3][name:7][data1:18][data2:18][data3:18]
var packetStructure = []uint{3, 7, 18, 18, 18}

// Communication packet type
type CommPacket struct {
	PacketType PacketType
	PacketId   uint8
	Data1      float32
	Data2      float32
	Data3      float32
}

// Reads no more than 32 bits at a time from
// a byte buffer between the start position i
// and end position j, inclusive
func readBits(buf []byte, i uint, j uint) uint32 {
	if i > j {
		j, i = i, j
	}
	iStart := i / 8
	jStart := j / 8
	iMod := i % 8
	jMod := j % 8
	if iStart >= uint(len(buf)) {
		return 0
	}
	if jStart >= uint(len(buf)) {
		jStart = uint(len(buf) - 1)
		jMod = 7
	}
	if iStart == jStart {
		return uint32((buf[iStart] >> iMod) & (0xff >> (7 - jMod + iMod)))
	}
	var res uint32 = 0
	ptr := 8 - iMod
	res |= uint32(buf[iStart] >> iMod)
	iStart++
	for ; iStart < jStart; iStart++ {
		res |= uint32(buf[iStart]) << ptr
		ptr += 8
	}
	res |= uint32(buf[jStart]&(0xff>>(7-jMod))) << ptr
	return res
}

func setBits(buf []byte, bits uint32, i uint, length uint) {
	bits &= 0xffffffff >> (0x20 - length)
	j := i + length - 1
	iStart := i / 8
	jStart := j / 8
	iMod := i % 8
	jMod := j % 8
	// Assume that the provided index and length are valid
	if iStart == jStart {
		val := byte(bits)
		buf[iStart] &= ^((0xff >> (7 - jMod)) & (0xff << iMod))
		buf[iStart] |= val << iMod
		return
	}
	slice := byte(bits)
	buf[iStart] &= ^(0xff >> iMod << iMod)
	buf[iStart] |= slice << iMod
	bitPtr := 8 - iMod
	iStart++
	for ; iStart < jStart; iStart++ {
		slice = byte(bits >> bitPtr)
		buf[iStart] = slice
		bitPtr += 8
	}
	buf[jStart] &= 0xff >> jMod << jMod
	buf[jStart] |= byte(bits >> bitPtr)
}

// Read segments of bits, no larger than
// 32 bits as defined by shape
func readSegments(buf []byte, shape []uint) []uint32 {
	shapeLen := uint(len(shape))
	bufLen := uint(len(buf))
	var i uint = 0
	var k uint = 0
	result := make([]uint32, shapeLen, shapeLen)
	for i < bufLen*8 && k < shapeLen {
		result[k] = readBits(buf, i, i+shape[k]-1)
		i += shape[k]
		k++
	}
	return result
}

func writeSegments(buf []byte, shape []uint, values []uint32) {
	shapeLen := uint(len(shape))
	bufLen := uint(len(buf))
	var k uint = 0
	var i uint = 0
	for i < bufLen*8 && k < shapeLen {
		setBits(buf, values[k], i, shape[k])
		i += shape[k]
		k++
	}
}

// Floating point specifications:
// [decimalPart:7][integerPart:10][sign:1]
func decodeFloat18(bits uint32) float32 {
	if bits == 0 {
		// Bit of a hack in place of a better compression algorithm
		// I'm too lazy to port https://goo.gl/ftEgQ1
		return float32(0)
	}
	// Expand bits to float32 LSB
	expanded := (bits & 0xfff) << 11
	expanded |= (((bits & 0x1f000) >> 12) + 0x70) << 23
	expanded |= (bits & 0x20000) << 14
	return math.Float32frombits(expanded)
}

func encodeFloat18(val float32) uint32 {
	// Compress bits
	bits := math.Float32bits(val)
	compressed := (bits & 0x7ff800) >> 11
	compressed |= (((bits & 0x7f800000) >> 23) - 0x70) << 12
	compressed |= (bits & 0x80000000) >> 14
	return compressed
}

// Read the packet data as defined by the packet
// structure. Function will try to read as much
// data as possible.
func ReadPacket(buf []byte) *CommPacket {
	var bitSegments = readSegments(buf, packetStructure)
	return &CommPacket{
		PacketType: PacketType(bitSegments[0]),
		PacketId:   uint8(bitSegments[1]),
		Data1:      decodeFloat18(bitSegments[2]),
		Data2:      decodeFloat18(bitSegments[3]),
		Data3:      decodeFloat18(bitSegments[4]),
	}
}

func WritePacket(packet *CommPacket) []byte {
	var bitSegments = make([]uint32, 5, 5)
	bitSegments[0] = uint32(packet.PacketType)
	bitSegments[1] = uint32(packet.PacketId)
	bitSegments[2] = encodeFloat18(packet.Data1)
	bitSegments[3] = encodeFloat18(packet.Data2)
	bitSegments[4] = encodeFloat18(packet.Data3)
	var buf = make([]byte, 8, 8)
	writeSegments(buf, packetStructure, bitSegments)
	return buf
}
