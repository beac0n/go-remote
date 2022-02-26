package util

import _ "golang.org/x/crypto/sha3"
import "crypto"

const ConfigDir = "/etc/go-remote"
const FilePathTimestamp = "/etc/go-remote/go-remote-timestamp"
const SocketPath = "/etc/go-remote/go-remote.sock"

const HashFunction = crypto.SHA3_512

const TimestampLen = 8

const AesKeySize = 32
const EncryptedDataLen = 36

const SecInNs = int64(1000000000)
