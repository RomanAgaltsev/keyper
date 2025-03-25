package transform

import (
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/RomanAgaltsev/keyper/internal/database/queries"
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

func UserToCreateUserParams(src *model.User) queries.CreateUserParams {
	return queries.CreateUserParams{
		Login:    src.Login,
		Password: src.Password,
	}
}

func DBToUser(src queries.User) *model.User {
	return &model.User{
		ID:        src.ID,
		Login:     src.Login,
		Password:  src.Password,
		CreatedAt: src.CreatedAt,
	}
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
		Data:      src.Data,
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
		Data:      src.Data,
		Comment:   &src.Comment,
		CreatedAt: timestamppb.New(src.CreatedAt),
		UpdatedAt: timestamppb.New(src.UpdatedAt),
	}

	return &result
}

func DBToSecret(src queries.Secret) *model.Secret {
	// TODO: copygen - Вынести из репозитория
	return &model.Secret{
		ID:        src.ID,
		Name:      src.Name,
		Type:      model.SecretType(src.Type),
		Metadata:  src.Metadata,
		Data:      src.Data,
		Comment:   *src.Comment,
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
		UserID:    src.UserID,
	}
}

func SecretToCreateSecretParams(src *model.Secret) queries.CreateSecretParams {
	return queries.CreateSecretParams{
		Name:     src.Name,
		Type:     queries.SecretType(src.Type),
		Metadata: src.Metadata,
		Data:     src.Data,
		Comment:  &src.Comment,
		UserID:   src.UserID,
	}
}

func SecretToUpdateSecretParams(src *model.Secret) queries.UpdateSecretParams {
	return queries.UpdateSecretParams{
		ID:        src.ID,
		Name:      src.Name,
		Type:      queries.SecretType(src.Type),
		Metadata:  src.Metadata,
		Data:      src.Data,
		Comment:   &src.Comment,
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
		UserID:    src.UserID,
	}
}

func ListSecretsRowToSecrets(src []queries.ListSecretsRow) model.Secrets {
	result := make([]*model.Secret, 0, len(src))
	for _, secret := range src {
		result = append(result, &model.Secret{
			ID:        secret.ID,
			Name:      secret.Name,
			Type:      model.SecretType(secret.Type),
			Metadata:  secret.Metadata,
			Comment:   *secret.Comment,
			CreatedAt: secret.CreatedAt,
			UpdatedAt: secret.UpdatedAt,
		})
	}

	return result
}
