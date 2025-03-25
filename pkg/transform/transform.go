package transform

import (
	"bytes"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/RomanAgaltsev/keyper/internal/model"
	pb "github.com/RomanAgaltsev/keyper/pkg/keyper/v1"
)

func PbToUser(src *pb.UserCredentials) *model.User {
	if src == nil {
		return nil
	}

	result := model.User{
		Login:    src.Login,
		Password: src.Password,
	}

	return &result
}

func PbToSecret(src *pb.Secret) *model.Secret {
	if src == nil {
		return nil
	}

	secretType := model.SecretTypeUNSPECIFIED
	switch src.Type {
	case 1:
		secretType = model.SecretTypeCREDENTIALS
	case 2:
		secretType = model.SecretTypeTEXT
	case 3:
		secretType = model.SecretTypeBINARY
	case 4:
		secretType = model.SecretTypeCARD
	}

	result := &model.Secret{
		ID:        uuid.MustParse(*src.Id),
		Name:      src.Name,
		Type:      secretType,
		Metadata:  src.Metadata,
		Data:      bytes.Join(src.Data, nil),
		Comment:   *src.Comment,
		CreatedAt: src.CreatedAt.AsTime(),
		UpdatedAt: src.UpdatedAt.AsTime(),
	}

	return result
}

func SecretToPb(src *model.Secret) *pb.Secret {
	secretID := src.ID.String()

	var secretType pb.SecretType = 0
	switch src.Type {
	case model.SecretTypeCREDENTIALS:
		secretType = 1
	case model.SecretTypeTEXT:
		secretType = 2
	case model.SecretTypeBINARY:
		secretType = 3
	case model.SecretTypeCARD:
		secretType = 4
	}

	result := pb.Secret{
		Id:        &secretID,
		Name:      src.Name,
		Type:      secretType,
		Metadata:  src.Metadata,
		Data:      append([][]byte(nil), src.Data),
		Comment:   &src.Comment,
		CreatedAt: timestamppb.New(src.CreatedAt),
		UpdatedAt: timestamppb.New(src.UpdatedAt),
	}

	return &result
}
