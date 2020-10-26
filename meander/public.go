package meander

type Facede interface {
	Public() interface{}
}

func Public(o interface{}) interface{} {
	//public()が実装されているか
	if p, ok := o.(Facede); ok {
		return p.Public()
	}
	return o
}
