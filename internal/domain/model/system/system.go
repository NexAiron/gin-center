package system

import baseModel "gin-center/internal/domain/model/base"

type System struct {
	baseModel.BaseModel
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

func NewSystem(name, description, version string) *System {
	return &System{
		Name:        name,
		Description: description,
		Version:     version,
	}
}
