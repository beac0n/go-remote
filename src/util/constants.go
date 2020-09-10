package util

import _ "golang.org/x/crypto/sha3"
import "crypto"

const KeySize = 4096
const HashFunction = crypto.SHA3_512

const TimestampLen = 8
const SaltLen = 374

const EncryptedDataLen = 512

const FilePathTimestamp = "./.timestamp"
const ServerSuffix = "server"
const ClientSuffix = "client"

const SecInNs = int64(1000000000)
