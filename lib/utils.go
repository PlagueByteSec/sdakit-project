package lib

// Result pool
type Pool map[string]string

// Add new entry to result pool
func (pool Pool) AddEntry(entry string) {
	(pool)[entry] = entry
}

// Verify that the entry does not exist in Pool
func (pool Pool) ContainsEntry(entry string) bool {
	_, exists := (pool)[entry]
	return exists
}
