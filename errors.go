package servicecatalog

import (
	"webtyp.com/orm"
	"webtyp.com/router"
)

type domainError string

func (e domainError) Error() string { return string(e) }

const (
	ErrNotFound              domainError = "service not found"
	ErrAlreadyExists         domainError = "service already exists"
	ErrSpecialtyNotFound     domainError = "specialty not found"
	ErrSpecialtyInUse        domainError = "specialty in use"
	ErrSpecialtyPrefixExists domainError = "specialty prefix already exists"
	ErrNoRecords             domainError = "service_catalog: no records to save"
	ErrNoIds                 domainError = "service_catalog: no ids given"
	ErrSupplyLineExists      domainError = "service_catalog: the recipe already has this supply"
	ErrSupplyLineNotFound    domainError = "service_catalog: recipe line not found"
	ErrQuantity              domainError = "service_catalog: the quantity must be greater than zero"
	ErrUnsupportedCurrency   domainError = "service_catalog: currency not accepted by this catalog"
)

// decodeError es un cuerpo de solicitud que no se pudo leer: siempre 400.
type decodeError struct{ err error }

func (d decodeError) Error() string { return msgDecode + d.err.Error() }

const msgDecode = "service_catalog: cannot read the request: "

// writeError responde a err con su estado y, siempre, su texto — lo único que ve un cliente
// HTTP. 400 = ilegible o entrada inválida · 404 = no encontrado · 409 = conflicto · 500 = cualquier otra cosa.
func writeError(ctx router.Context, err error) {
	ctx.WriteStatus(statusOf(err))
	ctx.Write([]byte(err.Error()))
}

func statusOf(err error) int {
	switch e := err.(type) {
	case decodeError, ValidationError:
		return 400
	case domainError:
		switch e {
		case ErrNotFound, ErrSpecialtyNotFound, ErrSupplyLineNotFound:
			return 404
		case ErrAlreadyExists, ErrSpecialtyPrefixExists, ErrSupplyLineExists:
			return 409
		case ErrSpecialtyInUse, ErrNoRecords, ErrNoIds, ErrQuantity, ErrUnsupportedCurrency:
			return 400
		}
	}
	if orm.IsNotFound(err) {
		return 404
	}
	return 500
}

// ValidationError envuelve un fallo de validación del modelo: 400.
type ValidationError struct {
	Err error
}

func (v ValidationError) Error() string {
	return v.Err.Error()
}
