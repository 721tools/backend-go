package opensealistener

import (
	"fmt"
	"strings"
)

func parseNFTID(id string) (found bool, chain, contract, token string) {
	str := strings.Split(id, "/")

	if len(str) == 3 {
		return true, str[0], str[1], str[2]
	}

	return false, "", "", ""
}

func format(contract, token string, behavior interface{}) string {
	return fmt.Sprintf("%s/%s/%d", contract, token, behavior)
}

func to(id string, behavior interface{}) string {
	found, _, contract, token := parseNFTID(id)
	if found {
		return format(contract, token, behavior)
	}
	return ""
}
