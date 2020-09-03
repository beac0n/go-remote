package util

const ServerKeyLen = 32
const ClientKeyLen = 64
const CryptoKeyLen = 32

const ServerKeyFileLen = ServerKeyLen + CryptoKeyLen
const ClientKeyFileLen = ClientKeyLen + CryptoKeyLen

const TimestampLen = 8
const SaltLen = 56
const MsgLen = TimestampLen + SaltLen

const SigLen = 64
const DataLen = SigLen + MsgLen

const EncryptedDataLen = 156

const FilePathTimestamp = "./.timestamp"
const ServerSuffix = "server"
const ClientSuffix = "client"

const SecInNs = int64(1000000000)