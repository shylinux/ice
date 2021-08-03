package ice

type Show struct {
	Name string
	Help string
}

type Shower interface {
	Show(show []*Show) []*Show
	ShortDef() string
	FieldDef() string
}
