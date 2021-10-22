package params

const defaultNilValue int = -99999

const (
	paramName string = "Name"
	paramId   string = "ID"
)

type Param interface {
	NilValue() int
	Value() interface{}
	ParamName() string
}

type param struct {
	n string
}

func (p param) ParamName() string {
	return p.n
}

func (p param) NilValue() int {
	return defaultNilValue
}

type nameParam struct {
	param
	name *string
}

func NameParam(name string) Param {
	return &nameParam{param{paramName}, &name}
}

func (p *nameParam) Value() interface{} {
	if p.name != nil {
		return *p.name
	}

	return p.NilValue()
}

type idParam struct {
	param
	id *int
}

func IDParam(id int) Param {
	return &idParam{param{paramId}, &id}
}

func (p *idParam) Value() interface{} {
	if p.id != nil {
		return p.id
	}

	return p.NilValue()
}