package payload

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

//type ProtocolDataUnit struct {
//	Data         []byte
//}


type (
	PayloadHandler interface {
		Encode()([]byte, error)
		Decode(buf []byte)(error)
    }
    ModBusPayload struct {
		FunctionCode uint8
		StartAddress uint16
		Quantity uint16

		//dataLength uint8
		////data []byte
		//isCRC bool
		//unitIdentifier uint8
		//CRCData	uint16
	}
)

func (mp *ModBusPayload) Encode()([]byte, error){
	return nil, nil
}
func (mp *ModBusPayload) Decode(buf []byte)(error){
	return nil
}

func (mp *ModBusPayload) ConvertBytesToMap(properties []string, bytesdataType string, buff []byte)  (map[string]float32, error){
	result := make(map[string]float32)
	byteBuffer := bytes.NewBuffer(buff)
	switch bytesdataType {
	case "uint16":
		var tmp uint16
		for _, v := range properties {
			if err := binary.Read(byteBuffer, binary.BigEndian, &tmp);err != nil {
				s := fmt.Sprintf("read error , %v", err)
				return nil, errors.New(s)
			}
			result[v] = float32(tmp)
			fmt.Printf("表示为十进制%d\n", tmp)
		}
	}
	return result, nil
}


func (mp *ModBusPayload) convertMapToBytes(properties []string, bytesdataType string, datamap map[string]float32)([]byte, error){
	result := make([]byte, 0)

	buffer := bytes.NewBuffer(result)
	switch bytesdataType {
	case "uint8":
		for _, v := range properties {
			if err := binary.Write(buffer, binary.BigEndian, uint8(datamap[v])); err != nil {
				s := fmt.Sprintf("Pack error , %v", err)
				return nil, errors.New(s)
			}
		}
	}
	return buffer.Bytes(), nil
}



const (
	// Bit access
	FuncCodeReadDiscreteInputs = 2
	FuncCodeReadCoils          = 1
	FuncCodeWriteSingleCoil    = 5
	FuncCodeWriteMultipleCoils = 15

	// 16-bit access
	FuncCodeReadInputRegisters         = 4
	FuncCodeReadHoldingRegisters       = 3
	FuncCodeWriteSingleRegister        = 6
	FuncCodeWriteMultipleRegisters     = 16
	FuncCodeReadWriteMultipleRegisters = 23
	FuncCodeMaskWriteRegister          = 22
	FuncCodeReadFIFOQueue              = 24
)

