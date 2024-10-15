package validate

import "regexp"

func ValidateEmail(email string) (bool, error) {
	regexstring := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	is_email, err := regexp.MatchString(regexstring, email)

	if err != nil {
		return is_email, err
	}

	return is_email, nil

}

func ValidateUsername(username string) (bool, error) {
	regexstring := `^[a-zA-Z_][a-zA-Z0-9._%+-]{0,8}$`

	is_valid, err := regexp.MatchString(regexstring, username)

	if err != nil {
		return is_valid, err
	}

	return is_valid, nil
}

func ValidatePassword(passwd string) (bool, error) {
	regexstring := `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d@$#&_]{8,}$`

	is_valid, err := regexp.MatchString(regexstring, passwd)

	if err != nil {
		return is_valid, err
	}

	return is_valid, nil
}
