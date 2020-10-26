package meander

type Facede interface {
	public() interface{}
}

func public(o interface{}) interface{} {
	//public()が実装されているか
	if p, ok := o.(Facede); ok {
		return p.public()
	}
	return o
}
