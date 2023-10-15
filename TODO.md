# TODO

## Client

- [ ] Go client
- [ ] Javascript client

## Server

- [x] Delete API
- [ ] Scan / Keys API
- [ ] Graceful shutdown
- [ ] Support redis tcp protocol
- [ ] Authentication
- [ ] Authorization
- [ ] Telemetry

## Storage engine

- [x] Delete tombstone
- [x] Backend interface
- [ ] Global index instead of 1 index per segment?
- [ ] Merging: delete tombstones and write hint file
- [ ] Snapshot isolation: MVCC

## Distributed

- [ ] Service discovery with serf
- [ ] Single leader replication with raft
- [ ] Configurable consistency modes

## Other

- [ ] Open config file if exists config.yaml
- [ ] Write README.md
- [ ] Create architecture diagram
- [ ] CI/CD
- [ ] Create helm package
- [ ] Homebrew
