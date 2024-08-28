package shared

func PoolInit(pools *PoolBase) {
	pools.PoolIPv4Addresses = make([]string, 0)
	pools.PoolIPv6Addresses = make([]string, 0)
	pools.PoolSubdomains = make([]string, 0)
}

func PoolRemoveDuplicates(pool []string) []string {
	temp := make(map[string]bool)
	revisedPool := make([]string, 0, len(pool))
	/*for _, value := range pool {
		if !temp[value] {
			revisedPool = append(revisedPool, value)
			temp[value] = true
		}
	}*/
	for idx := 0; idx < len(pool); idx++ {
		value := pool[idx]
		if !temp[value] {
			revisedPool = append(revisedPool, value)
			temp[value] = true
		}
	}
	return revisedPool
}

func PoolCleanup(pools *PoolBase) {
	PoolRemoveDuplicates(pools.PoolIPv4Addresses)
	PoolRemoveDuplicates(pools.PoolIPv6Addresses)
	PoolRemoveDuplicates(pools.PoolSubdomains)
}
