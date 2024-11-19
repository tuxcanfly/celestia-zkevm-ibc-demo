#!/bin/bash
simd init zk-ibc-demo --home testing/files/simapp-validator

simd keys add validator --keyring-backend test --recover 

exec simd start --home testing/files/simapp-validator