package converter

import (
    recordsv1 "github.com/RomanAgaltsev/keyper/pkg/keyper"
	"github.com/RomanAgaltsev/keyper/server/internal/model"
)

func ToRecordsFromServiceRecord(record *model.Record) *recordsv1.Record {
	//	credencials
	//
	//	return &recordsv1.GetResponse{
	//		Id:      uint64(record.ID),
	//		Type:    uint32(record.Type),
	//		Address: record.Address,
	//		Credentials:
	return &recordsv1.Record{}
}
