package misc

func ipSlice(ip [4]byte) []byte {
	return []byte{ip[0], ip[1], ip[2], ip[3]}
}

func ipArrary(ip []byte) [4]byte {
	return [4]byte{ip[0], ip[1], ip[2], ip[3]}
}
