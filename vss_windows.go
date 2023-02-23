//go:build windows
// +build windows

package vss

import (
	"errors"
	"fmt"

	"github.com/go-ole/go-ole"
)

type Snapshotter struct{}

func (v *Snapshotter) CreateSnapshot(drive string, timeout int, force bool) (*Snapshot, error) {
	successful := false

	if timeout < 180 {
		timeout = 180
	}

	// Initalize COM Library
	ole.CoInitialize(0)
	defer ole.CoUninitialize()
	vssBackupComponent, err := LoadAndInitVSS()
	if err != nil {
		return nil, err
	}

	if err := vssBackupComponent.SetContext(VSS_CTX_BACKUP); err != nil {
		vssBackupComponent.Release()
		return nil, err
	}

	if err := vssBackupComponent.SetBackupState(VSS_BT_COPY); err != nil {
		return nil, err
	}

	var async *IVssAsync

	if async, err = vssBackupComponent.GatherWriterMetadata(); err != nil {
		return nil, fmt.Errorf("VSS_GATHER - Shadow copy creation failed: GatherWriterMetadata, err: %s", err)
	} else if async == nil {
		return nil, fmt.Errorf("VSS_GATHER - Shadow copy creation failed: GatherWriterMetadata failed to return a valid IVssAsync object")
	}

	if err := async.Wait(timeout); err != nil {
		return nil, fmt.Errorf("VSS_GATHER - Shadow copy creation failed: GatherWriterMetadata didn't finish properly, err: %s", err)
	}

	async.Release()

	if isSupported, err := vssBackupComponent.IsVolumeSupported(drive); err != nil {
		vssBackupComponent.Release()
		return nil, fmt.Errorf("Snapshots are not supported for drive %s, err: %s", drive, err)
	} else if !isSupported {
		vssBackupComponent.Release()
		return nil, fmt.Errorf("Snapshots are not supported for drive %s, err: %s", drive, err)
	}

	var snapshotSetID ole.GUID
	var snapshotID ole.GUID

	if err = vssBackupComponent.StartSnapshotSet(&snapshotSetID); err != nil {
		return nil, fmt.Errorf("VSS_START - Shadow copy creation failed: StartSnapshotSet, err %s", err)
	}

	if err = vssBackupComponent.AddToSnapshotSet(drive, &snapshotID); err != nil {
		return nil, fmt.Errorf("VSS_ADD - Shadow copy creation failed: AddToSnapshotSet, err: %s", err)
	}

	if async, err = vssBackupComponent.PrepareForBackup(); err != nil {
		vssBackupComponent.AbortBackup()
		vssBackupComponent.Release()
		return nil, fmt.Errorf("VSS_PREPARE - Shadow copy creation failed: PrepareForBackup returned, err: %s", err)
	}
	if async == nil {
		vssBackupComponent.AbortBackup()
		vssBackupComponent.Release()
		return nil, fmt.Errorf("VSS_PREPARE - Shadow copy creation failed: PrepareForBackup failed to return a valid IVssAsync object")
	}

	if err := async.Wait(timeout); err != nil {
		vssBackupComponent.AbortBackup()
		vssBackupComponent.Release()
		return nil, fmt.Errorf("VSS_PREPARE - Shadow copy creation failed: PrepareForBackup didn't finish properly, err %s", err)

	}
	async.Release()

	if async, err = vssBackupComponent.DoSnapshotSet(); err != nil {
		vssBackupComponent.AbortBackup()
		vssBackupComponent.Release()
		return nil, fmt.Errorf("VSS_SNAPSHOT - Shadow copy creation failed: DoSnapshotSet, err: %s", err)
	}
	if async == nil {
		vssBackupComponent.AbortBackup()
		vssBackupComponent.Release()
		return nil, fmt.Errorf("VSS_SNAPSHOT - Shadow copy creation failed: DoSnapshotSet failed to return a valid IVssAsync object")
	}

	if err := async.Wait(timeout); err != nil {
		vssBackupComponent.AbortBackup()
		vssBackupComponent.Release()
		return nil, fmt.Errorf("VSS_SNAPSHOT - Shadow copy creation failed: DoSnapshotSet didn't finish properly, err: %s", err)
	}
	async.Release()

	// Gather Properties
	properties := VssSnapshotProperties{}

	if err = vssBackupComponent.GetSnapshotProperties(snapshotID, &properties); err != nil {
		vssBackupComponent.AbortBackup()
		vssBackupComponent.Release()
		return nil, fmt.Errorf("VSS_PROPERTIES - GetSnapshotProperties, err: %s", err)
	}
	details := SnapshotDetails{}
	details, err = ParseProperties(properties)

	snapshot := Snapshot{
		Id:      snapshotID.String(),
		Details: details,
	}

	// Delete Snapshot if an error occurs
	defer v.deleteOnFailure(&successful, snapshot.Id)

	deviceObjectPath := snapshot.Details.DeviceObject + `\`
	snapshot.DeviceObjectPath = deviceObjectPath

	// Check Snapshot is Complete
	if err := snapshot.Validate(); err != nil {
		return nil, err
	}

	// Cancel sheduled deletion
	successful = true

	return &snapshot, nil
}

func (v *Snapshotter) Details(id string) (*Snapshot, error) {
	return nil, errors.New("not yet implemented")
}

func (v *Snapshotter) DeleteSnapshot(snapshotId string) error {
	// Initalize COM Library
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	vssBackupComponent, err := LoadAndInitVSS()
	if err != nil || vssBackupComponent == nil {
		return err
	}

	defer vssBackupComponent.Release()

	id := ole.NewGUID(snapshotId)

	if deletedGUID, _, err := vssBackupComponent.DeleteSnapshots(*id); err != nil {
		vssBackupComponent.AbortBackup()
		vssBackupComponent.Release()
		return fmt.Errorf("VSS_DELETE - Failed to delete the shadow copy: %s\n", deletedGUID.String())
	}
	return nil
}

// Cleanup Method. Called if an error occurs during creation of a snapshot.
func (v *Snapshotter) deleteOnFailure(finished *bool, snapshotId string) {
	if *finished {
		return
	}
	fmt.Printf("Early deleteing snapshot %s, due to error\n", snapshotId)
	if err := v.DeleteSnapshot(snapshotId); err != nil {
		fmt.Printf("error deleting corrupted snapshot %v:\n", err)
	}
}
