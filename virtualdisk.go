/*
 * =============================================================================================
 * IBM Confidential
 * © Copyright IBM Corp. 2019-2021
 * The source code for this program is not published or otherwise divested of its trade secrets,
 * irrespective of what has been deposited with the U.S. Copyright Office.
 * =============================================================================================
 */

package v1beta1

import (
	storageapis "github.ibm.com/genctl/apis/storage/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VirtualDisk
// +k8s:openapi-gen=true
// +resource:path=virtualdisks,strategy=VirtualDiskStrategy
type VirtualDisk struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VirtualDiskSpec   `json:"spec"`
	Status VirtualDiskStatus `json:"status,omitempty"`
}

type BusType string

const (
	BusTypeVirtio         BusType = "virtio"
	BusTypeScsi           BusType = "scsi"
	BusTypeIde            BusType = "ide"
	BusTypeDefault        BusType = BusTypeVirtio
	BusTypeWindowsDefault BusType = BusTypeScsi
)

var BusTypeFromString = map[string]BusType{
	"virtio": BusTypeVirtio,
	"scsi":   BusTypeScsi,
	"ide":    BusTypeIde,
}

type ImageType string

const (
	ImageTypeRaw   ImageType = "RAW"
	ImageTypeQcow2 ImageType = "QCOW2"
	ImageTypeVpc   ImageType = "VPC"
	ImageTypeVhd   ImageType = "VHD"
)

type BackingSourceType string

const (
	BackingSourceTypeFile    BackingSourceType = "file"
	BackingSourceTypeBlock   BackingSourceType = "block"
	BackingSourceTypeDir     BackingSourceType = "dir"
	BackingSourceTypeNetwork BackingSourceType = "network"
	BackingSourceTypeVolume  BackingSourceType = "volume"
)

type VolumeAttachmentSecurityLevelType string

const (
	VolumeAttachmentSecurityLevelStandard VolumeAttachmentSecurityLevelType = ""
	VolumeAttachmentSecurityLevelKrb      VolumeAttachmentSecurityLevelType = "krb5p"
)

type OperatingSystem struct {
	// Name is the display name of the os (ex. Ubuntu Server 16 LTS)
	Name string `json:"name"`

	// Vendor is the os' vendor (ex. Canonical)
	Vendor string `json:"vendor"`

	// Version is the specific version of the os (ex. 16.04.3)
	Version string `json:"version"`
}

type VirtualDiskSpec struct {
	// Resource Name is the human readable name for the VirtualDisk
	ResourceName string `json:"resourceName"`

	// Name is used to uniquely identify a VirtualDisk
	Name string `json:"name"`

	// One and only one of Volume or VolumeSpec or LocalDiskSpec must be specified.
	// The resulting virtual disk will either be a Volume created from the VolumeSpec or
	// the actual Volume specified here or a StorageDevice created from the LocalDiskSpec.
	Volume        *apimachinerytypes.NamespacedName `json:"volume,omitempty"`
	VolumeSpec    *storageapis.VolumeSpec           `json:"volumeSpec,omitempty"`
	LocalDiskSpec *LocalDiskSpec                    `json:"localDiskSpec,omitempty"`

	// VirtualMachine points to the VM to which this VirtualDisk is attached to.  For VirtualDiskSpecs
	// included within a VirtualMachineSpec, this field can be left empty as it is always filled with
	// that VirtualMachine.  This field is needed for the attach after VM start case.
	VirtualMachine *apimachinerytypes.NamespacedName `json:"virtualMachine,omitempty"`

	// Specifies where this VirtualDisk is to be assigned.
	// +optional
	Node *Node `json:"node,omitempty"`

	// RequiredForBoot specifies that this disk is must be ready when a VM is booted up
	RequiredForBoot bool `json:"requiredForBoot,omitempty"`

	// Indicates whether this VirtualDisk is the boot volume for a VM.
	// There must be at least one and only one VirtualDisk that is bootable before a VM can run.
	// +optional
	Bootable bool `json:"bootable,omitempty"`

	// Set the vdisk bus type for VMs
	// +optional
	BusType BusType `json:"busType,omitempty"`

	// ImageFormat specifies the format of the image associated with the boot volume
	// +optional
	ImageFormat ImageType `json:"imageFormat"`

	// ImageOS specifies the Operating System of the image associated with the boot volume
	// +optional
	ImageOS *OperatingSystem `json:"imageOS,omitempty"`

	// ExternalVolumeID identifies where the volume resides
	// +optional
	ExternalVolumeID string `json:"externalVolumeID,omitempty"`

	// VolumeAttachmentSecurityLevel represents the desired level of security/encryption for volume attachment
	// +optional
	VolumeAttachmentSecurityLevel VolumeAttachmentSecurityLevelType `json:"volumeAttachmentSecurityLevel,omitempty"`

	// AutoDelete is a boolean specifying that the disk should be deleted when
	// the VM is deleted
	// +optional
	AutoDelete bool `json:"autoDelete,omitempty"`

	// MaxIOPS is directly passed to the StorageDevice as a parameter
	MaxIOPS int32 `json:"maxIOPS"`

	// BandwidthMbps is the total bandwidth of this VirtualDisk
	BandwidthMbps uint32 `json:"bandwidthMbps,omitempty"`
}

type LocalDiskSpec struct {
	// Node specifies where the NFS volume attachment will be performed.
	// +required
	Node Node `json:"node"`

	// RemoteAddress is set when a remote path is to be mounted locally
	// +optional
	RemoteAddress string `json:"remoteAddress,omitempty"`

	// SourcePath is where the disk is located currently
	// +required
	SourcePath string `json:"sourcePath"`

	// MountOptions is used only when RemoteAddress is set
	// +optional
	MountOptions string `json:"mountOptions,omitempty"`

	// ImageType specifies what the format of the image at SourcePath is.
	// It will be removed in the near future. Use ImageFormat in VirtualDiskSpec instead
	// +optional
	ImageType ImageType `json:"imageType"`

	// MaxIOPS is the maximum data transfer rate (1 IOP = 4 KiB/s) to reserve
	// for the disk, or zero if no limit
	// +optional
	MaxIOPS int32 `json:"maxIOPS"`

	// OperatingSystem is metadata relating to the os on this image (for boot disks)
	// +optional
	OperatingSystem *OperatingSystem `json:"operatingSystem,omitempty"`

	// ImageID is metadata relating to what image is the base of the file at SourcePath
	// +optional
	ImageID string `json:"imageID,omitempty"`
}

// Backing Source specification
type BackingSource struct {
	// Type specified the backing source type
	// +required
	Type BackingSourceType `json:"type"`

	// DevicePath specifies the device path where the backing source is attached.
	//   E.g. /mnt/cci_root/backing_files_001/034df24b-481f-4e79-8929-fb9811c0bd9b/8d6ceb77-040f-417d-9e39-4744eb78537b.qcow2"
	// +required
	DevicePath string `json:"devicePath"`

	// secretUUID specifies the UUID of the libvirt secret associated to the backing source.
	// +optional
	SecretUUID string `json:"secretUuid"`

	// EncryptionWdek specifies the encryption Wdek used at a specific layer in the backing chain.
	// When the StorageDevice is mounting a file to the host it needs access to Encryption Wdek and CRN information.
	// StorageDevice obtains this information by traversing the StorageLayers in a backing chain.
	// It uses this information to create a libvirt secretUUID which is passed from StorageDevice's BackingSources to VirtualDisk's BackingSources.
	// +optional
	EncryptionWdek string `json:"encryptionWdek"`

	// EncryptionWdek specifies the encryption Wdek used at a specific layer in the backing chain.
	// When the StorageDevice is mounting a file to the host it needs access to Encryption Wdek and CRN information.
	// StorageDevice obtains this information by traversing the StorageLayers in a backing chain.
	// It uses this information to create a libvirt secretUUID which is passed from StorageDevice's BackingSources to VirtualDisk's BackingSources.
	// +optional
	EncryptionKeyCrn string `json:"encryptionKeyCrn"`

	// IsSnapshot specifies whether or not this layer of the backing sources is originating from a snapshot layer object
	// +optional
	IsSnapshot bool `json:"isSnapshot"`

	// IsImage specifies whether or not this layer of the backing sources is originating from an Image. Needed for Custom Encrypted Images.
	// +optional
	IsImage bool `json:"isImage"`
}

type VirtualDiskStatus struct {
	// SpecErrors describes errors in the spec, if any.
	// While the spec has errors, the implementation will not
	// do anything more than identify them.
	// Example errors include referencing a deleted object.
	// This is not a historical list that grows over time;
	// rather, it is the current list.
	// Il8n TBD.
	// +optional
	// +patchStrategy=replace
	SpecErrors []string `json:"specErrors,omitempty" patchStrategy:"replace"`

	// Volume is a reference to the Volume object backing this virtual disk.  The Volume
	// may either be a clone of a snapshot or image or an actual Volume specified in
	// the VirtualDiskSpec.
	Volume apimachinerytypes.NamespacedName `json:"volume"`

	// VolumeProvisioned indicates whether the volume is actually ready or not.
	VolumeProvisioned bool `json:"volumeProvisioned,omitempty"`

	// StorageType is a mirror of the StorageType from the volume, so we know how to interface
	// with the volume.
	StorageType storageapis.StorageType `json:"storageType,omitempty"`

	// MaxIOPS is drawn from the status of the StorageDevice for projection to the region.
	MaxIOPS int32 `json:"maxIOPS,omitempty"`

	// Node is set when the virtualmachine has been assigned to start on a node.  This signals
	// that the storage attachment should be created.  Node is set to empty when the VM is stopped
	// or when a virtualdisk is detached from the VM.
	Node *Node `json:"node,omitempty"`

	// ReadyToAttach indicates whether storage attachment is known to be ready for use.
	ReadyToAttach bool `json:"readyToAttach,omitempty"`

	// Path represents the actual path to find the disk on the compute node.  It must be filled when ReadyToAttach==true.
	Path string `json:"path,omitempty"`

	// Attached indicates that the virtualdisk is attached to the VM.  When a virtualdisk is detached
	// from a VM or the VM is stopped, Attached is set to false.
	Attached bool `json:"attached,omitempty"`

	// HypervisorAlias stores the alias given to the VirtualDisk by the hypervisor when the device is attached.
	HypervisorAlias string `json:"hypervisorAlias,omitempty"`

	// PCIAddress stores the pci address (in BDF format) given to the VirtualNic by the hypervisor when the device is attached.
	// Example: "0000:00:4.0"
	PCIAddress string `json:"pciAddress,omitempty"`

	// CCWAddress stores the ccw address given to the VirtualDisk by the hypervisor when the device is attached.
	// Example: "fe:0:0000"
	CCWAddress string `json:"ccwAddress,omitempty"`

	// Failed indicates that the implementation knows this VirtualDisk will henceforth not be functional.
	// This covers all such cases.
	// +optional
	Failed bool `json:"failed,omitempty"`

	// StorageDevice is the name of the storagedevice that has been created
	// or will be created by this virtualdisk.
	// Only used when Volume or VolumeSpec is provided in the spec.
	// +optional
	StorageDevice string `json:"storageDevice,omitempty"`

	// Messages are intended for Ops consumption and can carry explanations of the current status.
	// This is not a historical list that grows over time;
	// rather, it is the current list.
	// Il8n TBD.
	// +optional
	// +patchStrategy=replace
	Messages []string `json:"messages,omitempty" patchStrategy:"replace"`

	// StatusReasons are intended to be decoded by the regional api server and contain
	// status codes describing the current status reason.
	// This is not a historical list that grows over time;
	// rather, it is the current list.
	// Il8n TBD.
	// +optional
	// +patchStrategy=replace
	StatusReasons []StatusReason `json:"statusReasons,omitempty" patchStrategy:"replace"`

	// Set the vdisk bus type for VMs
	// +optional
	BusType BusType `json:"busType,omitempty"`

	// OperatingSystem is metadata relating to the os on this image (for boot disks)
	// +optional
	OperatingSystem *OperatingSystem `json:"operatingSystem,omitempty"`

	// ImageID is metadata relating to what image is the base of the file at SourcePath
	// +optional
	ImageID string `json:"imageID,omitempty"`

	// ImageFormat specifies the format of the image associated with the boot volume
	// +optional
	ImageFormat ImageType `json:"imageFormat"`

	// ExternalVolumeID identifies where the volume resides
	// +optional
	ExternalVolumeID string `json:"externalVolumeID,omitempty"`

	// SecretDekUUID is the UUID of the secret to decrypt a file (if encrypted)
	// +optional
	SecretDekUUID string `json:"secretDekUuid,omitempty" vet:"ignore"`

	// List of backing sources.
	//  Only support Type 2 volume provision, implies there is at least one backing source.
	//  The order of backing source in the array is important. The backing file chain is linked based on
	//  the order of the backing source in the array; element 0 is the head and n is the tail.
	// +optional
	BackingSources []BackingSource `json:"backingSources,omitempty"`

	// BandwidthMbps is the total bandwidth of this VirtualDisk. This status is primarily used in projection to the
	// region.
	BandwidthMbps uint32 `json:"bandwidthMbps,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type VirtualDiskList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	// List of virtual disks.
	Items []VirtualDisk `json:"items"`
}
