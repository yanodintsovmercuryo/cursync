package filter

import (
	"github.com/yanodintsovmercuryo/cursync/pkg/string_utils"
)

// GetFilePatterns returns file patterns from flag value
func (f *Filter) GetFilePatterns(flagValue string) ([]string, error) {
	if flagValue == "" {
		return []string{}, nil
	}

	patterns := string_utils.SplitTrimFilter(flagValue, ",")
	return patterns, nil
}
