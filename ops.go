package servicecatalog

import (
	"webtyp.com/model"
	"webtyp.com/router"
	"webtyp.com/rpc"
)

// Nombres de las operaciones, tal como las ven el transporte y la vista: "<ModelName>.<name>".
const (
	OpListSpecialties = "list_specialties"
	OpGetSpecialty    = "get_specialty"
	OpUpsertSpecialty = "upsert_specialty"
	OpDeleteSpecialty = "delete_specialty"

	OpListServices      = "list_services"
	OpGetService        = "get_service"
	OpFindServiceBySKU  = "find_service_by_sku"
	OpCreateService     = "create_service"
	OpUpdateService     = "update_service"
	OpUpsertService     = "upsert_service"
	OpDeactivateService = "deactivate_service"
	OpDeleteService     = "delete_service"

	OpListCurrencies  = "list_currencies"
	OpListAgreements  = "list_agreements"
	OpUpsertAgreement = "upsert_agreement"
	OpDeleteAgreement = "delete_agreement"
)

// UpsertSpecialtyArgs transporta el lote plural para OpUpsertSpecialty.
// Ejemplo: args := UpsertSpecialtyArgs{Records: []Specialty{spec1, spec2}}
type UpsertSpecialtyArgs struct {
	Records []Specialty
}

func (a *UpsertSpecialtyArgs) ModelName() string     { return "upsert_specialty_args" }
func (a *UpsertSpecialtyArgs) Schema() []model.Field { return nil }
func (a *UpsertSpecialtyArgs) Pointers() []any       { return nil }
func (a *UpsertSpecialtyArgs) IsNil() bool           { return a == nil }
func (a *UpsertSpecialtyArgs) EncodeFields(w model.FieldWriter) {
	arr := w.Array("records", len(a.Records))
	for i := range a.Records {
		arr.Object(&a.Records[i])
	}
}
func (a *UpsertSpecialtyArgs) DecodeFields(r model.FieldReader) {
	if arr, ok := r.Array("records"); ok {
		a.Records = make([]Specialty, arr.Len())
		for i := 0; i < arr.Len(); i++ {
			arr.Object(i, &a.Records[i])
		}
	}
}

// DeleteSpecialtiesArgs transporta los IDs a eliminar para OpDeleteSpecialty.
// Ejemplo: args := DeleteSpecialtiesArgs{Ids: []string{"id1", "id2"}}
type DeleteSpecialtiesArgs struct {
	Ids []string
}

func (a *DeleteSpecialtiesArgs) ModelName() string     { return "delete_specialties_args" }
func (a *DeleteSpecialtiesArgs) Schema() []model.Field { return nil }
func (a *DeleteSpecialtiesArgs) Pointers() []any       { return nil }
func (a *DeleteSpecialtiesArgs) IsNil() bool           { return a == nil }
func (a *DeleteSpecialtiesArgs) EncodeFields(w model.FieldWriter) {
	arr := w.Array("ids", len(a.Ids))
	for _, id := range a.Ids {
		arr.String(id)
	}
}
func (a *DeleteSpecialtiesArgs) DecodeFields(r model.FieldReader) {
	if arr, ok := r.Array("ids"); ok {
		a.Ids = make([]string, arr.Len())
		for i := 0; i < arr.Len(); i++ {
			a.Ids[i] = arr.String(i)
		}
	}
}

// UpsertServiceArgs transporta el lote plural para OpUpsertService.
// Ejemplo: args := UpsertServiceArgs{Records: []Service{svc1, svc2}}
type UpsertServiceArgs struct {
	Records []Service
}

func (a *UpsertServiceArgs) ModelName() string     { return "upsert_item_args" }
func (a *UpsertServiceArgs) Schema() []model.Field { return nil }
func (a *UpsertServiceArgs) Pointers() []any       { return nil }
func (a *UpsertServiceArgs) IsNil() bool           { return a == nil }
func (a *UpsertServiceArgs) EncodeFields(w model.FieldWriter) {
	arr := w.Array("records", len(a.Records))
	for i := range a.Records {
		arr.Object(&a.Records[i])
	}
}
func (a *UpsertServiceArgs) DecodeFields(r model.FieldReader) {
	if arr, ok := r.Array("records"); ok {
		a.Records = make([]Service, arr.Len())
		for i := 0; i < arr.Len(); i++ {
			arr.Object(i, &a.Records[i])
		}
	}
}

// DeleteServicesArgs transporta los IDs a eliminar para OpDeleteService.
// Ejemplo: args := DeleteServicesArgs{Ids: []string{"id1", "id2"}}
type DeleteServicesArgs struct {
	Ids []string
}

func (a *DeleteServicesArgs) ModelName() string     { return "delete_items_args" }
func (a *DeleteServicesArgs) Schema() []model.Field { return nil }
func (a *DeleteServicesArgs) Pointers() []any       { return nil }
func (a *DeleteServicesArgs) IsNil() bool           { return a == nil }
func (a *DeleteServicesArgs) EncodeFields(w model.FieldWriter) {
	arr := w.Array("ids", len(a.Ids))
	for _, id := range a.Ids {
		arr.String(id)
	}
}
func (a *DeleteServicesArgs) DecodeFields(r model.FieldReader) {
	if arr, ok := r.Array("ids"); ok {
		a.Ids = make([]string, arr.Len())
		for i := 0; i < arr.Len(); i++ {
			a.Ids[i] = arr.String(i)
		}
	}
}

// UpdateServiceArgs transporta el parche de campos para OpUpdateService.
// Ejemplo: args := UpdateServiceArgs{Ids: []string{"id1"}, Fields: []string{"price"}, Record: Service{Price: 20000}}
type UpdateServiceArgs struct {
	Ids    []string
	Fields []string
	Record Service
}

func (a *UpdateServiceArgs) ModelName() string     { return "update_item_args" }
func (a *UpdateServiceArgs) Schema() []model.Field { return nil }
func (a *UpdateServiceArgs) Pointers() []any       { return nil }
func (a *UpdateServiceArgs) IsNil() bool           { return a == nil }
func (a *UpdateServiceArgs) EncodeFields(w model.FieldWriter) {
	ids := w.Array("ids", len(a.Ids))
	for _, id := range a.Ids {
		ids.String(id)
	}
	fields := w.Array("fields", len(a.Fields))
	for _, f := range a.Fields {
		fields.String(f)
	}
	w.Object("record", &a.Record)
}
func (a *UpdateServiceArgs) DecodeFields(r model.FieldReader) {
	if arr, ok := r.Array("ids"); ok {
		a.Ids = make([]string, arr.Len())
		for i := 0; i < arr.Len(); i++ {
			a.Ids[i] = arr.String(i)
		}
	}
	if arr, ok := r.Array("fields"); ok {
		a.Fields = make([]string, arr.Len())
		for i := 0; i < arr.Len(); i++ {
			a.Fields[i] = arr.String(i)
		}
	}
	r.Object("record", &a.Record)
}

func (m *Module) MountOperations(reg router.OperationRegistry) {
	reg.Operation(OpListSpecialties, m.opListSpecialties).Requires(model.ResourceOf(&SpecialtyModel), model.Read).Accepts(&ListSpecialtiesArgs{})
	reg.Operation(OpGetSpecialty, m.opGetSpecialty).Requires(model.ResourceOf(&SpecialtyModel), model.Read).Accepts(&GetSpecialtyArgs{})
	reg.Operation(OpUpsertSpecialty, m.opUpsertSpecialty).Requires(model.ResourceOf(&SpecialtyModel), model.Create|model.Update).Accepts(&UpsertSpecialtyArgs{})
	reg.Operation(OpDeleteSpecialty, m.opDeleteSpecialty).Requires(model.ResourceOf(&SpecialtyModel), model.Delete).Accepts(&DeleteSpecialtiesArgs{})

	reg.Operation(OpListServices, m.opListServices).Requires(model.ResourceOf(&ServiceModel), model.Read).Accepts(&ListServicesArgs{})
	reg.Operation(OpGetService, m.opGetService).Requires(model.ResourceOf(&ServiceModel), model.Read).Accepts(&GetServiceArgs{})
	reg.Operation(OpFindServiceBySKU, m.opFindServiceBySKU).Requires(model.ResourceOf(&ServiceModel), model.Read).Accepts(&FindServiceBySKUArgs{})
	reg.Operation(OpCreateService, m.opCreateService).Requires(model.ResourceOf(&ServiceModel), model.Create).Accepts(&Service{})
	reg.Operation(OpUpdateService, m.opUpdateService).Requires(model.ResourceOf(&ServiceModel), model.Update).Accepts(&UpdateServiceArgs{})
	reg.Operation(OpUpsertService, m.opUpsertService).Requires(model.ResourceOf(&ServiceModel), model.Create|model.Update).Accepts(&UpsertServiceArgs{})
	reg.Operation(OpDeactivateService, m.opDeactivateService).Requires(model.ResourceOf(&ServiceModel), model.Update).Accepts(&DeactivateServiceArgs{})
	reg.Operation(OpDeleteService, m.opDeleteService).Requires(model.ResourceOf(&ServiceModel), model.Delete).Accepts(&DeleteServicesArgs{})

	reg.Operation(OpListCurrencies, m.opListCurrencies).Requires(model.ResourceOf(&ServiceModel), model.Read).Accepts(&ListCurrenciesArgs{})

	reg.Operation(OpListAgreements, m.opListServiceAgreements).Requires(model.ResourceOf(&ServiceAgreementModel), model.Read).Accepts(&ListAgreementsArgs{})
	reg.Operation(OpUpsertAgreement, m.opUpsertServiceAgreement).Requires(model.ResourceOf(&ServiceAgreementModel), model.Create|model.Update).Accepts(&ServiceAgreement{})
	reg.Operation(OpDeleteAgreement, m.opDeleteServiceAgreement).Requires(model.ResourceOf(&ServiceAgreementModel), model.Delete).Accepts(&DeleteAgreementArgs{})

	reg.Operation(OpListServiceSupplies, m.opListServiceSupplies).Requires(model.ResourceOf(&ServiceSupplyModel), model.Read).Accepts(&ListServiceSuppliesArgs{})
	reg.Operation(OpUpsertServiceSupply, m.opUpsertServiceSupply).Requires(model.ResourceOf(&ServiceSupplyModel), model.Create|model.Update).Accepts(&ServiceSupply{})
	reg.Operation(OpDeleteServiceSupply, m.opDeleteServiceSupply).Requires(model.ResourceOf(&ServiceSupplyModel), model.Delete).Accepts(&DeleteServiceSupplyArgs{})
}

var _ router.OperationModule = (*Module)(nil)

func (m *Module) opListSpecialties(ctx router.Context) {
	var args ListSpecialtiesArgs
	if err := ctx.Decode(&args); err != nil {
		rpc.WriteError(ctx, ErrCannotRead)
		return
	}
	list, err := m.ListSpecialties(m.tenantID)
	if err != nil {
		rpc.WriteError(ctx, err)
		return
	}
	resList := make(SpecialtyList, len(list))
	for i := range list {
		resList[i] = &list[i]
	}
	if err := ctx.Encode(&resList); err != nil {
		rpc.WriteError(ctx, err)
	}
}

func (m *Module) opGetSpecialty(ctx router.Context) {
	var args GetSpecialtyArgs
	if err := ctx.Decode(&args); err != nil {
		rpc.WriteError(ctx, ErrCannotRead)
		return
	}
	spec, err := m.GetSpecialty(m.tenantID, args.Id)
	if err != nil {
		rpc.WriteError(ctx, err)
		return
	}
	if err := ctx.Encode(&spec); err != nil {
		rpc.WriteError(ctx, err)
	}
}

func (m *Module) opUpsertSpecialty(ctx router.Context) {
	var args UpsertSpecialtyArgs
	if err := ctx.Decode(&args); err != nil {
		rpc.WriteError(ctx, ErrCannotRead)
		return
	}
	if len(args.Records) == 0 {
		rpc.WriteError(ctx, ErrNoRecords)
		return
	}
	var last Specialty
	for _, spec := range args.Records {
		spec.TenantId = m.tenantID
		out, err := m.UpsertSpecialty(spec)
		if err != nil {
			rpc.WriteError(ctx, err)
			return
		}
		last = out
	}
	if err := ctx.Encode(&last); err != nil {
		rpc.WriteError(ctx, err)
	}
}

func (m *Module) opDeleteSpecialty(ctx router.Context) {
	var args DeleteSpecialtiesArgs
	if err := ctx.Decode(&args); err != nil {
		rpc.WriteError(ctx, ErrCannotRead)
		return
	}
	if len(args.Ids) == 0 {
		rpc.WriteError(ctx, ErrNoIds)
		return
	}
	for _, id := range args.Ids {
		if err := m.DeleteSpecialty(m.tenantID, id); err != nil {
			rpc.WriteError(ctx, err)
			return
		}
	}
	ctx.WriteStatus(200)
}

func (m *Module) opListServices(ctx router.Context) {
	var args ListServicesArgs
	if err := ctx.Decode(&args); err != nil {
		rpc.WriteError(ctx, ErrCannotRead)
		return
	}
	filter := ServiceFilter{SpecialtyId: args.SpecialtyId, ActiveOnly: args.ActiveOnly, Limit: args.Limit, Offset: args.Offset}
	items, err := m.ListServices(m.tenantID, filter)
	if err != nil {
		rpc.WriteError(ctx, err)
		return
	}
	list := make(ServiceList, len(items))
	for i := range items {
		list[i] = &items[i]
	}
	if err := ctx.Encode(&list); err != nil {
		rpc.WriteError(ctx, err)
	}
}

func (m *Module) opGetService(ctx router.Context) {
	var args GetServiceArgs
	if err := ctx.Decode(&args); err != nil {
		rpc.WriteError(ctx, ErrCannotRead)
		return
	}
	item, err := m.GetService(m.tenantID, args.Id)
	if err != nil {
		rpc.WriteError(ctx, err)
		return
	}
	if err := ctx.Encode(&item); err != nil {
		rpc.WriteError(ctx, err)
	}
}

func (m *Module) opFindServiceBySKU(ctx router.Context) {
	var args FindServiceBySKUArgs
	if err := ctx.Decode(&args); err != nil {
		rpc.WriteError(ctx, ErrCannotRead)
		return
	}
	item, err := m.FindBySKU(m.tenantID, args.Sku)
	if err != nil {
		rpc.WriteError(ctx, err)
		return
	}
	if err := ctx.Encode(&item); err != nil {
		rpc.WriteError(ctx, err)
	}
}

func (m *Module) opCreateService(ctx router.Context) {
	var item Service
	if err := ctx.Decode(&item); err != nil {
		rpc.WriteError(ctx, ErrCannotRead)
		return
	}
	item.TenantId = m.tenantID
	created, err := m.CreateService(item)
	if err != nil {
		rpc.WriteError(ctx, err)
		return
	}
	if err := ctx.Encode(&created); err != nil {
		rpc.WriteError(ctx, err)
	}
}

func (m *Module) opUpdateService(ctx router.Context) {
	var args UpdateServiceArgs
	if err := ctx.Decode(&args); err != nil {
		rpc.WriteError(ctx, ErrCannotRead)
		return
	}
	if len(args.Ids) == 0 {
		rpc.WriteError(ctx, ErrNoIds)
		return
	}
	for _, id := range args.Ids {
		existing, err := m.GetService(m.tenantID, id)
		if err != nil {
			rpc.WriteError(ctx, err)
			return
		}
		for _, f := range args.Fields {
			switch f {
			case "name":
				existing.Name = args.Record.Name
			case "description":
				existing.Description = args.Record.Description
			case "price":
				existing.Price = args.Record.Price
			case "sku":
				existing.Sku = args.Record.Sku
			case "is_active":
				existing.IsActive = args.Record.IsActive
			case "specialty_id":
				existing.SpecialtyId = args.Record.SpecialtyId
			}
		}
		if _, err := m.UpdateService(existing); err != nil {
			rpc.WriteError(ctx, err)
			return
		}
	}
	ctx.WriteStatus(200)
}

func (m *Module) opUpsertService(ctx router.Context) {
	var args UpsertServiceArgs
	if err := ctx.Decode(&args); err != nil {
		rpc.WriteError(ctx, ErrCannotRead)
		return
	}
	if len(args.Records) == 0 {
		rpc.WriteError(ctx, ErrNoRecords)
		return
	}
	var last Service
	for _, item := range args.Records {
		item.TenantId = m.tenantID
		var out Service
		var err error
		if item.Id == "" {
			out, err = m.CreateService(item)
		} else {
			out, err = m.UpdateService(item)
			if err != nil && model.KindOf(err) == model.ErrorNotFound {
				out, err = m.CreateService(item)
			}
		}
		if err != nil {
			rpc.WriteError(ctx, err)
			return
		}
		last = out
	}
	if err := ctx.Encode(&last); err != nil {
		rpc.WriteError(ctx, err)
	}
}

func (m *Module) opDeactivateService(ctx router.Context) {
	var args DeactivateServiceArgs
	if err := ctx.Decode(&args); err != nil {
		rpc.WriteError(ctx, ErrCannotRead)
		return
	}
	if err := m.DeactivateService(m.tenantID, args.Id); err != nil {
		rpc.WriteError(ctx, err)
		return
	}
	ctx.WriteStatus(200)
}

func (m *Module) opDeleteService(ctx router.Context) {
	var args DeleteServicesArgs
	if err := ctx.Decode(&args); err != nil {
		rpc.WriteError(ctx, ErrCannotRead)
		return
	}
	if len(args.Ids) == 0 {
		rpc.WriteError(ctx, ErrNoIds)
		return
	}
	for _, id := range args.Ids {
		if err := m.DeleteService(m.tenantID, id); err != nil {
			rpc.WriteError(ctx, err)
			return
		}
	}
	ctx.WriteStatus(200)
}

func (m *Module) opListServiceAgreements(ctx router.Context) {
	var args ListAgreementsArgs
	if err := ctx.Decode(&args); err != nil {
		rpc.WriteError(ctx, ErrCannotRead)
		return
	}
	items, err := m.ListServiceAgreements(m.tenantID, args.ServiceId)
	if err != nil {
		rpc.WriteError(ctx, err)
		return
	}
	list := make(ServiceAgreementList, len(items))
	for i := range items {
		list[i] = &items[i]
	}
	if err := ctx.Encode(&list); err != nil {
		rpc.WriteError(ctx, err)
	}
}

func (m *Module) opUpsertServiceAgreement(ctx router.Context) {
	var a ServiceAgreement
	if err := ctx.Decode(&a); err != nil {
		rpc.WriteError(ctx, ErrCannotRead)
		return
	}
	a.TenantId = m.tenantID
	out, err := m.UpsertServiceAgreement(a)
	if err != nil {
		rpc.WriteError(ctx, err)
		return
	}
	if err := ctx.Encode(&out); err != nil {
		rpc.WriteError(ctx, err)
	}
}

func (m *Module) opDeleteServiceAgreement(ctx router.Context) {
	var args DeleteAgreementArgs
	if err := ctx.Decode(&args); err != nil {
		rpc.WriteError(ctx, ErrCannotRead)
		return
	}
	if err := m.DeleteServiceAgreement(m.tenantID, args.Id); err != nil {
		rpc.WriteError(ctx, err)
		return
	}
	ctx.WriteStatus(200)
}

func (m *Module) opListCurrencies(ctx router.Context) {
	var args ListCurrenciesArgs
	if err := ctx.Decode(&args); err != nil {
		rpc.WriteError(ctx, ErrCannotRead)
		return
	}

	list := make(CurrencyOptionList, len(m.currencies))
	for i, c := range m.currencies {
		list[i] = &CurrencyOption{Code: c}
	}

	if err := ctx.Encode(&list); err != nil {
		rpc.WriteError(ctx, err)
	}
}
