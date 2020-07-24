package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"errors"
	)

func convertBytesToMap(properties []string, bytesdataType string, buff []byte)  (map[string]float32, error){
	result := make(map[string]float32)
	byteBuffer := bytes.NewBuffer(buff)
	switch bytesdataType {
	case "uint8":
		var tmp uint8
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


func convertMapToBytes(properties []string, bytesdataType string, datamap map[string]float32)([]byte, error){
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








func main() {
	oderList := []string{"v1", "v2", "v3", "v4"}
	dataType := "uint8"
	//buff := []byte{0x03, 0x04, 0x05, 0x06, 0x03, 0x04, 0x05, 0x06}

	//m, _ := convertBytesToMap(oderList, dataType, buff)
	//for k, v := range m {
	//	fmt.Printf("%s,%d\n", k, v)
	//}

	testdata := make(map[string]float32)
	i := 1
	for _, v := range oderList {
		testdata[v] = float32(i)
		i++
	}

	n, _ := convertMapToBytes(oderList, dataType, testdata)
	fmt.Printf("frame", n)

}
