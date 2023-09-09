package bun

type Effect struct {
	try     func()
	catch   *func(err any)
	finally *func()
}

func Try(fn func()) *Effect {
	return &Effect{
		try:   fn,
		catch: nil,
	}
}

func (effect *Effect) Catch(fn func(err any)) *Effect {
	effect.catch = &fn
	return effect
}

func (effect *Effect) Finally(fn func()) *Effect {
	effect.finally = &fn
	return effect
}

func (effect *Effect) Run() {
	defer func() {
		defer func() {
			if effect.finally != nil {
				(*effect.finally)()
			}
		}()
		if r := recover(); r != nil {
			if effect.catch != nil {
				(*effect.catch)(r)
				return
			}
			panic(r)
		}
	}()
	effect.try()
}
