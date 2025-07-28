package contract

type Module interface {
}

type ModuleName string

type Nameable interface {
	Name() ModuleName
}
