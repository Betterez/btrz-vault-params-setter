package btrzutils

import (
	"errors"
)

// VaultEntry - vault entry data
type VaultEntry struct {
	ServiceName string
	valuePairs  map[string]string
}

// IsEntryOK - checks if this entry is ok
func (ve *VaultEntry) IsEntryOK() bool {
	if ve.ServiceName == "" {
		return false
	}
	if len(ve.valuePairs) < 1 {
		return false
	}
	return true
}

// AddValuePair - add new value pair to the entry
func (ve *VaultEntry) AddValuePair(name, value string) error {
	if _, ok := ve.valuePairs[name]; ok {
		return errors.New("Value already exists")
	}
	ve.valuePairs[name] = value
	return nil
}

// HasName - HasName
func (ve *VaultEntry) HasName(name string) bool {
	_, ok := ve.valuePairs[name]
	return ok
}
