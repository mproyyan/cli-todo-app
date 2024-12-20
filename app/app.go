package app

import "time"

type List struct {
	Description   string
	CreatedAt     time.Time
	FormattedTime string
	Done          bool
}

type Todo struct {
	List []List
}
