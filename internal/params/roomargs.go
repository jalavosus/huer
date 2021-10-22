package params

type RoomArgs struct {
	name *string
	id   *int
}

func (a RoomArgs) HasName() bool {
	return a.name != nil && *a.name != ""
}

func (a RoomArgs) Name() string {
	if a.HasName() {
		return *a.name
	}

	return ""
}

func (a RoomArgs) HasID() bool {
	return a.id != nil
}

func (a RoomArgs) ID() int {
	if a.HasID() {
		id := a.id

		return *id
	}

	return defaultNilValue
}

func NewRoomArgs(opts ...Param) *RoomArgs {
	args := new(RoomArgs)

	for _, opt := range opts {
		val := opt.Value()
		if val == opt.NilValue() {
			continue
		}

		switch opt.ParamName() {
		case paramName:
			switch v := val.(type) {
			case *string:
				args.name = v
			case string:
				args.name = &v
			}
		case paramId:
			switch v := val.(type) {
			case *int:
				args.id = v
			case int:
				args.id = &v
			}
		default:
			continue
		}
	}

	return args
}