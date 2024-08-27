package shared

func PoolInit(pools *PoolBase) {
	pools.PoolIPv4Addresses = make([]string, 0)
	pools.PoolIPv6Addresses = make([]string, 0)
	pools.PoolSubdomains = make([]string, 0)
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

func PoolCleanup(pools *PoolBase) {
	poolRemoveDuplicates(pools.PoolIPv4Addresses)
	poolRemoveDuplicates(pools.PoolIPv6Addresses)
	poolRemoveDuplicates(pools.PoolSubdomains)
}
