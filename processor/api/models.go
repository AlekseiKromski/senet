package api

type Validation struct {
	Result   bool     `json:"result"`
	Messages []string `json:"messages"`
}

func newValidation() *Validation {
	return &Validation{
		Result:   true,
		Messages: []string{},
	}
}

func (v *Validation) appendMessage(message string) {
	v.Messages = append(v.Messages, message)
}

func (v *Validation) setResult(result bool) {
	v.Result = result
}

type Validator interface {
	validate() *Validation
}

//Models

type register struct {
	Username             string `json:"username"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"password_confirmation"`
}

func (r *register) validate() *Validation {
	v := newValidation()
	if len(r.Username) <= 8 {
		v.appendMessage("username less than 8")
		v.setResult(false)
	}

	if len(r.Password) <= 8 || len(r.PasswordConfirmation) <= 8 {
		v.appendMessage("password less than 8")
		v.setResult(false)
	}

	if r.Password != r.PasswordConfirmation {
		v.appendMessage("passwords are different")
		v.setResult(false)
	}

	return v
}

type login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (l *login) validate() *Validation {
	v := newValidation()
	if len(l.Username) <= 0 {
		v.appendMessage("no user")
		v.setResult(false)
	}

	if len(l.Password) <= 0 {
		v.appendMessage("no password")
		v.setResult(false)
	}

	return v
}

type createChat struct {
	Type  string   `json:"type"`
	Users []string `json:"users"`
}

func (cc *createChat) validate() *Validation {
	v := newValidation()
	if len(cc.Type) <= 0 {
		v.appendMessage("no user")
		v.setResult(false)
	}

	if len(cc.Users) <= 0 {
		v.appendMessage("no password")
		v.setResult(false)
	}

	return v
}
