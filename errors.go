package servicecatalog

import (
	"webtyp.com/model"
)

const (
	ErrNotFound           model.NotFound = "service not found"
	ErrSpecialtyNotFound  model.NotFound = "specialty not found"
	ErrSupplyLineNotFound model.NotFound = "service_catalog: recipe line not found"

	ErrAlreadyExists         model.Conflict = "service already exists"
	ErrSpecialtyPrefixExists model.Conflict = "specialty prefix already exists"
	ErrSupplyLineExists      model.Conflict = "service_catalog: the recipe already has this supply"

	ErrSpecialtyInUse      model.Invalid = "specialty in use"
	ErrNoRecords           model.Invalid = "service_catalog: no records to save"
	ErrNoIds               model.Invalid = "service_catalog: no ids given"
	ErrQuantity            model.Invalid = "service_catalog: the quantity must be greater than zero"
	ErrUnsupportedCurrency model.Invalid = "service_catalog: currency not accepted by this catalog"
	// ErrActiveNeedsPrice: un servicio activo se puede reservar y cobrar, así que necesita precio. Uno sin precio
	// todavía (por ejemplo, recién importado) se guarda inactivo hasta que alguien le ponga el precio.
	ErrActiveNeedsPrice model.Invalid = "an active service needs a price — set the price, or leave the service inactive until it has one"

	ErrCannotRead model.Invalid = "service_catalog: cannot read the request"
)
