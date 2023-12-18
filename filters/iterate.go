package filters

func IterateSets[T any](devices []T, size int, callback func([]T)) {
	if size <= 0 {
		return
	}

	if size > len(devices) {
		return
	}

	// The logic below is a simple unrolling of the recursive loops:
	//
	// n := len(devices)
	// for i := 0; i < n; i++
	//     for j := i+1; j < n; j++
	//         for k := j+1; k < n; k++
	//             ...
	//             for z := y+1; z < n; z++
	//                 callback({devices[i], devices[j], devices[k], ..., devices[z]})
	//
	// Where 'size' represents how many logical 'for' loops there are, 'level'
	// represents how many 'for' loops deep we are, 'indices' holds the loop
	// index at each level, and 'set' builds out the list of devices to pass to
	// the callback each time the bottom most level is reached.
	level := 0
	indices := make([]int, size)
	set := make([]T, size)

	for {
		if indices[level] == len(devices) {
			if level == 0 {
				break
			}

			level--
			indices[level]++
			continue
		}

		set[level] = devices[indices[level]]

		if level < (size - 1) {
			level++
			indices[level] = indices[level-1] + 1
			continue
		}

		callback(set)
		indices[level]++
	}
}

func SetCountPadding[T comparable](set []T) int {
	count := 0
	var noop T

	for i := range set {
		if set[i] == noop {
			count++
		}
	}

	return count
}

func SetContains[T comparable](set []T, d T) bool {
	for i := range set {
		if set[i] == d {
			return true
		}
	}
	return false
}

func SetCopyAndAddPadding[T comparable](set []T, size int) []T {
	if size <= 0 {
		return []T{}
	}
	var noop T

	sets := append([]T{}, set...)
	for len(sets)%size != 0 {
		sets = append(sets, noop)
	}

	return sets
}

func IteratePartitions[T comparable](devices []T, size int, callback func([][]T)) {
	if size <= 0 {
		return
	}

	if size > len(devices) {
		return
	}

	// Optimize for the case when size == 1.
	if size == 1 {
		for _, device := range devices {
			callback([][]T{[]T{device}})
		}
		return
	}

	devices = SetCopyAndAddPadding(devices, size)
	padding := SetCountPadding(devices)

	// We wrap the recursive call to make use of an 'accum' variable to
	// build out each partition as the recursion progresses.
	var iterate func(devices []T, size int, accum [][]T)
	iterate = func(devices []T, size int, accum [][]T) {
		// Padding should ensure that his never happens.
		if size > len(devices) {
			panic("Internal error in best effort allocation policy")
		}

		// Base case once we've reached 'size' number of devices.
		if size == len(devices) {
			callback(append(accum, devices))
			return
		}

		// For all other sizes and device lengths ...
		//
		// The code below is optimized to avoid considering duplicate
		// partitions, e.g. [[0,1],[2,3]] and [[2,3],[0,1]].
		//
		// This ensures that the _first_ device index of each set in a
		// partition is in increaing order, e.g. [[0...], [4...], [7...]] and
		// never [[0...], [7...], [4...].
		IterateSets(devices[1:], size-1, func(set []T) {
			set = append([]T{devices[0]}, set...)

			p := SetCountPadding(set)
			if !(p == 0 || p == padding) {
				return
			}

			remaining := []T{}
			for _, d := range devices {
				if !SetContains(set, d) {
					remaining = append(remaining, d)
				}
			}

			iterate(remaining, size, append(accum, set))
		})
	}

	iterate(devices, size, [][]T{})
}
