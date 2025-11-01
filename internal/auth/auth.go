package auth

type argonParams struct {
	threads   uint8
	saltLen   uint8
	time      uint32
	memory    uint32
	keyLength uint32
}

var authParams = argonParams{
	threads:   1,
	saltLen:   16,
	time:      2,
	memory:    16 * 1024,
	keyLength: 32,
}
