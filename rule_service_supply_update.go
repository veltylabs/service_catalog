package servicecatalog

import (
	"webtyp.com/crud"
	"webtyp.com/model"
	"webtyp.com/storage"
	"webtyp.com/orm"
)

func (next *ServiceSupply) Update(a crud.CurrentUser, tx crud.Tx, old *ServiceSupply, fields []string, m *Module) (crud.RuleAnswer, error) {
	if next.Quantity <= 0 {
		return 0, ErrQuantity
	}
	conflict := false
	err := tx.ReadAll(func() model.Model { return &ServiceSupply{} }, crud.Query{Conditions: []storage.Condition{orm.Eq(ServiceSupply_.ServiceId, next.ServiceId)}}, func(row model.Model) {
		if row.(*ServiceSupply).ItemId == next.ItemId && row.(*ServiceSupply).Id != next.Id {
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
