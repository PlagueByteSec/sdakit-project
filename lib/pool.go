package lib

var (
	IPv4Pool    = make([]string, 0)
	IPv6Pool    = make([]string, 0)
	PoolDomains = make([]string, 0)
)

func PoolContainsEntry(pool []string, value string) bool {
	for _, entry := range pool {
		if value == entry {
			return true
		}
	}
	return false
}
