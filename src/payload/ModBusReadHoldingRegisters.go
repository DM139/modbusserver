package payload

import (
	"bytes"
	//"container/list"
	"encoding/binary"
	"errors"
	"fmt"
)

type ReadHoldingRegistersHandler struct {
	*ModBusPayload
	Properties []string
	//datamap map[string]float32
	BytesdataType string
	//payloaduff []byte
	//valueList *list.List
}

func (rhr *ReadHoldingRegistersHandler) Encode() ([]byte, error){
	result := make([]byte, 0)
	buffer := bytes.NewBuffer(result)


	if rhr.Quantity < 1 || rhr.Quantity > 125 {
		s := fmt.Sprintf("modbus: quantity '%v' must be between '%v' and '%v',", rhr.Quantity, 1, 125)
		return nil, errors.New(s)
	}

	//if rhr.isCRC {
	//	if err := binary.Write(buffer, binary.BigEndian, rhr.unitIdentifier); err != nil {
	//		s := fmt.Sprintf("Pack unitIdentifier error , %v", err)
	//		return nil, errors.New(s)
	//	}
	//}

	if err := binary.Write(buffer, binary.BigEndian, rhr.FunctionCode); err != nil {
		s := fmt.Sprintf("Pack functionCode error , %v", err)
		return nil, errors.New(s)
	}

	if err := binary.Write(buffer, binary.BigEndian, rhr.StartAddress); err != nil {
		s := fmt.Sprintf("Pack startAddress error , %v", err)
		return nil, errors.New(s)
	}

	if err := binary.Write(buffer, binary.BigEndian, rhr.Quantity); err != nil {
		s := fmt.Sprintf("Pack quantity error , %v", err)
		return nil, errors.New(s)
	}

	//if rhr.isCRC {
	//	var crc crc
	//	crc.reset().pushBytes(buffer.Bytes())
	//	checksum := crc.value()
	//
	//	if err := binary.Write(buffer, binary.BigEndian, byte(checksum)); err != nil {
	//		s := fmt.Sprintf("Pack quantity error , %v", err)
	//		return nil, errors.New(s)
	//	}
	//
	//	if err := binary.Write(buffer, binary.BigEndian, byte(checksum >> 8)); err != nil {
	//		s := fmt.Sprintf("Pack quantity error , %v", err)
	//		return nil, errors.New(s)
	//	}
	//}


	return buffer.Bytes(), nil

}

func (rhr *ReadHoldingRegistersHandler) Decode(buf []byte)(error) {
	byteBuffer := bytes.NewBuffer(buf)
	var functionCode uint8
	var dataLength uint8

	_ = binary.Read(byteBuffer, binary.BigEndian, &functionCode)
	if functionCode != rhr.FunctionCode{
		s := fmt.Sprintf("modbus: payload functionCode is incorrect")
		return errors.New(s)
	}
	_ = binary.Read(byteBuffer, binary.BigEndian, &dataLength)
	dataLen := byteBuffer.Len()
	if dataLen != int(dataLength){
		s := fmt.Sprintf("modbus: payload dataLength is incorrect")
		return errors.New(s)
	}
	payloaduff := make([]byte, dataLength)
	_ = binary.Read(byteBuffer, binary.BigEndian, &payloaduff)

	m, err := rhr.ConvertBytesToMap(rhr.Properties, rhr.BytesdataType, payloaduff)
	for k, v := range m {
		fmt.Printf("%s,%d\n", k, v)
	}
	return err
}