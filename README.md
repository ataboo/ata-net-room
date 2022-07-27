# Ata Net Room

Websocket relay server with separated rooms.

## Setup

Generate SSL certs for dev:
`openssl ecparam -genkey -name secp384r1 -out server.key`
`openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650`