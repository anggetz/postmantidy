package core

import "strings"

var registeredUrl string

func IsUniqueUrl(item ItemStructure) bool {
	if strings.Contains(registeredUrl, item.Request.Url.Raw) {
		return false
	}

	registeredUrl += "[DELIM]" + item.Request.Url.Raw

	return true
}
