package dbErrors

type DbError struct {
	Code  string
	EN    string
	ZH    string
	Class string
}

func (this DbError) GetEN() string {
	return this.EN
}

func (this DbError) GetZH() string {
	return this.ZH
}

func (this DbError) GetClass() string {
	return this.GetClass()
}
