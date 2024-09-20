package datapools

import (
	"sort"
)

type PoolAction int

const (
	PoolAppend PoolAction = iota
	PoolCheck
	PoolReset
)

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

func PoolsCleanupCore(pools *PoolBase) {
	// General
	pools.PoolIPv4Addresses = poolRemoveDuplicates(pools.PoolIPv4Addresses)
	pools.PoolIPv6Addresses = poolRemoveDuplicates(pools.PoolIPv6Addresses)
	pools.PoolSubdomains = poolRemoveDuplicates(pools.PoolSubdomains)
	pools.PoolHttpSuccessSubdomains = poolRemoveDuplicates(pools.PoolHttpSuccessSubdomains)
}

func PoolsCleanupSummary(pools *PoolBase) {
	// Purpose/Service detection
	pools.PoolMailSubdomains = poolRemoveDuplicates(pools.PoolMailSubdomains)
	pools.PoolApiSubdomains = poolRemoveDuplicates(pools.PoolApiSubdomains)
	pools.PoolLoginSubdomains = poolRemoveDuplicates(pools.PoolLoginSubdomains)
	pools.PoolCmsSubdomains = poolRemoveDuplicates(pools.PoolCmsSubdomains)
	// Security testing
	pools.PoolCorsSubdomains = poolRemoveDuplicates(pools.PoolCorsSubdomains)
	pools.PoolCookieInjection = poolRemoveDuplicates(pools.PoolCookieInjection)
	pools.PoolHeaderInjection = poolRemoveDuplicates(pools.PoolHeaderInjection)
	pools.PoolRequestSmuggling = poolRemoveDuplicates(pools.PoolRequestSmuggling)
}

func ManagePool(action PoolAction, entry string, pool *[]string) bool {
	switch action {
	case PoolAppend: // Append value to the pool if it doesn't already exist
		idx := sort.SearchStrings(*pool, entry)
		if idx >= len(*pool) || (*pool)[idx] != entry {
			*pool = append(*pool, entry)
			sort.Strings(*pool)
		}
		return true
	case PoolCheck: // Check if the entry exists in pool
		idx := sort.SearchStrings(*pool, entry)
		return idx < len(*pool) && (*pool)[idx] == entry
	case PoolReset:
		if len(*pool) >= 1 && (*pool)[0] == "" {
			*pool = []string{}
		}
		return true
	default:
		return false
	}
}
