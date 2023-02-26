//go:build !windows
// +build !windows

package vss

import (
	"errors"
)

type Snapshotter struct{}

func (v *Snapshotter) CreateSnapshot(drive string, timeout int, force bool) (*Snapshot, error) {
	return nil, errors.New("not yet implemented")
}

func (v *Snapshotter) DeleteSnapshot(snapshotId string) error {
	return errors.New("not yet implemented")
}

func (v *Snapshotter) deleteOnFailure(finished *bool, snapshotId string) {
	if !*finished {
		v.DeleteSnapshot(snapshotId)
	}
}
