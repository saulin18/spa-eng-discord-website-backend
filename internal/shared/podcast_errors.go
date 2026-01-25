package shared

import "errors"

var (
	PodcastNotFound = errors.New("podcast not found")
	LinkReportNotFound = errors.New("link report not found")
)