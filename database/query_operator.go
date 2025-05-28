package database

type QueryOperator string

const (
	Eq          QueryOperator = "="
	Neq                       = "!"
	Like                      = "~"
	NotLike                   = "!~"
	Gt                        = ">"
	Gte                       = ">="
	Lt                        = "<"
	Lte                       = "<="
	In                        = "*"
	NotIn                     = "-"
	InSQ                      = "*s"
	NotInSQ                   = "-s"
	Between                   = "#"
	Contains                  = "@"
	NotContains               = "!@"
	Empty                     = "^"
	True                      = "t"
	False                     = "f"
)
