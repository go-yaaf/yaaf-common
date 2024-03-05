package database

type QueryOperator string

const (
	Eq       QueryOperator = "="
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
