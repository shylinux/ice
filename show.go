package ice

type Show struct {
	Name string
	Help string
}

type Shower interface {
	Show(key string, show []*Show) []*Show
	ShortDef() string
	FieldDef() string
}
