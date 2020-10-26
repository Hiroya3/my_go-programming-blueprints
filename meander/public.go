package meander

type Facede interface {
	Public() interface{}
}

func Public(o interface{}) interface{} {
	//public()が実装されているかチェック
	if p, ok := o.(Facede); ok {
		//実装されていたらPublicを呼び出し
		//(Publicで返り値の見え方を定義することでユーザーにjsonの内容に変化を与えない)
		return p.Public()
	}
	return o
}
