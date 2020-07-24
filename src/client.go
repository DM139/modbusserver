package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"modbusserver/src/payload"
	protocols "modbusserver/src/protocols"

	"io"
	"log"
	"net"
)

// Example command: go run client.go
func main() {
	functionCode := uint8(0x03)
	startAddress := uint16(0x0001)
	quantity := uint16(0x0008)

	properties := []string{"v1", "v2", "v3", "v4",
		"v5", "v6", "v7", "v8",
	}
	bytesdataType := "uint16"

	mbp := &payload.ModBusPayload{functionCode, startAddress, quantity}
	rhr := &payload.ReadHoldingRegistersHandler{ModBusPayload: mbp, Properties: properties, BytesdataType: bytesdataType}
	var ph payload.PayloadHandler = rhr


	conn, err := net.Dial("tcp", "127.0.0.1:5020")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	go func() {
		for {

			response, err := ClientDecode(conn)
			if err != nil {
				log.Printf("ClientDecode error, %v\n", err)
			}

			log.Printf("receive , %v, data:%v\n", response, []byte(response.Data))
			err1 := ph.Decode(response.Data);
			if err1 != nil {
				log.Printf("ClientDecode error, %v\n", err1)
			}
		}
	}()

	buf, err := ph.Encode();
	pbdata, err := ClientEncode(protocols.DefaultTransactionIdentifier, protocols.DefaultProtocolIdentifier, protocols.DefaultUnitIdentifier, buf)
	if err != nil {
		panic(err)
	}
	conn.Write(pbdata)

	//data = []byte("world")
	//pbdata, err = ClientEncode(protocols.DefaultTransactionIdentifier, protocols.DefaultProtocolIdentifier, protocols.DefaultUnitIdentifier, data)
	//if err != nil {
	//	panic(err)
	//}
	//conn.Write(pbdata)

	select {}
}

// ClientEncode :
func ClientEncode(transactionIdentifier, protocolIdentifier uint16, unitIdentifier uint8, data []byte) ([]byte, error) {
	result := make([]byte, 0)

	buffer := bytes.NewBuffer(result)

	if err := binary.Write(buffer, binary.BigEndian, transactionIdentifier); err != nil {
		s := fmt.Sprintf("Pack transactionIdentifier error , %v", err)
		return nil, errors.New(s)
	}

	if err := binary.Write(buffer, binary.BigEndian, protocolIdentifier); err != nil {
		s := fmt.Sprintf("Pack protocolIdentifier error , %v", err)
		return nil, errors.New(s)
	}
	dataLen := uint16(1 + len(data))
	if err := binary.Write(buffer, binary.BigEndian, dataLen); err != nil {
		s := fmt.Sprintf("Pack datalength error , %v", err)
		return nil, errors.New(s)
	}

	if err := binary.Write(buffer, binary.BigEndian, unitIdentifier); err != nil {
		s := fmt.Sprintf("Pack UnitIdentifier error , %v", err)
		return nil, errors.New(s)
	}

	if dataLen > 1 {
		if err := binary.Write(buffer, binary.BigEndian, data); err != nil {
			s := fmt.Sprintf("Pack data error , %v", err)
			return nil, errors.New(s)
		}
	}

	return buffer.Bytes(), nil
}

// ClientDecode :
func ClientDecode(rawConn net.Conn) (*protocols.TcpModbusProtocol, error) {
	newPackage := protocols.TcpModbusProtocol{}
	//log.Printf("data:%s\n",&newPackage.DataLength)

	headData := make([]byte, protocols.DefaultHeadLength)
	n, err := io.ReadFull(rawConn, headData)
	if n != protocols.DefaultHeadLength {
		return nil, err
	}

	// parse protocols header
	bytesBuffer := bytes.NewBuffer(headData)
	var TransactionIdentifier uint16
	binary.Read(bytesBuffer, binary.BigEndian, &TransactionIdentifier)
	binary.Read(bytesBuffer, binary.BigEndian, &newPackage.ProtocolIdentifier)
	binary.Read(bytesBuffer, binary.BigEndian, &newPackage.DataLength)
	binary.Read(bytesBuffer, binary.BigEndian, &newPackage.UnitIdentifier)
	newPackage.TransactionIdentifier = uint32(TransactionIdentifier)
	if newPackage.DataLength < 2 {
		return &newPackage, nil
	}

	data := make([]byte, newPackage.DataLength -1)
	dataNum, err2 := io.ReadFull(rawConn, data)
	if uint16(dataNum) != (newPackage.DataLength - 1) {
		s := fmt.Sprintf("read data error, %v", err2)
		return nil, errors.New(s)
	}

	newPackage.Data = data

	return &newPackage, nil
}