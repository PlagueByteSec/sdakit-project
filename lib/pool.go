package lib

type Pool struct {
	IPv4Pool    []string
	IPv6Pool    []string
	PoolDomains []string
}

func (pool *Pool) PoolInit() {
	pool.IPv4Pool = make([]string, 0)
	pool.IPv6Pool = make([]string, 0)
	pool.PoolDomains = make([]string, 0)
}

func PoolContainsEntry(pool []string, value string) bool {
	for _, entry := range pool {
		if value == entry {
			return true
		}
	}
	return false
}

func poolRemoveDuplicates(pool []string) []string {
	temp := make(map[string]bool)
	revisedPool := make([]string, 0, len(pool))
	for _, value := range pool {
		if !temp[value] {
			revisedPool = append(revisedPool, value)
			temp[value] = true
		}
	}
	return revisedPool
}

func (pool *Pool) PoolCleanup() {
	poolRemoveDuplicates(pool.IPv4Pool)
	poolRemoveDuplicates(pool.IPv6Pool)
	poolRemoveDuplicates(pool.PoolDomains)
}
