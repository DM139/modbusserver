package protocols

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"sync/atomic"
	"github.com/panjf2000/gnet"
)

// CustomLengthFieldProtocol : custom protocols
// custom protocols header contains Version, ActionType and DataLength fields
// its payload is Data field
type TcpModbusProtocol struct {
	TransactionIdentifier uint32
	ProtocolIdentifier uint16
	DataLength uint16
	UnitIdentifier uint8
	//FunctionCode uint8
	Data       []byte
}


// Encode ...
func (cc *TcpModbusProtocol) Encode(c gnet.Conn, buf []byte) ([]byte, error) {
	result := make([]byte, 0)

	buffer := bytes.NewBuffer(result)

	// take out the param
	item := c.Context().(TcpModbusProtocol)

	transactionId := atomic.AddUint32(&item.TransactionIdentifier, 1)
	item.TransactionIdentifier = transactionId
	if err := binary.Write(buffer, binary.BigEndian, uint16(transactionId)); err != nil {
		s := fmt.Sprintf("Pack TransactionIdentifier error , %v", err)
		return nil, errors.New(s)
	}

	if err := binary.Write(buffer, binary.BigEndian, item.ProtocolIdentifier); err != nil {
		s := fmt.Sprintf("Pack ProtocolIdentifier error , %v", err)
		return nil, errors.New(s)
	}

	dataLen := uint16(1 + len(buf))
	if err := binary.Write(buffer, binary.BigEndian, dataLen); err != nil {
		s := fmt.Sprintf("Pack datalength error , %v", err)
		return nil, errors.New(s)
	}

	if err := binary.Write(buffer, binary.BigEndian, item.UnitIdentifier); err != nil {
		s := fmt.Sprintf("Pack UnitIdentifier error , %v", err)
		return nil, errors.New(s)
	}

	if dataLen > 1 {
		if err := binary.Write(buffer, binary.BigEndian, buf); err != nil {
			s := fmt.Sprintf("Pack data error , %v", err)
			return nil, errors.New(s)
		}
	}

	return buffer.Bytes(), nil
}

// Decode ...
func (cc *TcpModbusProtocol) Decode(c gnet.Conn) ([]byte, error) {
	// parse header
	headerLen := DefaultHeadLength // uint16+uint16+uint16+uint8
	if size, header := c.ReadN(headerLen); size == headerLen {
		byteBuffer := bytes.NewBuffer(header)
		var TransactionIdentifier, ProtocolIdentifier, DataLength uint16
		var UnitIdentifier uint8
		_ = binary.Read(byteBuffer, binary.BigEndian, &TransactionIdentifier)
		_ = binary.Read(byteBuffer, binary.BigEndian, &ProtocolIdentifier)
		_ = binary.Read(byteBuffer, binary.BigEndian, &DataLength)
		_ = binary.Read(byteBuffer, binary.BigEndian, &UnitIdentifier)
		// to check the protocols version and actionType,
		// reset buffer if the version or actionType is not correct
		item := c.Context().(TcpModbusProtocol)
		if TransactionIdentifier != item.ProtocolIdentifier || ProtocolIdentifier != item.ProtocolIdentifier {
			c.ResetBuffer()
			log.Println("not normal protocols:", TransactionIdentifier, ProtocolIdentifier, DataLength, UnitIdentifier)
			return nil, errors.New("not normal protocols")
		}
		// parse payload
		dataLen := int(DataLength) //max int32 can contain 210MB payload
		protocolLen := headerLen + dataLen -1
		if dataSize, data := c.ReadN(protocolLen); dataSize == protocolLen {
			c.ShiftN(protocolLen)
			// log.Println("parse success:", data, dataSize)

			// return the payload of the data
			return data[headerLen:], nil
		}
		// log.Println("not enough payload data:", dataLen, protocolLen, dataSize)
		return nil, errors.New("not enough payload data")

	}
	// log.Println("not enough header data:", size)
	return nil, errors.New("not enough header data")
}

// default custom protocols const
const (
	DefaultHeadLength = 7

	//DefaultProtocolVersion = 0x8001 // test protocols version

	ActionPing = 0x0001 // ping
	ActionPong = 0x0002 // pong
	ActionData = 0x00F0 // business

	DefaultTransactionIdentifier = 0x0000
	DefaultProtocolIdentifier = 0x0000
	DefaultUnitIdentifier = 0x01
	//DefaultUnitIdentifier = 0x01
	//DefaultFunctionCode = 0x03

)

func isCorrectAction(actionType uint16) bool {
	switch actionType {
	case ActionPing, ActionPong, ActionData:
		return true
	default:
		return false
	}
}

