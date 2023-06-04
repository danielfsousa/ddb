# TODO

## Client

- [ ] Go client
- [ ] Javascript client

## Server

- [x] Delete API
- [ ] Graceful shutdown
- [ ] Support redis tcp protocol

## Storage engine

- [x] Delete tombstone
- [x] Backend interface
- [ ] Global index instead of 1 index per segment?
- [ ] Merging: delete tombstones and write hint file
- [ ] Snapshot isolation: MVCC

## Distributed

- [ ] Single leader replication with raft
- [ ] Service discovery with serf

## Other

- [ ] Open config file if exists config.yaml
- [ ] Write README.md
- [ ] Create architecture diagram
- [ ] CI/CD
- [ ] Create helm package
