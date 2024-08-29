package shared

import "github.com/fhAnso/Sentinel/v1/pkg"

func poolRemoveDuplicates(pool []string) []string {
	temp := make(map[string]bool)
	revisedPool := make([]string, 0, len(pool))
	for idx := 0; idx < len(pool); idx++ {
		value := pool[idx]
		if !temp[value] {
			revisedPool = append(revisedPool, value)
			temp[value] = true
		}
	}
	return revisedPool
}

func PoolAppendValue(subdomain string, pool *[]string) {
	if !pkg.IsInSlice(subdomain, *pool) {
		*pool = append(*pool, subdomain)
	}
}

func PoolInit(pools *PoolBase) {
	// CORE
	pools.PoolIPv4Addresses = make([]string, 0)
	pools.PoolIPv6Addresses = make([]string, 0)
	pools.PoolSubdomains = make([]string, 0)
	// SUMMARY
	pools.PoolMailSubdomains = make([]string, 0)
	pools.PoolApiSubdomains = make([]string, 0)
	pools.PoolLoginSubdomains = make([]string, 0)
	pools.PoolCorsSubdomains = make([]string, 0)
	pools.PoolHttpSuccessSubdomains = make([]string, 0)
}

func PoolsCleanupCore(pools *PoolBase) {
	poolRemoveDuplicates(pools.PoolIPv4Addresses)
	poolRemoveDuplicates(pools.PoolIPv6Addresses)
	poolRemoveDuplicates(pools.PoolSubdomains)
}

func PoolsCleanupSummary(pools *PoolBase) {
	poolRemoveDuplicates(pools.PoolMailSubdomains)
	poolRemoveDuplicates(pools.PoolApiSubdomains)
	poolRemoveDuplicates(pools.PoolLoginSubdomains)
	poolRemoveDuplicates(pools.PoolCorsSubdomains)
	poolRemoveDuplicates(pools.PoolHttpSuccessSubdomains)
}
