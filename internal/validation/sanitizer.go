package validation

import "github.com/microcosm-cc/bluemonday"

type Sanitizer struct {
	policy *bluemonday.Policy
}

func (s *Sanitizer) Sanitize(input string) string {
	return s.policy.Sanitize(input)
}

func NewSanitizer() *Sanitizer {
	return &Sanitizer{
		policy: bluemonday.UGCPolicy(),
	}
}
