package bun

type PromiseHandler func(res any) any

type Promise struct {
	parent        *Promise
	current       PromiseHandler
	then          *Promise
	exceptionally *Promise
	done          chan bool
}

func NewPromise(fn PromiseHandler) *Promise {
	return &Promise{
		parent:        nil,
		current:       fn,
		then:          nil,
		exceptionally: nil,
		done:          make(chan bool),
	}
}

func (promise *Promise) Then(fn PromiseHandler) *Promise {
	promise.then = &Promise{
		parent:        promise,
		current:       fn,
		then:          nil,
		exceptionally: nil,
		done:          nil,
	}
	return promise.then
}

func (promise *Promise) Exceptionally(fn PromiseHandler) *Promise {
	promise.exceptionally = &Promise{
		parent:        promise,
		current:       fn,
		then:          nil,
		exceptionally: nil,
	}
	return promise.exceptionally
}

func (promise *Promise) Run() {
	if promise.parent != nil {
		promise.parent.Run()
		return
	}
	promise.run(nil, promise.done)
}

func (promise *Promise) Await() {
	if promise.parent != nil {
		promise.parent.Await()
		return
	}

	promise.Run()

	select {
	case <-promise.done:
		return
	}
}

func (promise *Promise) run(res any, done chan bool) {
	go func() {
		Try(func() {
			result := promise.current(res)
			if promise.then != nil {
				promise.then.run(result, done)
				return
			}
			if promise.exceptionally != nil && promise.exceptionally.then != nil {
				promise.exceptionally.then.run(result, done)
				return
			}
			done <- true
		}).Catch(func(err any) {
			if promise.exceptionally != nil {
				promise.exceptionally.run(err, done)
				return
			}
			done <- true
		}).Run()
	}()
}
