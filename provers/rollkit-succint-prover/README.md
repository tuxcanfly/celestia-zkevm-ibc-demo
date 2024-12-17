# Rollkit Succint Prover

The Rollkit Succint Prover is a gRPC service that generates zero-knowledge proofs for Rollkit EVM state transitions and data membership. It is designed to work with IBC (Inter-Blockchain Communication) and specifically implements proofs compatible with the ICS-07 Tendermint client specification.

## Usage

> ⚠️ **Warning**
> This gRPC service is still under development and may not work as described

To run the server you will need to clone the repo and install rust and cargo. To run the node you also need to set the following environment variables:

- `TENDERMINT_RPC_URL` - the url of the tendermint chain you are proving.
- `EVM_RPC_URL` the json rpc url of the evm chain you are generating the proofs for.
- `EVM_CONTRACT_ADDRESS` - the evm address of the tendermint sp1 ics07 contract.

To then run the server (on port `:50001`):

```
cargo run
```


## Protobuf

gRPC depends on proto defined types. These are stored in `proto/prover/v1` from the root directory.
