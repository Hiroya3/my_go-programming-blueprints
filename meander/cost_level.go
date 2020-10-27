package meander

type Cost int8

const (
	//iotaで整数の連番を生成
	_     Cost = iota
	Cost1      //1
	Cost2      //2
	Cost3      //3
	Cost4      //4
	Cost5      //5
)

var costStrings = map[string]Cost{
	"$":     Cost1,
	"$$":    Cost2,
	"$$$":   Cost3,
	"$$$$":  Cost4,
	"$$$$$": Cost5,
}

func (l Cost) String() string {
	for s, v := range costStrings {
		if l == v {
			return s
		}
	}
	return "不正な値です"
}
