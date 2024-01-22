package persistence

import "github.com/fiskaly/coding-challenges/signing-service-challenge/api"

type Repository interface {
	Save(entity api.DeviceRequest) // TODO: Use generics or sth
	Update(entity api.DeviceRequest)
}
