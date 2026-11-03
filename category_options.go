package servicecatalog

import "webtyp.com/fmt"

// CategoryOptions entrega las especialidades del catálogo como opciones
// {id, nombre}: es lo que room_layout pide como categorías de un box
// (roomlayout.CategoryReader), sin que ninguno importe al otro.
//
// Ejemplo: opts, err := catalog.CategoryOptions(tenantID)
func (m *Module) CategoryOptions(tenantID string) ([]fmt.KeyValue, error) {
	list, err := m.ListSpecialties(tenantID)
	if err != nil {
		return nil, err
	}
	out := make([]fmt.KeyValue, 0, len(list))
	for _, sp := range list {
		out = append(out, fmt.KeyValue{Key: sp.Id, Value: sp.Name})
	}
	return out, nil
}
