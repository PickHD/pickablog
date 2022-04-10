package helper

import (
	"strings"

	"github.com/PickHD/pickablog/model"
)

// GenerateSlug responsible to transform title of blog become slugged
func GenerateSlug(title string) string {
	noSpecStr := model.NoSpecialChar.ReplaceAllString(title," ")
	splitAndLowStr := strings.Split(strings.ToLower(noSpecStr), " ")


	return strings.Join(splitAndLowStr, "-")
}