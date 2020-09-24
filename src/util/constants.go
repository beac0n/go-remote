package util

import _ "golang.org/x/crypto/sha3"
import "crypto"

const AesKeySize = 32

const HashFunction = crypto.SHA3_512

const TimestampLen = 8

const EncryptedDataLen = 36

const FilePathTimestamp = "./.timestamp"
const KeySuffix = ".key"

const SecInNs = int64(1000000000)
const MaxFileSizeMb = float64(5)
