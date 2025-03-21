package apiserver

import (
	"dialogue/internal/validator"
)

func (u *userForm) validateUserForm() *userForm {
	if u.FieldErrors == nil {
		u.FieldErrors = make(map[string]string)
	}

	//Basic validation checks.
	u.CheckField(validator.NotBlank(u.Nickname), "nickname", "This field cannot be blank")
	u.CheckField(validator.NotBlank(u.Email), "email", "This field cannot be blank")
	u.CheckField(validator.Matches(u.Email, validator.EmailRX), "email", "This field must be a valid email address")
	u.CheckField(validator.NotBlank(u.Password), "password", "This field cannot be blank")
	u.CheckField(validator.MinChars(u.Password, 8), "password", "This field must be at least 8 characters long")
	return u
}

func (u *userLoginForm) validateUserForm() *userLoginForm {
	//Basic validation checks.
	u.CheckField(validator.NotBlank(u.Email), "email", "This field cannot be blank")
	u.CheckField(validator.Matches(u.Email, validator.EmailRX), "email", "This field must be a valid email address")
	u.CheckField(validator.NotBlank(u.Password), "password", "This field cannot be blank")
	return u
}

func (p *accountPasswordUpdateForm) validatePasswordUpdateForm() *accountPasswordUpdateForm {
	//Basic validations check.
	p.CheckField(validator.NotBlank(p.CurrentPassword), "currentPassword", "This field cannot be blank")
	p.CheckField(validator.NotBlank(p.NewPassword), "newPassword", "This field cannot be blank")
	p.CheckField(validator.MinChars(p.NewPassword, 8), "newPassword", "This field must be at least 8 characters long")
	p.CheckField(validator.NotBlank(p.NewPasswordConfirmation), "newPasswordConfirmation", "This field cannot be blank")
	p.CheckField(p.NewPassword == p.NewPasswordConfirmation, "newPasswordConfirmation", "Passwords do not match")
	return p
}

func (s *storyForm) validateStoryForm() *storyForm {
	s.CheckField(validator.NotBlank(s.Title), "title", "This field cannot be blank")
	s.CheckField(validator.NotBlank(s.Content), "content", "This field cannot be blank")
	return s
}
