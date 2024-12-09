# ZK-EVM X Celestia Token Transfer Demo

> ⚠️ **Warning**
> This repository is a work in progress and under active development.

This repo exists to showcase transferring tokens to and from a Cosmos SDK chain (representing Celestia) and a ZK proveable EVM using the [IBC-Eureka solidity contracts](https://github.com/cosmos/solidity-ibc-eureka/tree/main/src). The diagram below is meant to detail the components involved and, at a high level, how they interact with one another.

![mvp-zk-accounts](./mvp-zk-accounts.png)

For more information refer to the [architecture document](./ARCHITECTURE.md). Note that the design is subject to change.

## Contributing

1. Install the [requirements](https://github.com/cosmos/solidity-ibc-eureka?tab=readme-ov-file#requirements) listed in the solidity-ibc-eureka README
1. Install [Docker](https://docs.docker.com/get-docker/)
1. Fork this repo and clone it
1. Set up the git submodule for solidity-ibc-eureka via:

    ```shell
    git submodule init
    git submodule update

    # TODO: remove this fetch and checkout after https://github.com/celestiaorg/celestia-zkevm-ibc-demo/issues/17 is resolved
    git fetch --all
    git checkout d952d312065f9c3689a9f3100dc98cc0c86c2f77
    ```

## Local Development

```shell
# Start a local development environment by running
docker compose up --detach

# Deploy smart contracts on the EVM roll-up
make deploy-contracts
```
