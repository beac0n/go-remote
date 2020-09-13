package util

import _ "golang.org/x/crypto/sha3"
import "crypto"

const AesKeySize = 32
const RsaKeySize = 4096

const HashFunction = crypto.SHA3_512

const TimestampLen = 8
const SaltLen = 2040
const TotalDataLen = TimestampLen + SaltLen

const EncryptedDataLen = 2644

const FilePathTimestamp = "./.timestamp"
const ServerSuffix = "server"
const ClientSuffix = "client"

const SecInNs = int64(1000000000)
const MaxFileSizeMb = float64(5)
