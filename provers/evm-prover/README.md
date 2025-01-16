# EVM Prover

The EVM Prover is a gRPC service that generates zero-knowledge proofs for EVM state transitions. It is designed to work with IBC (Inter-Blockchain Communication) and specifically implements proofs compatible with the ICS-07 Tendermint client specification.

## Usage

> [!WARNING] 
> This gRPC service is still under development and may not work as described

To run the server you will need to clone the repo and install rust and cargo. To run the node you also need to set the following environment variables:

- `TENDERMINT_RPC_URL` - the url of the tendermint chain you are proving.
- `RPC_URL` the json rpc url of the evm chain you are generating the proofs for.
- `CONTRACT_ADDRESS` - the evm address of the tendermint sp1 ics07 contract.

To then run the server (on port `:50051`):

```
cargo run
```

To use the SP1 Prover Network you should also populate the `SP1_PROVER` and `SP1_PRIVATE_KEY` environment variables. You can also use a `.env` file for all environment variables


## Protobuf

gRPC depends on proto defined types. These are stored in `proto/prover/v1` from the root directory.
