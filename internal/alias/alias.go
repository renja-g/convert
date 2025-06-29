package alias

var aliases = map[string]string{
	".jpg": ".jpeg",
}

func Resolve(s string) string {
	if resolved, ok := aliases[s]; ok {
		return resolved
	}
	return s
}
