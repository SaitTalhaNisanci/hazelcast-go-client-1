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

package record

import (
	"sync/atomic"
	"time"

	"github.com/hazelcast/hazelcast-go-client/internal/nearcache"
)

type AbstractNearCacheRecord struct {
	partitionID    int32
	sequence       int64
	uuid           atomic.Value
	value          atomic.Value
	expirationTime atomic.Value
	creationTime   atomic.Value
	accessTime     atomic.Value
	recordState    int64
	accessHit      int64
}

func NewAbstractNearCacheRecord(value interface{}, creationTime time.Duration,
	expirationTime time.Duration) *AbstractNearCacheRecord {
	a := &AbstractNearCacheRecord{}
	a.value.Store(value)
	a.creationTime.Store(creationTime)
	a.expirationTime.Store(expirationTime)
	a.accessTime.Store(nearcache.TimeNotSet)
	a.uuid.Store("")
	atomic.StoreInt64(&a.recordState, nearcache.ReadPermitted)
	return a
}

func (a *AbstractNearCacheRecord) Value() interface{} {
	return a.value.Load()
}

func (a *AbstractNearCacheRecord) SetValue(value interface{}) {
	a.value.Store(value)
}

func (a *AbstractNearCacheRecord) SetCreationTime(time time.Duration) {
	a.creationTime.Store(time)
}

func (a *AbstractNearCacheRecord) SetAccessTime(time time.Duration) {
	a.accessTime.Store(time)
}

func (a *AbstractNearCacheRecord) IsIdleAt(maxIdleTime time.Duration, now time.Duration) bool {
	if maxIdleTime > 0 {
		accessTime := a.accessTime.Load().(time.Duration)
		creationTime := a.creationTime.Load().(time.Duration)
		if accessTime > nearcache.TimeNotSet {
			return accessTime+maxIdleTime < now
		}
		return creationTime+maxIdleTime < now

	}
	return false
}

func (a *AbstractNearCacheRecord) IncrementAccessHit() {
	atomic.AddInt64(&a.accessHit, 1)
}

func (a *AbstractNearCacheRecord) RecordState() int64 {
	return atomic.LoadInt64(&a.recordState)
}

func (a *AbstractNearCacheRecord) CasRecordState(expect int64, update int64) {
	atomic.CompareAndSwapInt64(&a.recordState, expect, update)
}

func (a *AbstractNearCacheRecord) PartitionID() int32 {
	return atomic.LoadInt32(&a.partitionID)
}

func (a *AbstractNearCacheRecord) SetPartitionID(partitionID int32) {
	atomic.StoreInt32(&a.partitionID, partitionID)
}

func (a *AbstractNearCacheRecord) InvalidationSequence() int64 {
	return atomic.LoadInt64(&a.sequence)
}

func (a *AbstractNearCacheRecord) SetInvalidationSequence(sequence int64) {
	atomic.StoreInt64(&a.sequence, sequence)
}

func (a *AbstractNearCacheRecord) SetUUID(UUID string) {
	a.uuid.Store(UUID)
}

func (a *AbstractNearCacheRecord) HasSameUUID(UUID string) bool {
	uuid := a.uuid.Load().(string)
	return uuid != "" && UUID != "" && uuid == UUID
}