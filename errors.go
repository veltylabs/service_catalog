package servicecatalog

import (
	"webtyp.com/model"
)

const (
	ErrNotFound           model.NotFound = "servicio_no_encontrado"
	ErrSpecialtyNotFound  model.NotFound = "especialidad_no_encontrada"
	ErrSupplyLineNotFound model.NotFound = "receta_no_encontrada"

	ErrAlreadyExists         model.Conflict = "servicio_duplicado"
	ErrSpecialtyPrefixExists model.Conflict = "prefijo_duplicado"
	ErrSupplyLineExists      model.Conflict = "receta_duplicada"

	ErrSpecialtyInUse      model.Invalid = "especialidad en uso"
	ErrQuantity            model.Invalid = "cantidad debe ser mayor a cero"
	ErrUnsupportedCurrency model.Invalid = "moneda no soportada"
	// ErrActiveNeedsPrice: un servicio activo se puede reservar y cobrar, así que necesita precio. Uno sin precio
	// todavía (por ejemplo, recién importado) se guarda inactivo hasta que alguien le ponga el precio.
	ErrActiveNeedsPrice model.Invalid = "servicio activo necesita precio"

	ErrCannotRead model.Invalid = "service_catalog: cannot read the request"
)
