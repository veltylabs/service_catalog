package servicecatalog

import "webtyp.com/crud"

func (next *ServiceAgreement) Update(a crud.CurrentUser, tx crud.Tx, old *ServiceAgreement, fields []string, m *Module) (crud.RuleAnswer, error) {
	return crud.ApplyAndPublish, nil
}
