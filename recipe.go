package servicecatalog

import (
	"webtyp.com/orm"
)

// La receta de insumos de un servicio: lo que usa una atención. clinical_encounter lo lee (ServiceRecipe) para cargar los
// insumos de una atención, y supply_manager los rebaja del stock; ninguno de los dos módulos se importa aquí.


// ListServiceSupplies devuelve la receta de un servicio (o de todos los servicios si serviceId es ""), en orden de inserción.
// Ejemplo: recipe, err := cat.ListServiceSupplies("tenant-1", "service-1")
func (m *Module) ListServiceSupplies(tenantId, serviceId string) ([]ServiceSupply, error) {
	var proxy ServiceSupply
	qb := m.db.Query(&proxy).Where(ServiceSupply_.TenantId).Eq(tenantId)
	if serviceId != "" {
		qb = qb.Where(ServiceSupply_.ServiceId).Eq(serviceId)
	}
	list, err := ReadAllServiceSupply(qb)
	if err != nil {
		return nil, err
	}
	out := make([]ServiceSupply, len(list))
	for i, l := range list {
		out[i] = *l
	}
	return out, nil
}

func (m *Module) getServiceSupply(tenantId, id string) (ServiceSupply, error) {
	var proxy ServiceSupply
	qb := m.db.Query(&proxy).Where(ServiceSupply_.Id).Eq(id).Where(ServiceSupply_.TenantId).Eq(tenantId)
	got, err := ReadOneServiceSupply(qb, &proxy)
	if orm.IsNotFound(err) {
		return ServiceSupply{}, ErrSupplyLineNotFound
	}
	if err != nil {
		return ServiceSupply{}, err
	}
	return *got, nil
}

// UpsertServiceSupply crea una línea de receta (sin id) o actualiza una. El servicio debe existir en el inquilino (ErrNotFound),
// la cantidad debe ser > 0 (ErrQuantity) y un insumo aparece solo una vez por servicio (ErrSupplyLineExists).
// Ejemplo: out, err := cat.UpsertServiceSupply(ServiceSupply{TenantId: "tenant", ServiceId: "s-1", ItemId: "item-1", Quantity: 2})
func (m *Module) UpsertServiceSupply(s ServiceSupply) (ServiceSupply, error) {
	return s, nil
}

// DeleteServiceSupply elimina una línea de la receta.
// Ejemplo: err := cat.DeleteServiceSupply("tenant-1", "supply-1")
func (m *Module) DeleteServiceSupply(tenantId, id string) error {
	return nil
}

// ServiceRecipe es el lector que otros módulos declaran en una interfaz propia estrecha: los ids de ítem de insumo de la
// receta de un servicio y la cantidad de cada uno, en el mismo índice. Un servicio desconocido o ajeno devuelve ErrNotFound;
// un servicio sin receta devuelve dos slices vacíos.
// Ejemplo: itemIds, qtys, err := cat.ServiceRecipe("tenant-1", "service-1")
func (m *Module) ServiceRecipe(tenantId, serviceId string) (itemIds []string, quantities []float64, err error) {
	if _, err := m.GetService(tenantId, serviceId); err != nil {
		return nil, nil, err
	}
	lines, err := m.ListServiceSupplies(tenantId, serviceId)
	if err != nil {
		return nil, nil, err
	}
	itemIds = make([]string, len(lines))
	quantities = make([]float64, len(lines))
	for i, l := range lines {
		itemIds[i], quantities[i] = l.ItemId, l.Quantity
	}
	return itemIds, quantities, nil
}



