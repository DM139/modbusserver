package test

import (
	"modbusserver/src/payload"
	"testing"
)

func TestEncode(t *testing.T){

	functionCode := uint8(0x03)
	startAddress := uint16(0x0001)
	quantity := uint16(0x0008)

	properties := []string{"v1", "v2", "v3", "v4"}
	bytesdataType := "uint8"

	mbp := &payload.ModBusPayload{functionCode, startAddress, quantity}
	rhr := &payload.ReadHoldingRegistersHandler{ModBusPayload: mbp, Properties: properties, BytesdataType: bytesdataType}
	var ph payload.PayloadHandler = rhr
	buf, err := ph.Encode();
	if err != nil{
		t.Error(err)
	}
	t.Log(buf)

}

func TestDecode(t *testing.T){
	functionCode := uint8(0x03)
	startAddress := uint16(0x0001)
	quantity := uint16(0x0008)

	properties := []string{"v1", "v2", "v3", "v4",
							"v5", "v6", "v7", "v8",
							}
	bytesdataType := "uint16"
	buff := []byte{0x03, 0x10,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
	}

	mbp := &payload.ModBusPayload{functionCode, startAddress, quantity}
	rhr := &payload.ReadHoldingRegistersHandler{ModBusPayload: mbp, Properties: properties, BytesdataType: bytesdataType}
	var ph payload.PayloadHandler = rhr
	err := ph.Decode(buff);
	if err != nil{
		t.Error(err)
	}
}
