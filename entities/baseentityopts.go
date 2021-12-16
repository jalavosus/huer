package entities

func ApplyOpts(e HueEntity, opts ...BaseEntityOpt) {
	for _, opt := range opts {
		opt.apply(e)
	}
}

type baseEntityOptFunc struct {
	f func(HueEntity)
}

func newBaseEntityOptFunc(f func(HueEntity)) *baseEntityOptFunc {
	return &baseEntityOptFunc{f}
}

func (lo *baseEntityOptFunc) apply(opts HueEntity) {
	lo.f(opts)
}

type BaseEntityOpt interface {
	apply(HueEntity)
}

func EntityName(name string) BaseEntityOpt {
	return newBaseEntityOptFunc(func(e HueEntity) {
		e.SetName(name)
	})
}

func EntityId(id int) BaseEntityOpt {
	return newBaseEntityOptFunc(func(e HueEntity) {
		e.SetId(id)
	})
}

func EntityUid(uid string) BaseEntityOpt {
	return newBaseEntityOptFunc(func(e HueEntity) {
		e.SetUid(uid)
	})
}

func EntityHuer(h Huer) BaseEntityOpt {
	return newBaseEntityOptFunc(func(e HueEntity) {
		e.SetHue(h)
	})
}