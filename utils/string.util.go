package utils

func GenerateRandomString(stringLength int) string {
	letterRunes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	randomString := make([]rune, stringLength)
	for i := range randomString {
		randomString[i] = rune(letterRunes[RandomInt(0, len(letterRunes)-1)])
	}
	return string(randomString)
}
