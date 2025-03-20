package internal

//go:generate mockgen --build_flags=--mod=mod -destination=./mock/user/mock_repository.go -package=mock github.com/RomanAgaltsev/keyper/internal/app/keyper-srv/service UserRepository
//go:generate mockgen --build_flags=--mod=mod -destination=./mock/secret/mock_repository.go -package=mock github.com/RomanAgaltsev/keyper/internal/app/keyper-srv/service SecretRepository
