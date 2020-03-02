# MSP utility

This utility uses `fabric-tools` container to generate artifacts for a Hyperledger Fabric network, including genesis block, transactions for channel creation and updates, etc.

## Start Fabric-tools

Example:

```bash
cd ./msp
./msp-util.sh start -p netop1 -t k8s
```

This starts a fabric-tools container using the config file [netop1.env](../config/netop1.env), which must be configured and put in the [config](../config) folder.  The container will be running in `docker-desktop` Kubernetes on Mac.  Non-Mac users should specify a different `-t` value to run it in another environment supported by your platform. The following command prints out the supported options:

```bash
./msp-util.sh -h

Usage:
  msp-util.sh <cmd> [-p <property file>] [-t <env type>] [-o <consensus type>] [-c <channel name>]
    <cmd> - one of the following commands
      - 'start' - start tools container to run msp-util
      - 'shutdown' - shutdown tools container for the msp-util
      - 'bootstrap' - generate bootstrap genesis block and test channel tx defined in network spec
      - 'genesis' - generate genesis block of specified consensus type, with argument '-o <consensus type>'
      - 'channel' - generate channel creation tx for specified channel name, with argument '-c <channel name>'
      - 'mspconfig' - print MSP config json for adding to a network, output in '/Users/yxu/work/DovetailDemo/dovetail-contrib/hyperledger-fabric/operation/netop1.com/tool'
      - 'orderer-config' - print orderer RAFT consenter config for adding to a network, with arguments -s <start-seq> [-e <end-seq>]
      - 'build-cds' - build chaincode cds package from flogo model, with arguments -m <model-json> [-v <version>]
      - 'build-app' - build linux executable from flogo model, with arguments -m <model-json>
    -p <property file> - the .env file in config folder that defines network properties, e.g., netop1 (default)
    -t <env type> - deployment environment type: one of 'docker', 'k8s' (default), 'aws', 'az', or 'gcp'
    -o <consensus type> - 'solo' or 'etcdraft' used with the 'genesis' command
    -c <channel name> - name of a channel, used with the 'channel' command
    -s <start seq> - start sequence number (inclusive) for orderer config
    -e <end seq> - end sequence number (exclusive) for orderer config
    -m <model json> - Flogo model json file
    -v <cc version> - version of chaincode
  msp-util.sh -h (print this message)
```

* Docker-compose users can use option `-t docker`
* Azure users can refer instructions in the folder [az](../az) to run it from an Azure `bastion` VM instance, which uses default option `-t az`.
* AWS users can refer instructions in the folder [aws](../aws) to run it from an Amazon `bastion` EC2 instance, which uses default option `-t aws`.
* Google cloud users can refer instructions in the folder [gcp](../gcp) to run it from a GCP `bastion` VM instance, which uses default option `-t gcp`.

## Bootstrap artifacts of Fabric network

Example:

```bash
cd ./msp
./msp-util.sh bootstrap -p netop1 -t k8s
```

This uses the `tool` container to generate the genesis block and channel transactions for a test channel `mychannel` as specified in the network specification file [netop1.env](../config/netop1.env). The result is stored in the folder [netop1.com/tool](../netop1.com/tool) as `genesis.block`, `channel.tx`, and `anchors.tx`.

When this command is used on `AWS`, `Azure` or `GCP`, the generated files will be stored in a cloud file system mounted on the `bastion` host, e.g., a mounted folder `/mnt/share/netop1.com/tool` in an `EFS` file system on `AWS` or an `Azure Files` storage on `Azure` or a `Filestore` volume on `GCP`.

## Generate genesis block of a specified consensus type

Example:

```bash
cd ./msp
./msp-util.sh genesis -p netop1 -t k8s -o solo
```

This will create a genesis block `solo-genesis.block` used for `solo` consensus type.

## Generate channel transactions for a specified channel name

Example:

```bash
cd ./msp
./msp-util.sh channel -p netop1 -t k8s -c testchan
```

This will create the channel creation and anchor transactions for a channel named `testchan`.  The created files are named `testchan.tx` and `testchan-anchors.tx`.

## Shutdown and cleanup

Example:

```bash
cd ./msp
./msp-util.sh shutdown -p netop1 -t k8s
```

This shuts down the `tool` container.

## TODO

More operations will be supported by this utility, including

* Update transaction for adding new orderers;
* Update transaction for adding new peer organizations;
* Update transaction for adding new orderer organizations
