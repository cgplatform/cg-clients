package exceptions

import "errors"

var (
	INVALID_TOKEN             = errors.New("invalid_token")
	USER_NOT_EXISTS           = errors.New("user_not_exists")
	USER_NOT_AUTHORIZED       = errors.New("not_authorized")
	USER_NOT_VERIFIED         = errors.New("user_not_verified")
	INVALID_EMAIL_OR_PASSWORD = errors.New("email_or_password_invalid")
	INVALID_EMAIL             = errors.New("repeat_email")
	WRONG_PASSWORD            = errors.New("wrong_password")
)
