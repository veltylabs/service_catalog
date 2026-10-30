package servicecatalog

import (
	"webtyp.com/events"
	"webtyp.com/fmt"
	"webtyp.com/model"
	"webtyp.com/orm"
	"webtyp.com/time"
)

// Temas de los eventos de dominio, publicados después de cada cambio exitoso.
const (
	TopicSpecialtyCreated = "catalog.specialty.created"
	TopicSpecialtyUpdated = "catalog.specialty.updated"
	TopicSpecialtyDeleted = "catalog.specialty.deleted"

	TopicServiceCreated          = "catalog.service.created"
	TopicServiceUpdated          = "catalog.service.updated"
	TopicServiceDeactivated      = "catalog.service.deactivated"
	TopicServiceDeleted          = "catalog.service.deleted"
	TopicServiceAgreementCreated = "catalog.service_agreement.created"
	TopicServiceAgreementUpdated = "catalog.service_agreement.updated"
	TopicServiceAgreementDeleted = "catalog.service_agreement.deleted"
)

type Deps struct {
	IDs        model.IDGenerator // requerido — el módulo NUNCA construye un generador
	Publisher  events.Publisher  // opcional — nil desactiva la publicación de eventos
	TenantID   string            // requerido — el tenant de esta instalación; las operaciones nunca lo leen de la petición
	Currencies []string          // monedas que acepta el catálogo (ISO 4217, tres letras mayúsculas), la primera es la de los servicios nuevos; vacío = solo CurrencyCLP
}

type Module struct {
	tenantID   string
	db         *orm.DB
	ids        model.IDGenerator
	pub        events.Publisher
	currencies []string
}

// Nombres de las dependencias que New exige: son dato (código), no texto a traducir;
// por eso son variables y no constantes (langc traduce las constantes que recibe fmt.Err).
var (
	errDepsIDs      = "service_catalog: Deps.IDs"
	errDepsTenantID = "service_catalog: Deps.TenantID"
)

// Ejemplo: module, err := servicecatalog.New(db, deps)
func New(db *orm.DB, deps Deps) (*Module, error) {
	if deps.IDs == nil {
		return nil, fmt.Err(errDepsIDs, "is required")
	}
	if deps.TenantID == "" {
		return nil, fmt.Err(errDepsTenantID, "is required")
	}

	currencies := deps.Currencies
	if len(currencies) == 0 {
		currencies = []string{CurrencyCLP}
	}
	for _, c := range currencies {
		if len(c) != 3 {
			// simplistic check, real check for 3 upper-case ASCII
			return nil, fmt.Err("service_catalog: Deps.Currencies has an invalid code")
		}
		for i := 0; i < 3; i++ {
			if c[i] < 'A' || c[i] > 'Z' {
				return nil, fmt.Err("service_catalog: Deps.Currencies has an invalid code")
			}
		}
	}

	return &Module{tenantID: deps.TenantID, db: db, ids: deps.IDs, pub: deps.Publisher, currencies: currencies}, nil
}

// Métodos de Especialidad

// Ejemplo: spec, err := cat.GetSpecialty("tenant-1", "spec-id")
func (m *Module) GetSpecialty(tenantId, id string) (Specialty, error) {
	var spec Specialty
	qb := m.db.Query(&spec).Where(Specialty_.Id).Eq(id)
	if tenantId != "" {
		qb = qb.Where(Specialty_.TenantId).Eq(tenantId)
	}
	_, err := ReadOneSpecialty(qb, &spec)
	if err != nil {
		if orm.IsNotFound(err) {
			return Specialty{}, ErrSpecialtyNotFound
		}
		return Specialty{}, err
	}
	return spec, nil
}

// Ejemplo: spec, err := cat.GetSpecialtyByPrefix("tenant-1", "MD")
func (m *Module) GetSpecialtyByPrefix(tenantId, prefix string) (Specialty, error) {
	var spec Specialty
	qb := m.db.Query(&spec).Where(Specialty_.Prefix).Eq(prefix).Where(Specialty_.TenantId).Eq(tenantId)
	_, err := ReadOneSpecialty(qb, &spec)
	if err != nil {
		if orm.IsNotFound(err) {
			return Specialty{}, ErrSpecialtyNotFound
		}
		return Specialty{}, err
	}
	return spec, nil
}

// Ejemplo: list, err := cat.ListSpecialties("tenant-1")
func (m *Module) ListSpecialties(tenantId string) ([]Specialty, error) {
	var spec Specialty
	qb := m.db.Query(&spec).Where(Specialty_.TenantId).Eq(tenantId)
	results, err := ReadAllSpecialty(qb)
	if err != nil {
		return nil, err
	}
	list := make([]Specialty, len(results))
	for i, r := range results {
		list[i] = *r
	}
	return list, nil
}

// Ejemplo: spec, err := cat.UpsertSpecialty(Specialty{TenantId: "tenant-1", Prefix: "MD", Name: "Medicina", IsPublished: true})
func (m *Module) UpsertSpecialty(spec Specialty) (Specialty, error) {
	action := model.ActionCreate
	if spec.Id != "" {
		action = model.ActionUpdate
	}
	if spec.Id == "" {
		spec.Id = m.ids.NewID()
	}
	spec.UpdatedAt = time.NowNano()

	if err := spec.Validate(action); err != nil {
		return Specialty{}, err
	}

	// Validar la unicidad del prefijo por inquilino
	existingPrefix, err := m.GetSpecialtyByPrefix(spec.TenantId, spec.Prefix)
	if err == nil {
		if existingPrefix.Id != spec.Id {
			return Specialty{}, ErrSpecialtyPrefixExists
		}
	} else if model.KindOf(err) != model.ErrorNotFound {
		return Specialty{}, err
	}

	if action == model.ActionCreate {
		if err := m.db.Create(&spec); err != nil {
			return Specialty{}, err
		}
		if m.pub != nil {
			m.pub.Publish(events.Event{Topic: TopicSpecialtyCreated, Payload: &spec})
		}
		return spec, nil
	}

	// Verificar que la especialidad existe y pertenece al inquilino
	_, err = m.GetSpecialty(spec.TenantId, spec.Id)
	if err != nil {
		if model.KindOf(err) == model.ErrorNotFound {
			if err := m.db.Create(&spec); err != nil {
				return Specialty{}, err
			}
			if m.pub != nil {
				m.pub.Publish(events.Event{Topic: TopicSpecialtyCreated, Payload: &spec})
			}
			return spec, nil
		}
		return Specialty{}, err
	}

	if err := m.db.Update(&spec, orm.Eq(Specialty_.Id, spec.Id), orm.Eq(Specialty_.TenantId, spec.TenantId)); err != nil {
		return Specialty{}, err
	}
	if m.pub != nil {
		m.pub.Publish(events.Event{Topic: TopicSpecialtyUpdated, Payload: &spec})
	}
	return spec, nil
}

// Ejemplo: err := cat.DeleteSpecialty("tenant-1", "spec-id")
func (m *Module) DeleteSpecialty(tenantId, id string) error {
	spec, err := m.GetSpecialty(tenantId, id)
	if err != nil {
		return err
	}

	// Verificar si algún catalog.service está usando esta especialidad
	var item Service
	qb := m.db.Query(&item).Where(Service_.SpecialtyId).Eq(id)
	if spec.TenantId != "" {
		qb = qb.Where(Service_.TenantId).Eq(spec.TenantId)
	}
	items, err := ReadAllService(qb)
	if err != nil {
		return err
	}
	if len(items) > 0 {
		return ErrSpecialtyInUse
	}

	if spec.TenantId != "" {
		if err := m.db.Delete(&spec, orm.Eq(Specialty_.Id, id), orm.Eq(Specialty_.TenantId, spec.TenantId)); err != nil {
			return err
		}
	} else {
		if err := m.db.Delete(&spec, orm.Eq(Specialty_.Id, id)); err != nil {
			return err
		}
	}
	if m.pub != nil {
		m.pub.Publish(events.Event{Topic: TopicSpecialtyDeleted, Payload: &spec})
	}
	return nil
}

// Métodos de Servicio

// Ejemplo: svc, err := cat.GetService("tenant-1", "svc-id")
func (m *Module) GetService(tenantId, id string) (Service, error) {
	var item Service
	qb := m.db.Query(&item).Where(Service_.Id).Eq(id)
	if tenantId != "" {
		qb = qb.Where(Service_.TenantId).Eq(tenantId)
	}
	_, err := ReadOneService(qb, &item)
	if err != nil {
		if orm.IsNotFound(err) {
			return Service{}, ErrNotFound
		}
		return Service{}, err
	}
	return item, nil
}

// Ejemplo: svc, err := cat.FindBySKU("tenant-1", "MD0001")
func (m *Module) FindBySKU(tenantId, sku string) (Service, error) {
	var item Service
	qb := m.db.Query(&item).Where(Service_.Sku).Eq(sku).Where(Service_.TenantId).Eq(tenantId)
	_, err := ReadOneService(qb, &item)
	if err != nil {
		if orm.IsNotFound(err) {
			return Service{}, ErrNotFound
		}
		return Service{}, err
	}
	return item, nil
}

// Ejemplo: list, err := cat.ListServices("tenant-1", ServiceFilter{ActiveOnly: true})
func (m *Module) ListServices(tenantId string, filter ServiceFilter) ([]Service, error) {
	var item Service
	qb := m.db.Query(&item).Where(Service_.TenantId).Eq(tenantId)
	if filter.SpecialtyId != "" {
		qb = qb.Where(Service_.SpecialtyId).Eq(filter.SpecialtyId)
	}

	if filter.ActiveOnly {
		qb = qb.Where(Service_.IsActive).Eq(true)
	}
	if filter.Limit > 0 {
		qb = qb.Limit(int(filter.Limit))
	}
	if filter.Offset > 0 {
		qb = qb.Offset(int(filter.Offset))
	}
	results, err := ReadAllService(qb)
	if err != nil {
		return nil, err
	}
	items := make([]Service, len(results))
	for i, r := range results {
		items[i] = *r
	}
	return items, nil
}

// Ejemplo: svc, err := cat.CreateService(Service{TenantId: "tenant-1", Sku: "MD0001", Name: "Consulta", Price: 25000, Currency: "CLP", IsActive: true, TaxTreatment: TaxExempt})
func (m *Module) CreateService(item Service) (Service, error) {
	if item.Id == "" {
		item.Id = m.ids.NewID()
	}
	item.UpdatedAt = time.NowNano()

	if err := item.Validate(model.ActionCreate); err != nil {
		return Service{}, err
	}

	validCurrency := false
	for _, c := range m.currencies {
		if item.Currency == c {
			validCurrency = true
			break
		}
	}
	if !validCurrency {
		return Service{}, ErrUnsupportedCurrency
	}

	if item.TaxTreatment != TaxTaxed && item.TaxTreatment != TaxExempt {
		return Service{}, &model.FieldError{Field: "tax_treatment", Err: fmt.Err("invalid tax treatment")}
	}

	// Validar la unicidad del SKU por inquilino
	existing, err := m.FindBySKU(item.TenantId, item.Sku)
	if err == nil {
		if existing.Id != "" {
			return Service{}, ErrAlreadyExists
		}
	} else if model.KindOf(err) != model.ErrorNotFound {
		return Service{}, err
	}

	if err := m.db.Create(&item); err != nil {
		return Service{}, err
	}
	if m.pub != nil {
		m.pub.Publish(events.Event{Topic: TopicServiceCreated, Payload: &item})
	}
	return item, nil
}

// Ejemplo: svc, err := cat.UpdateService(existingSvc)
func (m *Module) UpdateService(item Service) (Service, error) {
	if err := item.Validate(model.ActionUpdate); err != nil {
		return Service{}, err
	}

	validCurrency := false
	for _, c := range m.currencies {
		if item.Currency == c {
			validCurrency = true
			break
		}
	}
	if !validCurrency {
		return Service{}, ErrUnsupportedCurrency
	}

	if item.TaxTreatment != TaxTaxed && item.TaxTreatment != TaxExempt {
		return Service{}, &model.FieldError{Field: "tax_treatment", Err: fmt.Err("invalid tax treatment")}
	}

	// Verificar que el ítem existe y pertenece al inquilino
	existing, err := m.GetService(item.TenantId, item.Id)
	if err != nil {
		return Service{}, err
	}
	if item.TenantId == "" {
		item.TenantId = existing.TenantId
	}

	item.UpdatedAt = time.NowNano()
	if err := m.db.Update(&item, orm.Eq(Service_.Id, item.Id), orm.Eq(Service_.TenantId, item.TenantId)); err != nil {
		return Service{}, err
	}
	if m.pub != nil {
		m.pub.Publish(events.Event{Topic: TopicServiceUpdated, Payload: &item})
	}
	return item, nil
}

// Ejemplo: err := cat.DeactivateService("tenant-1", "svc-id")
func (m *Module) DeactivateService(tenantId, id string) error {
	item, err := m.GetService(tenantId, id)
	if err != nil {
		return err
	}
	item.IsActive = false
	item.UpdatedAt = time.NowNano()
	if err := m.db.Update(&item, orm.Eq(Service_.Id, item.Id), orm.Eq(Service_.TenantId, item.TenantId)); err != nil {
		return err
	}
	if m.pub != nil {
		m.pub.Publish(events.Event{Topic: TopicServiceDeactivated, Payload: &item})
	}
	return nil
}

// Ejemplo: err := cat.DeleteService("tenant-1", "svc-id")
func (m *Module) DeleteService(tenantId, id string) error {
	item, err := m.GetService(tenantId, id)
	if err != nil {
		return err
	}
	if item.TenantId != "" {
		if err := m.db.Delete(&item, orm.Eq(Service_.Id, item.Id), orm.Eq(Service_.TenantId, item.TenantId)); err != nil {
			return err
		}
	} else {
		if err := m.db.Delete(&item, orm.Eq(Service_.Id, item.Id)); err != nil {
			return err
		}
	}
	if m.pub != nil {
		m.pub.Publish(events.Event{Topic: TopicServiceDeleted, Payload: &item})
	}
	return nil
}

// Ejemplo: exists, err := cat.ServiceExists("tenant-1", "svc-id")
func (m *Module) ServiceExists(tenantId, serviceId string) (bool, error) {
	item, err := m.GetService(tenantId, serviceId)
	if err != nil {
		if model.KindOf(err) == model.ErrorNotFound {
			return false, nil
		}
		return false, err
	}
	return item.IsActive, nil
}

// ServiceCard devuelve nombre, precio de lista, moneda y estado de un servicio
// del catálogo. Satisface CatalogReader.ServiceCard de appointment_booking
// estructuralmente. Ausente, de otro tenant o que no es servicio ⇒ ErrNotFound.
// Ejemplo: name, price, currency, active, err := cat.ServiceCard("tenant-1", "svc-id")
func (m *Module) ServiceCard(tenantId, serviceId string) (name string, price float64, currency string, active bool, err error) {
	if tenantId == "" || serviceId == "" {
		return "", 0, "", false, ErrNotFound
	}
	item, e := m.GetService(tenantId, serviceId)
	if e != nil {
		return "", 0, "", false, e
	}

	return item.Name, item.Price, item.Currency, item.IsActive, nil
}

// Ejemplo: list, err := cat.ListServiceAgreements("tenant-1", "svc-id")
func (m *Module) ListServiceAgreements(tenantId, serviceId string) ([]ServiceAgreement, error) {
	var a ServiceAgreement
	qb := m.db.Query(&a).Where(ServiceAgreement_.TenantId).Eq(tenantId)
	if serviceId != "" {
		qb = qb.Where(ServiceAgreement_.ServiceId).Eq(serviceId)
	}
	results, err := ReadAllServiceAgreement(qb)
	if err != nil {
		return nil, err
	}
	items := make([]ServiceAgreement, len(results))
	for i, r := range results {
		items[i] = *r
	}
	return items, nil
}

// Ejemplo: agr, err := cat.GetServiceAgreement("tenant-1", "agr-id")
func (m *Module) GetServiceAgreement(tenantId, id string) (ServiceAgreement, error) {
	var a ServiceAgreement
	qb := m.db.Query(&a).Where(ServiceAgreement_.Id).Eq(id).Where(ServiceAgreement_.TenantId).Eq(tenantId)
	_, err := ReadOneServiceAgreement(qb, &a)
	if err != nil {
		if orm.IsNotFound(err) {
			return ServiceAgreement{}, ErrNotFound
		}
		return ServiceAgreement{}, err
	}
	return a, nil
}

// Ejemplo: agr, err := cat.UpsertServiceAgreement(ServiceAgreement{TenantId: "tenant-1", ServiceId: "svc-id", Insurer: "Isapre", Price: 15000, IsActive: true})
func (m *Module) UpsertServiceAgreement(a ServiceAgreement) (ServiceAgreement, error) {
	action := model.ActionCreate
	if a.Id != "" {
		action = model.ActionUpdate
	}
	if a.Id == "" {
		a.Id = m.ids.NewID()
	}
	a.UpdatedAt = time.NowNano()

	if err := a.Validate(action); err != nil {
		return ServiceAgreement{}, err
	}

	if action == model.ActionCreate {
		if err := m.db.Create(&a); err != nil {
			return ServiceAgreement{}, err
		}
		if m.pub != nil {
			m.pub.Publish(events.Event{Topic: TopicServiceAgreementCreated, Payload: &a})
		}
		return a, nil
	}

	// Verificar que el acuerdo existe y pertenece al inquilino
	_, err := m.GetServiceAgreement(a.TenantId, a.Id)
	if err != nil {
		return ServiceAgreement{}, err
	}

	if err := m.db.Update(&a, orm.Eq(ServiceAgreement_.Id, a.Id), orm.Eq(ServiceAgreement_.TenantId, a.TenantId)); err != nil {
		return ServiceAgreement{}, err
	}
	if m.pub != nil {
		m.pub.Publish(events.Event{Topic: TopicServiceAgreementUpdated, Payload: &a})
	}
	return a, nil
}

// Ejemplo: err := cat.DeleteServiceAgreement("tenant-1", "agr-id")
func (m *Module) DeleteServiceAgreement(tenantId, id string) error {
	a, err := m.GetServiceAgreement(tenantId, id)
	if err != nil {
		return err
	}

	if err := m.db.Delete(&a, orm.Eq(ServiceAgreement_.Id, id), orm.Eq(ServiceAgreement_.TenantId, tenantId)); err != nil {
		return err
	}
	if m.pub != nil {
		m.pub.Publish(events.Event{Topic: TopicServiceAgreementDeleted, Payload: &a})
	}
	return nil
}

// ModelName es la identidad de este módulo: el transporte califica cada op
// debajo como "service_catalog.<name>" en la red, y los dos constructores
// NewXView de view.go pasan esta misma constante como view.Ops.Module para que
// el cliente componga el nombre calificado idéntico.
const ModelName = "service_catalog"

func (m *Module) ModelName() string { return ModelName }
