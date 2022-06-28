// Copyright 2022. Motty Cohen
//
// Database Query operator
//
package database

type queryOperator string

const (
	Eq       queryOperator = "="
	Neq                    = "!"
	Like                   = "~"
	Gt                     = ">"
	Gte                    = ">="
	Lt                     = "<"
	Lte                    = "<="
	In                     = "*"
	NotIn                  = "-"
	Between                = "#"
	Contains               = "@"
)
