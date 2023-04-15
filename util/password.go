package util

func ComparePassword(userPassword string, passwordToCompare string) bool {
	if userPassword == passwordToCompare {
		return true
	}

	return false
}
