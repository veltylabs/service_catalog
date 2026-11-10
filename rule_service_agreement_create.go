package servicecatalog

import "webtyp.com/crud"

func (p *ServiceAgreement) Create(a crud.CurrentUser, tx crud.Tx, m *Module) (crud.RuleAnswer, error) {
	return crud.ApplyAndPublish, nil
}
