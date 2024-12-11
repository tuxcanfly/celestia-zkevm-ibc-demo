# ZK-EVM X Celestia Token Transfer Demo

> ⚠️ **Warning**
> This repository is a work in progress and under active development.

This repo exists to showcase transferring tokens to and from a Cosmos SDK chain (representing Celestia) and a ZK proveable EVM using the [IBC-Eureka solidity contracts](https://github.com/cosmos/solidity-ibc-eureka/tree/main/src). The diagram below is meant to detail the components involved and, at a high level, how they interact with one another.

![mvp-zk-accounts](./mvp-zk-accounts.png)

For more information refer to the [architecture document](./ARCHITECTURE.md). Note that the design is subject to change.

## Contributing

1. Complete the solidity-ibc-eureka [requirements](https://github.com/cosmos/solidity-ibc-eureka?tab=readme-ov-file#requirements)
    1. After you `cp .env.example .env` you will need to set the `PRIVATE_KEY=0x82bfcfadbf1712f6550d8d2c00a39f05b33ec78939d0167be2a737d691f33a6a`
1. Install [Docker](https://docs.docker.com/get-docker/)
1. Fork this repo and clone it
1. Set up the git submodule for solidity-ibc-eureka via:

    ```shell
    git submodule init
    git submodule update
    ```

### Local development

```shell
# Start a local development environment by running
docker compose up --detach

# Deploy smart contracts on the EVM roll-up
make deploy-contracts

# Build the demo binary
make build-demo

# Run the demo binary
make run-demo # NOTE: this currently doesn't work
```

While deploying contracts, if you hit an error like: `[Revert] vm.envString: environment variable "E2E_FAUCET_ADDRESS" not found` then comment out the lines that use that environment variable from `./solidity-ibc-eureka/E2ETestDeploy.s.sol`.

### Helpful commands

```shell
# See the running contains
docker ps

# You can view the logs from a running container via Docker UI or:
docker logs beacond
docker logs celestia-network-bridge
docker logs celestia-network-validator
docker logs simapp-validator
docker logs reth
```
