# Ata Net Room

Websocket relay server with separated rooms.

## Setup

Generate SSL certs for dev:
`openssl genrsa -out server.key 2048`
`openssl ecparam -genkey -name secp384r1 -out server.key`