package util

func XorEncrypt(data []byte, key byte) []byte {
	encryptedData := make([]byte, len(data))
	copy(encryptedData, data)
	for i := 0; i < len(encryptedData); i++ {
		encryptedData[i] = encryptedData[i] ^ key
	}
	return encryptedData
}

func XorDecrypt(data []byte, key byte) []byte {
	return XorEncrypt(data, key) // XOR加密和解密操作是相同的
}
