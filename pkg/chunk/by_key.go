package chunk

// ByKey allow you to sort chunks by ID
type ByKey []Chunk

func (cs ByKey) Len() int           { return len(cs) }
func (cs ByKey) Swap(i, j int)      { cs[i], cs[j] = cs[j], cs[i] }
func (cs ByKey) Less(i, j int) bool { return cs[i].externalKey() < cs[j].externalKey() }

// unique will remove duplicates from the input.
// list must be sorted.
func unique(cs ByKey) ByKey {
	if len(cs) == 0 {
		return ByKey{}
	}

	result := make(ByKey, 1, len(cs))
	result[0] = cs[0]
	i, j := 0, 1
	for j < len(cs) {
		if result[i].externalKey() == cs[j].externalKey() {
			j++
			continue
		}
		result = append(result, cs[j])
		i++
		j++
	}
	return result
}

// merge will merge & dedupe two lists of chunks.
// list musts be sorted and not contain dupes.
func merge(a, b ByKey) ByKey {
	result := make(ByKey, 0, len(a)+len(b))
	i, j := 0, 0
	for i < len(a) && j < len(b) {
		if a[i].externalKey() < b[j].externalKey() {
			result = append(result, a[i])
			i++
		} else if a[i].externalKey() > b[j].externalKey() {
			result = append(result, b[j])
			j++
		} else {
			result = append(result, a[i])
			i++
			j++
		}
	}
	for ; i < len(a); i++ {
		result = append(result, a[i])
	}
	for ; j < len(b); j++ {
		result = append(result, b[j])
	}
	return result
}

// nWayUnion will merge and dedupe n lists of chunks.
// lists must be sorted and not contain dupes.
func nWayUnion(sets []ByKey) ByKey {
	l := len(sets)
	switch l {
	case 0:
		return ByKey{}
	case 1:
		return sets[0]
	case 2:
		return merge(sets[0], sets[1])
	default:
		var (
			split = l / 2
			left  = nWayUnion(sets[:split])
			right = nWayUnion(sets[split:])
		)
		return nWayUnion([]ByKey{left, right})
	}
}

// nWayIntersect will interesct n sorted lists of chunks.
func nWayIntersect(sets []ByKey) ByKey {
	l := len(sets)
	switch l {
	case 0:
		return ByKey{}
	case 1:
		return sets[0]
	case 2:
		var (
			left, right = sets[0], sets[1]
			i, j        = 0, 0
			result      = []Chunk{}
		)
		for i < len(left) && j < len(right) {
			if left[i].externalKey() == right[j].externalKey() {
				result = append(result, left[i])
			}

			if left[i].externalKey() < right[j].externalKey() {
				i++
			} else {
				j++
			}
		}
		return result
	default:
		var (
			split = l / 2
			left  = nWayIntersect(sets[:split])
			right = nWayIntersect(sets[split:])
		)
		return nWayIntersect([]ByKey{left, right})
	}
}
