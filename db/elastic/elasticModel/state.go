package elasticModel

type ESCount struct {
	Count int64 `json:"count"`
}

type ESState struct {
	All struct {
		Primaries struct {
			Store struct {
				SizeInBytes int64 `json:"size_in_bytes"`
			} `json:"store"`
		} `json:"primaries"`
		Total struct {
			Store struct {
				SizeInBytes int64 `json:"size_in_bytes"`
			} `json:"store"`
		} `json:"total"`
	} `json:"_all"`
}
