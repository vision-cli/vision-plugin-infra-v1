package placeholders

import (
	"regexp"

	"github.com/barkimedes/go-deepcopy"
	api_v1 "github.com/vision-cli/api/v1"
)

const (
	ArgsCommandIndex = 0
	ArgsNameIndex    = 1
	// include any other arg indexes here
)

var nonAlphaRegex = regexp.MustCompile(`[^a-zA-Z]+`)

func SetupPlaceholders(req api_v1.PluginRequest) (*api_v1.PluginPlaceholders, error) {
	var err error
	p, err := deepcopy.Anything(&req.Placeholders)
	if err != nil {
		return nil, err
	}
	return p.(*api_v1.PluginPlaceholders), nil
}

func clearString(str string) string {
	return nonAlphaRegex.ReplaceAllString(str, "")
}
