package constants

import "errors"

var (
	EmptyRequiredVar = errors.New("empty required variable in category")
	NoNewsFoundError = errors.New("no news found")

	ErrorPingElastic = errors.New("error pinging elasticsearch")
)
