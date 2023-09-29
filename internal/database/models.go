package database

type FizzBuzzRequest struct {
	Int1  int64  `json:"int1"`
	Int2  int64  `json:"int2"`
	Limit int64  `json:"limit"`
	Str1  string `json:"str1"`
	Str2  string `json:"str2"`
}

type MostUsedRequest struct {
	Int1  int64 `json:"int1" db:"value_1"`
	Int2  int64 `json:"int2" db:"value_2"`
	Hints int64 `json:"hints" db:"count"`
}
