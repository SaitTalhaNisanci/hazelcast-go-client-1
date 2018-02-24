// Copyright (c) 2008-2018, Hazelcast, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package protocol

import (
	. "github.com/hazelcast/hazelcast-go-client/internal/serialization"

	. "github.com/hazelcast/hazelcast-go-client/internal/common"
)

type mapGetAllCodec struct {
}

func (self *mapGetAllCodec) CalculateSize(args ...interface{}) (dataSize int) {
	dataSize += StringCalculateSize(args[0].(*string))
	dataSize += INT_SIZE_IN_BYTES
	for _, keysItem := range args[1].([]*Data) {
		dataSize += DataCalculateSize(keysItem)
	}
	return
}
func (self *mapGetAllCodec) EncodeRequest(args ...interface{}) (request *ClientMessage) {
	// Encode request into clientMessage
	request = NewClientMessage(nil, self.CalculateSize(args...))
	request.SetMessageType(MAP_GETALL)
	request.IsRetryable = false
	request.AppendString(args[0].(*string))
	request.AppendInt(len(args[1].([]*Data)))
	for _, keysItem := range args[1].([]*Data) {
		request.AppendData(keysItem)
	}
	request.UpdateFrameLength()
	return
}

func (self *mapGetAllCodec) DecodeResponse(clientMessage *ClientMessage, toObject ToObject) (parameters interface{}, err error) {

	responseSize := clientMessage.ReadInt32()
	response := make([]*Pair, responseSize)
	for responseIndex := 0; responseIndex < int(responseSize); responseIndex++ {
		var responseItem = &Pair{}
		responseItemKey, err := toObject(clientMessage.ReadData())
		if err != nil {
			return nil, err
		}
		responseItemVal, err := toObject(clientMessage.ReadData())
		if err != nil {
			return nil, err
		}
		responseItem.key = responseItemKey
		responseItem.value = responseItemVal
		response[responseIndex] = responseItem
	}
	parameters = response

	return
}
