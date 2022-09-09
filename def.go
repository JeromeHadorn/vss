package vss

import "time"

type Snapshot struct {
	Id               string
	DeviceObjectPath string
	Details          SnapshotDetails
	BaseFolder       string
	Drive            string
}

type SnapshotDetails struct {
	Id                 string
	ProviderID         string
	Status             string
	DeviceObject       string
	VolumeName         string
	OriginatingMachine string
	ServiceMachine     string
	ExposedName        string
	Attributes         interface{}
	InstallDate        time.Time
}

type ISnapshotter interface {
	CreateSnapshot(drive string, timeout int, force bool) (*Snapshot, error)
	Details(id string) (*Snapshot, error)
	DeleteSnapshot(id string) error
}
