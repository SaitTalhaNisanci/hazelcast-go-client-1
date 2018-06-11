// Copyright (c) 2008-2018, Hazelcast, Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License")
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

package discovery

import (
	"log"

	"time"

	"github.com/hazelcast/hazelcast-go-client/core"
)

// HzCloudAddrTranslator is used to translate private addresses to public addresses.
type HzCloudAddrTranslator struct {
	cloudDiscovery  *HazelcastCloud
	privateToPublic map[string]core.Address
}

// NewHzCloudAddrTranslator returns a HzCloudAddrTranslator with the given parameters.
func NewHzCloudAddrTranslator(cloudToken string, connectionTimeout time.Duration) *HzCloudAddrTranslator {
	return NewHzCloudAddrTranslatorWithCloudDisc(
		NewHazelcastCloud(
			cloudToken,
			connectionTimeout,
		),
	)
}

// NewHzCloudAddrTranslatorWithCloudDisc returns a HzCloudAddrTranslator with the given parameters.
func NewHzCloudAddrTranslatorWithCloudDisc(cloudDisc *HazelcastCloud) *HzCloudAddrTranslator {
	return &HzCloudAddrTranslator{
		cloudDiscovery: cloudDisc,
	}
}

// Translate translates the given addr to its public address.
func (at *HzCloudAddrTranslator) Translate(addr core.Address) core.Address {
	if addr == nil {
		return nil
	}

	if publicAddr, found := at.privateToPublic[addr.String()]; found {
		return publicAddr
	}

	at.Refresh()

	if publicAddr, found := at.privateToPublic[addr.String()]; found {
		return publicAddr
	}

	return nil
}

// Refresh refreshes the internal lookup table.
func (at *HzCloudAddrTranslator) Refresh() {
	privateToPublic, err := at.cloudDiscovery.discoverNodes()
	if err != nil {
		log.Println("Failed to load addresses from hazelcast.cloud ", err)
	} else {
		at.privateToPublic = privateToPublic
	}
}
