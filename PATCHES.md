# Patches

- `storage/sys_unix.go`: treat files as allocated when the filesystem lacks `SEEK_DATA` support.
- `torrent-piece-request-order.go`: disable the pending-piece debug assertion.
- `tracker/udp/client.go`: serialize connects, protect connection state and retry an expired connection ID without invalidating a newer one.
- `config.go`, `torrent.go`, `peer.go`, `peerconn.go`: optionally cancel requests for unwanted pieces and distinguish cancellation acknowledgements from rejected requests.
- `client-tracker-announcer.go`, `torrent.go`: a re-added torrent starts a new announce session, so it announces Started instead of inheriting the dropped torrent's interval.
