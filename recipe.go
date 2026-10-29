package servicecatalog

import (
	"webtyp.com/model"
	"webtyp.com/orm"
	"webtyp.com/router"
	"webtyp.com/rpc"
	"webtyp.com/time"
)

// La receta de insumos de un servicio: lo que usa una atención. clinical_encounter lo lee (ServiceRecipe) para cargar los
// insumos de una atención, y supply_manager los rebaja del stock; ninguno de los dos módulos se importa aquí.
const (
	OpListServiceSupplies = "list_service_supplies"
	OpUpsertServiceSupply = "upsert_service_supply"
	OpDeleteServiceSupply = "delete_service_supply"
)

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
	action := model.ActionCreate
	var existing ServiceSupply
	if s.Id != "" {
		var err error
		if existing, err = m.getServiceSupply(s.TenantId, s.Id); err != nil {
			return ServiceSupply{}, err
		}
		action = model.ActionUpdate
	} else {
		s.Id = m.ids.NewID()
	}
	s.UpdatedAt = time.NowNano()
	if err := s.Validate(action); err != nil {
		return ServiceSupply{}, err
	}
	if s.Quantity <= 0 {
		return ServiceSupply{}, ErrQuantity
	}
	if _, err := m.GetService(s.TenantId, s.ServiceId); err != nil {
		return ServiceSupply{}, err
	}
	if action == model.ActionCreate || existing.ItemId != s.ItemId || existing.ServiceId != s.ServiceId {
		lines, err := m.ListServiceSupplies(s.TenantId, s.ServiceId)
		if err != nil {
			return ServiceSupply{}, err
		}
		for _, l := range lines {
			if l.ItemId == s.ItemId && l.Id != s.Id {
				return ServiceSupply{}, ErrSupplyLineExists
			}
		}
	}
	if action == model.ActionCreate {
		if err := m.db.Create(&s); err != nil {
			return ServiceSupply{}, err
		}
		return s, nil
	}
	if err := m.db.Update(&s, orm.Eq(ServiceSupply_.Id, s.Id), orm.Eq(ServiceSupply_.TenantId, s.TenantId)); err != nil {
		return ServiceSupply{}, err
	}
	return s, nil
}

// DeleteServiceSupply elimina una línea de la receta.
// Ejemplo: err := cat.DeleteServiceSupply("tenant-1", "supply-1")
func (m *Module) DeleteServiceSupply(tenantId, id string) error {
	s, err := m.getServiceSupply(tenantId, id)
	if err != nil {
		return err
	}
	return m.db.Delete(&s, orm.Eq(ServiceSupply_.Id, s.Id), orm.Eq(ServiceSupply_.TenantId, s.TenantId))
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

func (m *Module) opListServiceSupplies(ctx router.Context) {
	var args ListServiceSuppliesArgs
	if err := ctx.Decode(&args); err != nil {
		rpc.WriteError(ctx, ErrCannotRead)
		return
	}
	items, err := m.ListServiceSupplies(m.tenantID, args.ServiceId)
	if err != nil {
		rpc.WriteError(ctx, err)
		return
	}
	list := make(ServiceSupplyList, len(items))
	for i := range items {
		list[i] = &items[i]
	}
	if err := ctx.Encode(&list); err != nil {
		rpc.WriteError(ctx, err)
	}
}

func (m *Module) opUpsertServiceSupply(ctx router.Context) {
	var s ServiceSupply
	if err := ctx.Decode(&s); err != nil {
		rpc.WriteError(ctx, ErrCannotRead)
		return
	}
	s.TenantId = m.tenantID
	out, err := m.UpsertServiceSupply(s)
	if err != nil {
		rpc.WriteError(ctx, err)
		return
	}
	if err := ctx.Encode(&out); err != nil {
		rpc.WriteError(ctx, err)
	}
}

func (m *Module) opDeleteServiceSupply(ctx router.Context) {
	var args DeleteServiceSupplyArgs
	if err := ctx.Decode(&args); err != nil {
		rpc.WriteError(ctx, ErrCannotRead)
		return
	}
	if err := m.DeleteServiceSupply(m.tenantID, args.Id); err != nil {
		rpc.WriteError(ctx, err)
		return
	}
	ctx.WriteStatus(200)
}
