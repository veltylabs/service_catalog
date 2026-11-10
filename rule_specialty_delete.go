package servicecatalog

import (
	"webtyp.com/crud"
	"webtyp.com/model"
	"webtyp.com/storage"
	"webtyp.com/orm"
)

func (p *Specialty) Delete(a crud.CurrentUser, tx crud.Tx, m *Module) (crud.RuleAnswer, error) {
	inUse := false
	err := tx.ReadAll(func() model.Model { return &Service{} }, crud.Query{Conditions: []storage.Condition{orm.Eq(Service_.SpecialtyId, p.Id)}}, func(row model.Model) {
		inUse = true
	})
	if err != nil {
		return 0, err
	}
	if inUse {
		return 0, ErrSpecialtyInUse
	}
	return crud.ApplyAndPublish, nil
}
