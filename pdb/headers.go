package pdb

import "time"

const PalmDBHeaderLength = 78 // 0x4F

type PalmDBHeader struct {
	Name               [32]byte
	FileAttributes     uint16
	Version            uint16
	CreationTime       uint32 // Timestamp
	ModificationTime   uint32 // Timestamp
	BackupTime         uint32 // Timestamp
	ModificationNumber uint32
	AppInfo            uint32
	SortInfo           uint32
	Type               [4]byte // BOOK
	Creator            [4]byte // MOBI
	LastRecordUID      uint32
	NextRecordList     uint32 // Always zero
	NumRecords         uint16
}

func NewPalmDBHeader(name string, dateTime time.Time, numRecords uint16, lastRecordUID uint32) PalmDBHeader {
	nameBytes := [32]byte{}
	copy(nameBytes[:31], name)
	palmTime := convertToPalmTime(dateTime)

	return PalmDBHeader{
		Name:               nameBytes,
		FileAttributes:     0,
		Version:            0,
		CreationTime:       palmTime,
		ModificationTime:   palmTime,
		BackupTime:         palmTime,
		ModificationNumber: 0,
		AppInfo:            0,
		SortInfo:           0,
		Type:               [4]byte{'B', 'O', 'O', 'K'},
		Creator:            [4]byte{'M', 'O', 'B', 'I'},
		LastRecordUID:      lastRecordUID,
		NextRecordList:     0,
		NumRecords:         numRecords,
	}
}

const RecordHeaderLength = 8 // 0x08

type RecordHeader struct {
	Offset    uint32
	Attribute byte
	Skip      byte
	UniqueID  uint16
}
