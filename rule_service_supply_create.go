package servicecatalog

import (
	"webtyp.com/crud"
	"webtyp.com/model"
	"webtyp.com/storage"
	"webtyp.com/orm"
)

func (p *ServiceSupply) Create(a crud.CurrentUser, tx crud.Tx, m *Module) (crud.RuleAnswer, error) {
	if p.Quantity <= 0 {
		return 0, ErrQuantity
	}
	conflict := false
	err := tx.ReadAll(func() model.Model { return &ServiceSupply{} }, crud.Query{Conditions: []storage.Condition{orm.Eq(ServiceSupply_.ServiceId, p.ServiceId)}}, func(row model.Model) {
		if row.(*ServiceSupply).ItemId == p.ItemId {
			conflict = true
		}
	})
	if err != nil {
		return 0, err
	}
	if conflict {
		return 0, ErrSupplyLineExists
	}
	return crud.ApplyAndPublish, nil
}
