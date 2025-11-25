package database

// QueryOperator represents a database query operator (e.g., =, !=, >, <).
// It defines the type of comparison or operation to be performed in a query filter.
type QueryOperator string

const (
	Eq          QueryOperator = "="  // Eq represents the Equal operator.
	Neq         QueryOperator = "!"  // Neq represents the Not Equal operator.
	Like        QueryOperator = "~"  // Like represents the Like operator (Regex or Wildcard).
	NotLike     QueryOperator = "!~" // NotLike represents the Not Like operator.
	Gt          QueryOperator = ">"  // Gt represents the Greater Than operator.
	Gte         QueryOperator = ">=" // Gte represents the Greater Than or Equal operator.
	Lt          QueryOperator = "<"  // Lt represents the Less Than operator.
	Lte         QueryOperator = "<=" // Lte represents the Less Than or Equal operator.
	In          QueryOperator = "*"  // In represents the In operator (match one of the values).
	NotIn       QueryOperator = "-"  // NotIn represents the Not In operator.
	InSQ        QueryOperator = "*s" // InSQ represents the In Sub Query operator.
	NotInSQ     QueryOperator = "-s" // NotInSQ represents the Not In Sub Query operator.
	Between     QueryOperator = "#"  // Between represents the Between operator (inclusive).
	Contains    QueryOperator = "@"  // Contains represents the Contains operator (for array fields).
	NotContains QueryOperator = "!@" // NotContains represents the Not Contains operator.
	Empty       QueryOperator = "^"  // Empty represents the Is Empty operator (null or empty).
	True        QueryOperator = "t"  // True represents the Is True operator.
	False       QueryOperator = "f"  // False represents the Is False operator.
	WithFlag    QueryOperator = "&"  // WithFlag represents the With Flag operator (bitwise AND).
	WithNoFlag  QueryOperator = "!&" // WithNoFlag represents the With No Flag operator.
)
