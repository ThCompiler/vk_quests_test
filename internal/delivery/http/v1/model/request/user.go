package request

import (
	"github.com/miladibra10/vjson"
	"vk_quests/internal/pkg/evjson"
)

type User struct {
	Name string `json:"name" swaggertype:"string" example:"User"`
}

func ValidateUser(data []byte) error {
	schema := evjson.NewSchema(
		vjson.String("name").Required(),
	)
	return schema.ValidateBytes(data)
}
