// Code generated by noice. DO NOT EDIT.
package errors

import (
	bytes "bytes"
	cherry "github.com/containerum/cherry"
	template "text/template"
)

const ()

// ErrAdminRequired error
// User is not admin and has no permissions
func ErrAdminRequired(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Admin access required", StatusHTTP: 403, ID: cherry.ErrID{SID: "volume-manager", Kind: 0x1}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrRequiredHeadersNotProvided(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Required headers not provided", StatusHTTP: 400, ID: cherry.ErrID{SID: "volume-manager", Kind: 0x2}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

// ErrRequestValidationFailed error
// Validation error when parsing request
func ErrRequestValidationFailed(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Request validation failed", StatusHTTP: 400, ID: cherry.ErrID{SID: "volume-manager", Kind: 0x3}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrInternal(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Internal error", StatusHTTP: 500, ID: cherry.ErrID{SID: "volume-manager", Kind: 0x4}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrDatabase(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Database error", StatusHTTP: 500, ID: cherry.ErrID{SID: "volume-manager", Kind: 0x5}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrResourceNotExists(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Resource not exists", StatusHTTP: 404, ID: cherry.ErrID{SID: "volume-manager", Kind: 0x6}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrResourceAlreadyExists(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Resource already exists", StatusHTTP: 400, ID: cherry.ErrID{SID: "volume-manager", Kind: 0x7}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrQuotaExceeded(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Resource quota exceeded", StatusHTTP: 400, ID: cherry.ErrID{SID: "volume-manager", Kind: 0x8}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrNoFreeStorages(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "No free storages found for volume", StatusHTTP: 507, ID: cherry.ErrID{SID: "volume-manager", Kind: 0x9}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrStorageDelete(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Can`t delete storage with volumes", StatusHTTP: 400, ID: cherry.ErrID{SID: "volume-manager", Kind: 0xa}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}

func ErrAccessDenied(params ...func(*cherry.Err)) *cherry.Err {
	err := &cherry.Err{Message: "Access to resource denied", StatusHTTP: 403, ID: cherry.ErrID{SID: "volume-manager", Kind: 0xb}, Details: []string(nil), Fields: cherry.Fields(nil)}
	for _, param := range params {
		param(err)
	}
	for i, detail := range err.Details {
		det := renderTemplate(detail)
		err.Details[i] = det
	}
	return err
}
func renderTemplate(templText string) string {
	buf := &bytes.Buffer{}
	templ, err := template.New("").Parse(templText)
	if err != nil {
		return err.Error()
	}
	err = templ.Execute(buf, map[string]string{})
	if err != nil {
		return err.Error()
	}
	return buf.String()
}