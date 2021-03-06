// Copyright © 2017 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: BSD-2-Clause
package vds

import (
	"fmt"
	"path"
	"strings"
	"sync"

	"github.com/vmware/virtual-security-module/config"
	"github.com/vmware/virtual-security-module/util"
)

const inMemoryDSType = "InMemoryDataStore"

func init() {
	if err := DataStoreRegistrar.Register(inMemoryDSType, NewInMemoryDS()); err != nil {
		panic(fmt.Sprintf("Failed to register data store type %v: %v", inMemoryDSType, err))
	}
}

// An implementation of a datastore in-memory.
// Useful for testing. Not recommended for production!!
type InMemoryDS struct {
	entryMap map[string]*DataStoreEntry
	mutex    sync.Mutex
}

func NewInMemoryDS() *InMemoryDS {
	return &InMemoryDS{
		entryMap: make(map[string]*DataStoreEntry),
	}
}

func (ds *InMemoryDS) Init(storeConfig *config.DataStoreConfig) error {
	return nil
}

func (ds *InMemoryDS) CompleteInit(storeConfig *config.DataStoreConfig) error {
	return nil
}

func (ds *InMemoryDS) Initialized() bool {
	return true
}

func (ds *InMemoryDS) CreateEntry(entry *DataStoreEntry) error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	_, ok := ds.entryMap[entry.Id]
	if ok {
		return util.ErrAlreadyExists
	}

	buf := make([]byte, len(entry.Data))
	copy(buf, entry.Data)

	dsEntry := &DataStoreEntry{
		Id:       entry.Id,
		Data:     buf,
		MetaData: entry.MetaData,
	}

	ds.entryMap[entry.Id] = dsEntry

	return nil
}

func (ds *InMemoryDS) ReadEntry(entryId string) (*DataStoreEntry, error) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	entry, ok := ds.entryMap[entryId]
	if !ok {
		return nil, util.ErrNotFound
	}

	buf := make([]byte, len(entry.Data))
	copy(buf, entry.Data)

	dsEntry := &DataStoreEntry{
		Id:       entry.Id,
		Data:     buf,
		MetaData: entry.MetaData,
	}

	return dsEntry, nil
}

func (ds *InMemoryDS) DeleteEntry(entryId string) error {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	_, ok := ds.entryMap[entryId]
	if !ok {
		return util.ErrNotFound
	}

	delete(ds.entryMap, entryId)

	return nil
}

func (ds *InMemoryDS) SearchChildEntries(parentEntryId string) ([]*DataStoreEntry, error) {
	ds.mutex.Lock()
	defer ds.mutex.Unlock()

	pattern := parentEntryId + "/?*"
	if strings.HasSuffix(parentEntryId, "/") {
		pattern = parentEntryId + "?*"
	}

	dsEntries := make([]*DataStoreEntry, 0)

	for entryId, entry := range ds.entryMap {
		matched, err := path.Match(pattern, entryId)
		if err != nil {
			return dsEntries, err
		}

		if matched {
			buf := make([]byte, len(entry.Data))
			copy(buf, entry.Data)

			dsEntry := &DataStoreEntry{
				Id:       entry.Id,
				Data:     buf,
				MetaData: entry.MetaData,
			}

			dsEntries = append(dsEntries, dsEntry)
		}
	}

	return dsEntries, nil
}

func (ds *InMemoryDS) Type() string {
	return inMemoryDSType
}

func (ds *InMemoryDS) Location() string {
	return ""
}
