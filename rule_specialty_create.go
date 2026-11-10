package servicecatalog

import (
	"webtyp.com/crud"
	"webtyp.com/model"
)

func (p *Specialty) Create(a crud.CurrentUser, tx crud.Tx, m *Module) (crud.RuleAnswer, error) {
	conflict := false
	err := tx.ReadAll(func() model.Model { return &Specialty{} }, crud.Query{}, func(row model.Model) {
		if row.(*Specialty).Prefix == p.Prefix {
			conflict = true
		}
	})
	if err != nil {
		return 0, err
	}
	if conflict {
		return 0, ErrSpecialtyPrefixExists
	}
	return crud.ApplyAndPublish, nil
}
