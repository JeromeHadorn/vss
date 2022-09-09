package vss

import (
	"errors"
	"fmt"
	"time"

	"github.com/jeromehadorn/vss/api"
)

type Snapshot struct {
	Id               string
	DeviceObjectPath string
	Details          SnapshotDetails
	BaseFolder       string
	Drive            string
}

func (s Snapshot) Validate() error {
	if s.Id == "" {
		return fmt.Errorf("snapshot is missing Id property")
	}
	if s.DeviceObjectPath == "" {
		return errors.New("snapshot is missing DeviceObjectPath")
	}
	return nil
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

func (d SnapshotDetails) Validate() error {
	if d.Id == "" {
		return fmt.Errorf("snapshot details are missing ShadowCopyId property")
	}
	if d.ProviderID == "" {
		return fmt.Errorf("snapshot details are missing Provider property")
	}
	if d.Status == "" {
		return fmt.Errorf("snapshot details are missing Status property")
	}
	if d.DeviceObject == "" {
		return fmt.Errorf("snapshot details are missing ShadowCopyVolume property")
	}
	if d.VolumeName == "" {
		return fmt.Errorf("snapshot details are missing OriginalVolume property")
	}
	if d.OriginatingMachine == "" {
		return fmt.Errorf("snapshot details are missing OriginatingMachine property")
	}
	if d.ServiceMachine == "" {
		return fmt.Errorf("snapshot details are missing ServiceMachine property")
	}
	if d.ExposedName == "" {
		return fmt.Errorf("snapshot details are missing ExposedName property")
	}
	if d.Attributes == nil {
		return fmt.Errorf("snapshot details are missing Attributes property")
	}
	if d.InstallDate.IsZero() {
		return fmt.Errorf("snapshot details are missing InstallDate property")
	}
	return nil
}

func ParseProperties(p api.VssSnapshotProperties) (SnapshotDetails, error) {
	return SnapshotDetails{
		Id:                 p.GetSnapshotId(),
		ProviderID:         p.GetProviderId(),
		Status:             p.GetSnapshotStatus(),
		DeviceObject:       p.GetSnapshotDeviceObject(),
		VolumeName:         p.GetOriginalVolume(),
		OriginatingMachine: p.GetOriginatingMachine(),
		ServiceMachine:     p.GetServiceMachine(),
		ExposedName:        p.GetExposedName(),
		Attributes:         p.GetSnapshotAttributes(),
		InstallDate:        p.GetCreationTimeStamp(),
	}, nil
}

type ISnapshotter interface {
	CreateSnapshot(drive string, timeout int, force bool) (*Snapshot, error)
	Details(id string) (*Snapshot, error)
	DeleteSnapshot(id string) error
}
