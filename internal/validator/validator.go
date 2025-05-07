package validator

type Validator struct {
	Errors map[string]string
}

func NewValidator() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (vld *Validator) Valid() bool {
	return len(vld.Errors) == 0
}

func (vld *Validator) AddError(key, value string) {
	if _, ok := vld.Errors[key]; !ok {
		vld.Errors[key] = value
	}
}

func (vld *Validator) CheckError(ok bool, key, value string) {
	if !ok {
		vld.AddError(key, value)
	}
}
