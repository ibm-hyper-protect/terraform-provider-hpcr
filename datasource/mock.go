package datasource

import (
	"github.com/terraform-provider-hpcr/fp"
)

type resourceDataMock struct {
	data map[string]any
}

func (mock resourceDataMock) GetOk(key string) (any, bool) {
	value, exists := mock.data[key]
	return value, exists
}

func (mock resourceDataMock) SetId(value string) {
	// noop
}

func (mock resourceDataMock) Set(key string, value any) error {
	mock.data[key] = value
	return nil
}

func CreateResourceDataMock(data map[string]any) fp.ResourceData {
	return resourceDataMock{data: data}
}
