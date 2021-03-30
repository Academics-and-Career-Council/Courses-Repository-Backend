package models

import (
	"net/http"
	"net/url"
)

type AuthenticationSession struct {
	Subject      string                 `json:"subject"`
	Extra        map[string]interface{} `json:"extra"`
	Header       http.Header            `json:"header"`
	MatchContext MatchContext           `json:"match_context"`
}

type MatchContext struct {
	RegexpCaptureGroups []string `json:"regexp_capture_groups"`
	URL                 *url.URL `json:"url"`
}
