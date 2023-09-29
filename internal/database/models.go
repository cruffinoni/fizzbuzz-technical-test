package database

type Task struct {
	ID   int64  `json:"id"`
	Task string `json:"task"`
	Done bool   `json:"done"`
}

type FizzBuzzRequest struct {
	Int1  int64  `json:"int1"`
	Int2  int64  `json:"int2"`
	Limit int64  `json:"limit"`
	Str1  string `json:"str1"`
	Str2  string `json:"str2"`
}
