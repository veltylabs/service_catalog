package servicecatalog

import "webtyp.com/crud"

func (m *Module) StoreSpecialty() crud.Store {
	return m.store
}

func (m *Module) StoreService() crud.Store {
	return m.store
}

func (m *Module) StoreServiceAgreement() crud.Store {
	return m.store
}

func (m *Module) StoreServiceSupply() crud.Store {
	return m.store
}
