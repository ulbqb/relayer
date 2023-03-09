# LBM Sims Demo

![demo](./images/demo.gif)

>logging in Relayer has changed since this gifs creation

While the relayer is under active development, it is meant primarily as a learning
tool to better understand the Inter-Blockchain Communication (IBC) protocol. In
that vein, the following lbmsim-lbmsim demonstrates the core functionality which will
remain even after the changes:

```bash
# ensure go and jq are installed 
# Go Documentation: https://golang.org/doc/install
# jq Documentation: https://stedolan.github.io/jq/download

# First, download and build the lbm sim source code so we have a working blockchain to test against
$ make get-lbmsim build-lbmsim

# lbmsim-lbmsim creates two lbm sim chains with data directories in this repo
# it also builds and configures the relayer for operations with those chains
$ ./examples/demo/scripts/lbmsim-lbmsim
# NOTE: If you want to stop the two lbm sim chains running in the background use `killall simd`

# At this point the relayer --home directory is ready for normal operations between
# lbmsim-0 and lbmsim-1. Looking at the folder structure of the relayer at this point is helpful
# NOTE: to install tree try `brew install tree` on mac or `apt install tree` on linux
$ tree ~/.relayer

# See if the chains are ready to relay over
$ rly chains list

# See the current status of the path you will relay over
$ rly paths list

# Now you can connect the two chains with one command:
$ rly tx link lbmsim-lbmsim -d -t 3s

# Check the token balances on both chains
$ rly q balance lbmsim-0
$ rly q bal lbmsim-1

# Then send some tokens between the chains
$ rly tx transfer lbmsim-0 lbmsim-1 1000000samoleans $(rly chains address lbmsim-1) channel-0

# Relay packets/acknowledgments. 
# Running `rly start lbmsim-lbmsim` essentially loops these two commands
$ rly tx relay-pkts lbmsim-lbmsim channel-0 -d
$ rly tx relay-acks lbmsim-lbmsim channel-0 -d

# See that the transfer has completed
$ rly q bal lbmsim-0
$ rly q bal lbmsim-1

# Send the tokens back to the account on lbmsim-0
$ rly tx transfer lbmsim-1 lbmsim-0 1000000ibc/27A6394C3F9FF9C9DCF5DFFADF9BB5FE9A37C7E92B006199894CF1824DF9AC7C $(rly chains addr lbmsim-0) channel-0
$ rly tx relay-pkts lbmsim-lbmsim channel-0 -d
$ rly tx relay-acks lbmsim-lbmsim channel-0 -d

# See that the return trip has completed
$ rly q bal lbmsim-0
$ rly q bal lbmsim-1

# NOTE: you will see the stake balances decreasing on each chain. This is to pay for fees
# You can change the amount of fees you are paying on each chain in the configuration.
```

---

[<-- Pruning Settings](./node_pruning.md) - [LBM Sim and Gaia Demo -->](./lbmsim-gaia-demo.md)