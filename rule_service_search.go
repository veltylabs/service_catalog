package servicecatalog

import (
	"webtyp.com/crud"
	"webtyp.com/storage"
	"webtyp.com/orm"
)

func (p *Service) Search(a crud.CurrentUser, filter *Service, m *Module) ([]storage.Condition, error) {
	var conds []storage.Condition
	if filter.SpecialtyId != "" {
		conds = append(conds, orm.Eq(Service_.SpecialtyId, filter.SpecialtyId))
	}
	if filter.ActiveOnly {
		conds = append(conds, orm.Eq(Service_.IsActive, true))
	}
	if filter.Sku != "" {
		conds = append(conds, orm.Eq(Service_.Sku, filter.Sku))
	}
	return conds, nil
}
