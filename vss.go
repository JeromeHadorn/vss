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

func (v *Snapshotter) Details(id string) (*Snapshot, error) {
	return nil, errors.New("not yet implemented")
}

func (v *Snapshotter) DeleteSnapshot(id string) error {
	return errors.New("not yet implemented")
}

// Cleanup Method. Called if an error occurs during creation of a snapshot.
func (v *Snapshotter) deleteOnFailure(finished *bool, snapshotId string) {

}
