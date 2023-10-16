# TODO

## Client

- [ ] Go client
- [ ] CLI client
- [ ] Javascript / Typescript client

## Server

- [x] Delete API
- [ ] Scan / Keys API
- [ ] Graceful shutdown
- [ ] Authentication
- [ ] Authorization
- [ ] Telemetry
- [ ] Support redis tcp protocol

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
- [ ] Sharding

## Other

- [ ] Fix precommit golangci-lint errors due to checking staging area only
- [ ] Open config file if exists config.yaml
- [ ] Write README.md
- [ ] Create architecture diagram
- [ ] CI/CD
- [ ] Create helm package
- [ ] Homebrew
