package utils

import (
	"reflect"
	"regexp"
	"strings"
)

func GenerateRandomString(stringLength int) string {
	letterRunes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	randomString := make([]rune, stringLength)
	for i := range randomString {
		randomString[i] = rune(letterRunes[RandomInt(0, len(letterRunes)-1)])
	}
	return string(randomString)
}

func Slugify(text string) string {
	text = strings.TrimSpace(strings.ToLower(text))

	re := regexp.MustCompile(`[^a-z0-9\s-]`)
	text = re.ReplaceAllString(text, "")

	text = strings.ReplaceAll(text, " ", "-")
	re = regexp.MustCompile(`-+`)
	text = re.ReplaceAllString(text, "-")

	return text
}

func GetStructName[T any](dto T) string {
	t := reflect.TypeOf(dto)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return "Not a struct"
	}

	return t.Name()
}
