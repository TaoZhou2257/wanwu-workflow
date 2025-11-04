package openapi3_util

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/getkin/kin-openapi/openapi3"
)

func LoadFromData(ctx context.Context, data []byte) (*openapi3.T, error) {
	doc, err := openapi3.NewLoader().LoadFromData(data)
	if err != nil {
		return nil, err
	}
	if err = ValidateDoc(ctx, doc); err != nil {
		return nil, err
	}
	return doc, err
}

func ValidateSchema(ctx context.Context, data []byte) error {
	_, err := LoadFromData(ctx, data)
	if err != nil {
		return err
	}
	return nil
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

func FilterSchemaOperations(ctx context.Context, data []byte, operationIDs []string) ([]byte, error) {
	doc, err := LoadFromData(ctx, data)
	if err != nil {
		return nil, err
	}
	return FilterDocOperations(doc, operationIDs).MarshalJSON()
}

func FilterDocOperations(doc *openapi3.T, operationIDs []string) *openapi3.T {
	paths := doc.Paths
	doc.Paths = nil
	for path, pathItem := range paths {
		for method, operation := range pathItem.Operations() {
			if slices.Contains(operationIDs, operation.OperationID) {
				doc.AddOperation(path, method, operation)
			}
		}
	}
	return doc
}
