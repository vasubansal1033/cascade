package cascade

// MergeSSTables merges two sets of SSTables into a single sorted SSTable at outPath.
// The first slice (newer) takes priority over the second (older) on key conflicts.
// Tombstones in newer that shadow a live entry in older are dropped from the output.
func MergeSSTables(newer, older []*SSTable, outPath string, counter *IOCounter) (*SSTable, error) {
	return nil, nil
}
