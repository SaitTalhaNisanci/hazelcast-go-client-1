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
	. "github.com/hazelcast/hazelcast-go-client/internal/common"
)

type listSubCodec struct {
}

func (self *listSubCodec) CalculateSize(args ...interface{}) (dataSize int) {
	// Calculates the request payload size
	dataSize += StringCalculateSize(args[0].(*string))
	dataSize += INT32_SIZE_IN_BYTES
	dataSize += INT32_SIZE_IN_BYTES
	return
}
func (self *listSubCodec) EncodeRequest(args ...interface{}) (request *ClientMessage) {
	// Encode request into clientMessage
	request = NewClientMessage(nil, self.CalculateSize(args))
	request.SetMessageType(LIST_SUB)
	request.IsRetryable = true
	request.AppendString(args[0].(*string))
	request.AppendInt32(args[1].(int32))
	request.AppendInt32(args[2].(int32))
	request.UpdateFrameLength()
	return
}

func (self *listSubCodec) DecodeResponse(clientMessage *ClientMessage, toObject ToObject) (parameters interface{}, err error) {
	// Decode response from client message

	responseSize := clientMessage.ReadInt32()
	response := make([]interface{}, responseSize)
	for responseIndex := 0; responseIndex < int(responseSize); responseIndex++ {
		responseItem, err := toObject(clientMessage.ReadData())
		if err != nil {
			return nil, err
		}
		response[responseIndex] = responseItem
	}
	parameters = response

	return
}
