//go:build windows
// +build windows
package vss

import (
	"errors"
)

type Snapshotter struct{}

func (v *Snapshotter) CreateSnapshot(drive string, timeout int, force bool) (*Snapshot, error) {
	return nil, errors.New("no snapshot can be created on non windows systems")
}

func (v *Snapshotter) Details(id string) (*Snapshot, error) {
	return nil, errors.New("no snapshot details can be returned on a non windows systems")
}

func (v *Snapshotter) DeleteSnapshot(id string) error {
	return errors.New("no snapshot can be deleted on non windows systems")
}
