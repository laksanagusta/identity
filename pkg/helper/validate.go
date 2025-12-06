package helper

import (
	"strings"

	"github.com/laksanagusta/identity/pkg/errorhelper"
)

const FilterTimeLayout = "2006-01-02T15:04:05.000-07:00"

func ValidateSort(sortableField map[string]struct{}, sort string) error {
	sortList := strings.Split(sort, " ")
	if len(sortList) != 2 {
		return errorhelper.BadRequestMap(map[string][]string{
			"sort": {"invalid"},
		})
	}

	if _, ok := sortableField[sortList[0]]; !ok {
		return errorhelper.BadRequestMap(map[string][]string{
			"sort": {"invalid"},
		})
	}

	if strings.ToLower(sortList[1]) != "asc" && strings.ToLower(sortList[1]) != "desc" {
		return errorhelper.BadRequestMap(map[string][]string{
			"sort": {"invalid"},
		})
	}

	return nil
}

func ValidateSortV2(sortableField map[string]string, sort string) (string, error) {
	sortList := strings.Split(sort, " ")
	if len(sortList) != 2 {
		return "", errorhelper.BadRequestMap(map[string][]string{
			"sort": {"invalid"},
		})
	}

	val, ok := sortableField[sortList[0]]
	if !ok {
		return "", errorhelper.BadRequestMap(map[string][]string{
			"sort": {"invalid"},
		})
	}

	if strings.ToLower(sortList[1]) != "asc" && strings.ToLower(sortList[1]) != "desc" {
		return "", errorhelper.BadRequestMap(map[string][]string{
			"sort": {"invalid"},
		})
	}

	return val, nil
}
