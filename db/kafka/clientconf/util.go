package clientconf

func boolKit(cp string, df bool) bool {
	if cp == "" {
		return df
	}
	return cp == "true"
}
