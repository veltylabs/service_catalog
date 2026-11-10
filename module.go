package servicecatalog

import (
	"webtyp.com/events"
	"webtyp.com/fmt"
	"webtyp.com/model"
	"webtyp.com/orm"
)

// Temas de los eventos de dominio, publicados después de cada cambio exitoso.
const (

)

type Deps struct {
	IDs        model.IDGenerator // requerido — el módulo NUNCA construye un generador
	Publisher  events.Publisher  // opcional — nil desactiva la publicación de eventos
	TenantID   string            // requerido — el tenant de esta instalación; las operaciones nunca lo leen de la petición
	Currencies []string          // monedas que acepta el catálogo (ISO 4217, tres letras mayúsculas), la primera es la de los servicios nuevos; vacío = solo CurrencyCLP
}

type Module struct {
	tenantID   string
	store      *orm.Store
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

	return &Module{tenantID: deps.TenantID, store: orm.TenantStore(db, deps.TenantID, deps.IDs), db: db, ids: deps.IDs, pub: deps.Publisher, currencies: currencies}, nil
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
	return m.listSpecialties(tenantId, "", "", 0)
}

func (m *Module) listSpecialties(tenantId, search, after string, limit int64) ([]Specialty, error) {
	var spec Specialty
	qb := m.db.Query(&spec).Where(Specialty_.TenantId).Eq(tenantId)
	if search != "" {
		qb = qb.Where(Specialty_.Name).Like("%" + search + "%")
	}
	if after != "" {
		qb = qb.Where(Specialty_.Id).Gt(after)
	}
	qb = qb.OrderBy(Specialty_.Id).Asc()
	if limit > 0 {
		qb = qb.Limit(int(limit))
	}
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
	return spec, nil
}

// Ejemplo: err := cat.DeleteSpecialty("tenant-1", "spec-id")
func (m *Module) DeleteSpecialty(tenantId, id string) error {
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
func (m *Module) ListServices(tenantId string, filter *Service) ([]Service, error) {
	var item Service
	qb := m.db.Query(&item).Where(Service_.TenantId).Eq(tenantId)
	if filter.SearchText != "" {
		qb = qb.Where(Service_.Name).Like("%" + filter.SearchText + "%")
	}
	if filter.After != "" {
		qb = qb.Where(Service_.Id).Gt(filter.After)
	}
	qb = qb.OrderBy(Service_.Id).Asc()

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
	return item, nil
}

// Ejemplo: svc, err := cat.UpdateService(existingSvc)
func (m *Module) UpdateService(item Service) (Service, error) {
	return item, nil
}

// Ejemplo: err := cat.DeactivateService("tenant-1", "svc-id")
func (m *Module) DeactivateService(tenantId, id string) error {
	return nil
}

// Ejemplo: err := cat.DeleteService("tenant-1", "svc-id")
func (m *Module) DeleteService(tenantId, id string) error {
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
	return a, nil
}

// Ejemplo: err := cat.DeleteServiceAgreement("tenant-1", "agr-id")
func (m *Module) DeleteServiceAgreement(tenantId, id string) error {
	return nil
}

// ModelName es la identidad de este módulo: el transporte califica cada op
// debajo como "service_catalog.<name>" en la red, y los dos constructores
// NewXView de view.go pasan esta misma constante como view.Ops.Module para que
// el cliente componga el nombre calificado idéntico.




// checkPrice: un servicio activo necesita un precio mayor que cero; uno inactivo puede esperar su precio en 0.
//
// Ejemplo: err := checkPrice(Service{IsActive: true, Price: 0}) // ErrActiveNeedsPrice
func checkPrice(item Service) error {
	if item.IsActive && item.Price <= 0 {
		return &model.FieldError{Field: "price", Err: ErrActiveNeedsPrice}
	}
	return nil
}
