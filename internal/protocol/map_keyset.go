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

import ()

type mapKeySetCodec struct {
}

func (self *mapKeySetCodec) CalculateSize(args ...interface{}) (dataSize int) {
	dataSize += StringCalculateSize(args[0].(*string))
	return
}
func (self *mapKeySetCodec) EncodeRequest(args ...interface{}) (request *ClientMessage) {
	// Encode request into clientMessage
	request = NewClientMessage(nil, self.CalculateSize(args...))
	request.SetMessageType(MAP_KEYSET)
	request.IsRetryable = true
	request.AppendString(args[0].(*string))
	request.UpdateFrameLength()
	return
}

func (self *mapKeySetCodec) DecodeResponse(clientMessage *ClientMessage, toObject ToObject) (parameters interface{}, err error) {

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
