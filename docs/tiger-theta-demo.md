# Mychain and Gaia Testnet Demo

![demo](./images/demo.gif)

>logging in Relayer has changed since this gifs creation

While the relayer is under active development, it is meant primarily as a learning
tool to better understand the Inter-Blockchain Communication (IBC) protocol. In
that vein, the following tiger-theta demonstrates the core functionality which will
remain even after the changes:

```bash
# ensure go and jq are installed 
# Go Documentation: https://golang.org/doc/install
# jq Documentation: https://stedolan.github.io/jq/download

# First, download and build the lbm sim and gaia source code so we have a working blockchain to test against
$ make get-lbmsim build-lbmsim

# tiger-theta creates the lbm sim and gaia-based chains with data directories in this repo
# it also builds and configures the relayer for operations with those chains
$ ./examples/demo/scripts/tiger-theta
# NOTE: If you want to stop the lbm sim and gaia-based chains running in the background use `killall simd && killall gaiad`

# At this point the relayer --home directory is ready for normal operations between
# tiger-0 and theta-testnet-001. Looking at the folder structure of the relayer at this point is helpful
# NOTE: to install tree try `brew install tree` on mac or `apt install tree` on linux
$ tree ~/.relayer

# See if the chains are ready to relay over
$ rly chains list

# See the current status of the path you will relay over
$ rly paths list

# Now you can connect the two chains with one command:
# NOTE: Fetching client states is heavy, so add --overide to avoid it.
$ rly tx link tiger-theta -d -t 20s --override

# Check the token balances on both chains
$ rly q balance tiger-0
$ rly q bal theta-testnet-001

# Then send some tokens between the chains
# NOTE: Change correct channel id.
$ rly tx transfer tiger-0 theta-testnet-001 1000000tiger $(rly chains address theta-testnet-001) channel-0

# Relay packets/acknowledgments. 
# Running `rly start tiger-theta` essentially loops these two commands
$ rly tx relay-pkts tiger-theta channel-0 -d
$ rly tx relay-acks tiger-theta channel-0 -d

# See that the transfer has completed
$ rly q bal tiger-0
$ rly q bal theta-testnet-001

# Send the tokens back to the account on tiger-0
# NOTE: Change correct token id.
$ rly tx transfer theta-testnet-001 tiger-0 1000000ibc/27A6394C3F9FF9C9DCF5DFFADF9BB5FE9A37C7E92B006199894CF1824DF9AC7C $(rly chains addr tiger-0) channel-0
$ rly tx relay-pkts tiger-theta channel-0 -d
$ rly tx relay-acks tiger-theta channel-0 -d

# See that the return trip has completed
$ rly q bal tiger-0
$ rly q bal theta-testnet-001

# NOTE: you will see the stake balances decreasing on each chain. This is to pay for fees
# You can change the amount of fees you are paying on each chain in the configuration.
```

---

[<-- Pruning Settings](./node_pruning.md)