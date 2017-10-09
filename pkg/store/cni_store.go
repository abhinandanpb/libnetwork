package store

import (
	"encoding/json"

	"github.com/docker/libnetwork/api"
	"github.com/docker/libnetwork/datastore"
	"github.com/sirupsen/logrus"

	"github.com/docker/libnetwork/client"
)

const (
	cniPrefix = "cni"
)

type CniStore struct {
	PodName          string
	PodNamespace     string
	InfraContainerID string
	SandboxID        string
	EndpointID       string
	SandboxMeta      api.SandboxMetadata
	dbIndex          uint64
	dbExists         bool
}

// Key provides the Key to be used in KV Store
func (cs *CniStore) Key() []string {
	return []string{cniPrefix, cs.PodName, cs.PodNamespace}
}

// KeyPrefix returns the immediate parent key that can be used for tree walk
func (cs *CniStore) KeyPrefix() []string {
	return []string{cniPrefix}
}

// Value marshals the data to be stored in the KV store
func (cs *CniStore) Value() []byte {
	b, err := json.Marshal(cs)
	if err != nil {
		logrus.Warnf("failed to marshal cni store: %v", err)
		return nil
	}
	return b
}

// SetValue unmarshalls the data from the KV store.
func (cs *CniStore) SetValue(value []byte) error {
	return json.Unmarshal(value, cs)
}

// Index returns the latest DB Index as seen by this object
func (cs *CniStore) Index() uint64 {
	return cs.dbIndex
}

// SetIndex method allows the datastore to store the latest DB Index into this object
func (cs *CniStore) SetIndex(index uint64) {
	cs.dbIndex = index
	cs.dbExists = true
}

// Exists method is true if this object has been stored in the DB.
func (cs *CniStore) Exists() bool {
	return cs.dbExists
}

// Skip provides a way for a KV Object to avoid persisting it in the KV Store
func (cs *CniStore) Skip() bool {
	return false
}

func (cs *CniStore) New() datastore.KVObject {
	return &CniStore{}
}

func (cs *CniStore) CopyTo(o datastore.KVObject) error {
	dstCs := o.(*CniStore)
	dstCs.PodName = cs.PodName
	dstCs.PodNamespace = cs.PodNamespace
	dstCs.InfraContainerID = cs.InfraContainerID
	dstCs.SandboxID = cs.SandboxID
	dstCs.EndpointID = cs.EndpointID
	dstCs.SandboxMeta = cs.SandboxMeta
	return nil
}

// DataScope method returns the storage scope of the datastore
func (cs *CniStore) DataScope() string {
	return datastore.LocalScope
}

func CopySandboxMetadata(sbConfig client.SandboxCreate, externalKey string) api.SandboxMetadata {
	var meta api.SandboxMetadata
	meta.ContainerID = sbConfig.ContainerID
	meta.HostName = sbConfig.HostName
	meta.DomainName = sbConfig.DomainName
	meta.HostsPath = sbConfig.HostsPath
	meta.ResolvConfPath = sbConfig.ResolvConfPath
	meta.DNS = sbConfig.DNS
	meta.UseExternalKey = sbConfig.UseExternalKey
	meta.UseDefaultSandbox = sbConfig.UseDefaultSandbox
	meta.ExposedPorts = sbConfig.ExposedPorts
	meta.PortMapping = sbConfig.PortMapping
	meta.ExternalKey = externalKey
	//TODO: skipped extrahosts
	return meta
}
