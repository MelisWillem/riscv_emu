package riscv

type Fetcher interface {
	Fetch(pc int32) (int32, bool)
}

type ArrayFetcher struct {
	data []int32
}

func (f ArrayFetcher) Fetch(pc int32) (int32, bool) {
	if pc < int32(len(f.data)) {
		return f.data[pc], true
	}
	return 0, false
}

func NewArrayFetch(data []int32) Fetcher {
	return ArrayFetcher{data}
}
