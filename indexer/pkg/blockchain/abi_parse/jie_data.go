package abi_parse

import (
	"github.com/721tools/backend-go/indexer/pkg/utils/hex"
)

var registers []Parse

func init() {
	registers = append(registers, NewERC721(), NewERC20(), NewSeaPort(), NewBlur())
}

type JieData struct {
	tagMethodHandle   map[string]map[string][]JieHandle
	tagEventLogHandle map[string]map[string][]JieHandle
	methodHandle      map[string][]JieHandle
	eventLogHandle    map[string][]JieHandle
}

func NewJieData() *JieData {
	data := &JieData{
		methodHandle:      make(map[string][]JieHandle),
		eventLogHandle:    make(map[string][]JieHandle),
		tagMethodHandle:   make(map[string]map[string][]JieHandle),
		tagEventLogHandle: make(map[string]map[string][]JieHandle)}
	for _, register := range registers {
		if _, ok := data.tagMethodHandle[register.Tag()]; !ok {
			data.tagMethodHandle[register.Tag()] = make(map[string][]JieHandle, 0)
		}
		if _, ok := data.tagEventLogHandle[register.Tag()]; !ok {
			data.tagEventLogHandle[register.Tag()] = make(map[string][]JieHandle, 0)
		}

		for methodID, method := range register.Methods() {
			if _, ok := data.methodHandle[methodID]; !ok {
				data.methodHandle[methodID] = make([]JieHandle, 0)
			}
			if _, ok := data.tagMethodHandle[register.Tag()][methodID]; !ok {
				data.tagMethodHandle[register.Tag()][methodID] = make([]JieHandle, 0)
			}

			data.tagMethodHandle[register.Tag()][methodID] = append(data.tagMethodHandle[register.Tag()][methodID], method)
			data.methodHandle[methodID] = append(data.methodHandle[methodID], method)
		}
		for eventSig, event := range register.EventLogs() {
			if _, ok := data.eventLogHandle[eventSig]; !ok {
				data.eventLogHandle[eventSig] = make([]JieHandle, 0)
			}
			if _, ok := data.tagEventLogHandle[register.Tag()][eventSig]; !ok {
				data.tagEventLogHandle[register.Tag()][eventSig] = make([]JieHandle, 0)
			}
			data.tagEventLogHandle[register.Tag()][eventSig] = append(data.tagEventLogHandle[register.Tag()][eventSig], event)
			data.eventLogHandle[eventSig] = append(data.eventLogHandle[eventSig], event)
		}
	}
	return data
}

func (j *JieData) JieMethod(methodID hex.Hex, input hex.Hex) (methodName string, methodArgs Args) {
	if methodID == nil || input == nil {
		return "", nil
	}
	handles, ok := j.methodHandle[methodID.HexStr()]
	if !ok {
		return "", nil
	}
	for _, handle := range handles {
		methodName, methodArgs = handle(input)
		if methodName != "" {
			return
		}
	}
	return "", nil
}

func (j *JieData) JieEventLogs(topic0 hex.Hex, input hex.Hex) (eventLogName string, eventLogArgs Args) {
	handles, ok := j.eventLogHandle[topic0.HexStr()]
	if !ok {
		return "", nil
	}
	for _, handle := range handles {
		eventLogName, eventLogArgs = handle(input)
		if eventLogName == "" {
			continue
		}
		return
	}
	return "", nil
}

func (j *JieData) JieMethodWithTag(tag string, methodID hex.Hex, input hex.Hex) (methodName string, methodArgs Args) {
	m, ok := j.tagMethodHandle[tag]
	if !ok {
		return "", nil
	}
	handles, ok := m[methodID.HexStr()]
	if !ok {
		return "", nil
	}
	for _, handle := range handles {
		methodName, methodArgs = handle(input)
		if methodName == "" {
			continue
		}
		return
	}
	return "", nil
}

func (j *JieData) JieEventLogsWithTag(tag string, topic0 hex.Hex, input hex.Hex) (eventName string, eventArgs Args) {
	m, ok := j.tagEventLogHandle[tag]
	if !ok {
		return "", nil
	}
	handles, ok := m[topic0.HexStr()]
	if !ok {
		return "", nil
	}
	for _, handle := range handles {
		eventName, eventArgs = handle(input)
		if eventName == "" {
			continue
		}
		return
	}
	return "", nil
}
