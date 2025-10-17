package openapi3_util

import (
	"context"
	"errors"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

func LoadFromData(data []byte) (*openapi3.T, error) {
	return openapi3.NewLoader().LoadFromData(data)
}

func ValidateSchema(ctx context.Context, data []byte) error {
	doc, err := LoadFromData(data)
	if err != nil {
		return err
	}
	return ValidateDoc(ctx, doc)
}

func ValidateDoc(ctx context.Context, doc *openapi3.T) error {
	if doc == nil {
		return errors.New("schema nil")
	}
	// check servers
	if len(doc.Servers) == 0 {
		return errors.New("schema servers empty")
	}
	// check operationId
	for path, pathItem := range doc.Paths {
		for method, operation := range pathItem.Operations() {
			if operation.OperationID == "" {
				return fmt.Errorf("schema path(%v) method(%v) operationId empty", path, method)
			}
		}
	}
	return doc.Validate(ctx)
}
