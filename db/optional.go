package db

// Optional initialises a variable given a suitable loader function
type Optional struct {
	val   interface{}
	block func() error
}

// Get returns the variable (initialised or not)
func (o *Optional) Get() interface{} {
	return o.val
}

// Nil checks whether the variable is initialised yet
func (o *Optional) Nil() bool {
	return o.val == nil
}

// Block will wait for the loader to finish initialising the variable, and return it after
func (o *Optional) Block() (interface{}, error) {
	err := o.block()
	return o.Get(), err
}

// NewOptional creates a new optional with the given loader func
func NewOptional(f func() (interface{}, error)) *Optional {
	ch := make(chan error, 1)
	optional := &Optional{
		val: nil,
		block: func() error {
			result := <-ch
			ch <- result
			return result
		},
	}

	go func() {
		val, err := f()
		optional.val = val
		ch <- err
	}()

	return optional
}
