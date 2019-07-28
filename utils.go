package persistentconn

func tupleListToMap(tl [][]string) map[string]string {
	res := make(map[string]string)
	for _, val := range tl {
		res[val[0]] = val[1]
	}
	return res
}

func contains(list []string, target string) bool {
	for _, mem := range list {
		if mem == target {
			return true
		}
	}
	return false
}
