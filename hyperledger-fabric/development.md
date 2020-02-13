# Setup Development Environment

Dovetail fabric extensions can be used in one of the following 2 modeling environments:

- [TIBCO Flogo® Enterprise v2.8.0](https://docs.tibco.com/products/tibco-flogo-enterprise-2-8-0)
- [TIBCO Cloud Integration (TCI)](https://cloud.tibco.com/)

## Prerequisite for local development

Following are packages required for setting up development evironment locally on Mac or Linux.

- Download [TIBCO Flogo® Enterprise 2.8.0](https://edelivery.tibco.com/storefront/eval/tibco-flogo-enterprise/prod11810.html) only if you do not use TCI.
- Download Hyperledger Fabric samples and executables of latest production release as described [here](https://github.com/hyperledger/fabric-samples/tree/release-1.4). Current release works with Fabric release 1.4.4.
- Clone [dovetail-contrib](https://github.com/TIBCOSoftware/dovetail-contrib) with Flogo extension for Hyperledger Fabric

For example, you can install these packages to a demo directory by using the following commands:

```bash
cd $HOME/work/DovetailDemo
curl -sSL http://bit.ly/2ysbOFE | bash -s
export PATH=$HOME/work/DovetailDemo/fabric-samples/bin:$PATH
export FAB_PATH=$HOME/work/DovetailDemo/fabric-samples
git clone https://github.com/TIBCOSoftware/dovetail-contrib.git
```

This is all you need to start developing chaincode and client apps, and build artifacts using a [dovetail-tools](./dovetail-tools) docker container.  Refer the [IOU](samples/iou) sample for step-by-step instructions about how to implement, build, deploy and run a blockchain application.

If you want to build artifacts locally, instead of using the `dovetail-tools` docker container, you will need to install and configure the following tools:

- [Install Go](https://golang.org/doc/install).  The current release requires Go 1.12.x to build Hyperledger Fabric chaincode.
- Download and install [flogo-cli](https://github.com/project-flogo/cli)

This document assumes that you [install Go](https://golang.org/doc/install) first, and then configure the dev environment as follows, i.e.,

```bash
go get -u github.com/project-flogo/cli/...

# rebuild `fabric-tools`, assuming dovetail is installed at $HOME/work/DovetailDemo/dovetail-contrib
cd $HOME/work/DovetailDemo/dovetail-contrib/hyperledger-fabric/fabric-tools
go install

# assuming Flogo Enterprise is installed at $HOME/work/DovetailDemo/flogo/2.8
export FE_HOME=$HOME/work/DovetailDemo/flogo/2.8
```

This configuration should be good to build and test all [samples](samples) locally.  Mac users can also update the `docker-compose-cli.yaml` to speed up chaincode installation as described [here](https://docs.docker.com/compose/compose-file/#caching-options-for-volume-mounts-docker-for-mac), i.e.,

```bash
cd $HOME/work/DovetailDemo/fabric-samples/first-network
sed -i -e "s/github.com\/chaincode.*/github.com\/chaincode:cached/" ./docker-compose-cli.yaml
```

## Configure TIBCO Flogo® Enterprise

If you have the license for the `TIBCO Flogo® Enterprise`, you can use it to edit models of the Dovetail samples.  We use the [marble](samples/marble) sample to describe the initial setup of Flogo Enterprise UI when you start the first app model.

- Start TIBCO Flogo® Enterprise as described in [User's Guide](https://docs.tibco.com/pub/flogo/2.8.0/doc/pdf/TIB_flogo_2.8_users_guide.pdf?id=2)
- Open http://localhost:8090 in Chrome web browser.
- Open [Extensions](http://localhost:8090/wistudio/extensions) link, and upload [`fabricExtension.zip`](fabricExtension.zip).  Note that you can generate this `zip` by using the script [`zip-fabric.sh`](zip-fabric.sh).
- Upload [`fabclientExtension.zip`](fabclientExtension.zip).  Note that you can generate this `zip` by using the script [`zip-fabclient.sh`](zip-fabclient.sh).
- Create new Flogo App of name `marble_app` and choose `Import app` to import the model [`marble_app.json`](samples/marble/marble_app.json)
- Optionally, you can then add or update the flow models in the browser.
- After you are done editing, export the Flogo App, and copy the downloaded model file, i.e., [`marble_app.json`](marble_app.json) to the [marble](samples/marble) sample folder.

Note: if a client app uses `General` triggers/acrivities included by the `TIBCO Flogo® Enterprise`, you need to use the following script to enable go-module for these components:

```bash
cd $HOME/work/DovetailDemo/dovetail-contrib/fe-generator
./init-gomod.sh ${FE_HOME}
```

## Modeling with TIBCO Cloud Integration (TCI)

If you are already a subscriber of [TIBCO Cloud Integration (TCI)](https://cloud.tibco.com/), or you plan to sign-up for a TCI trial, you can use TCI to edit app models exported from `Dovetail` or `TIBCO Flogo Enterprise`.  Refer to [Modeling with TCI](tci) for more detailed instructions.

## Build and deploy chaincode to Hyperledger Fabric

It is simplest to build chaincode using the `dovetail-tools` docker container as described in [README](dovetail-tools/README.md).  Refer to the [IOU](samples/iou) sample on how the simple build process work.

The other sampples describes local build process.  We use the [marble](samples/marble) sample to describe the steps to deploy and invoke chaincode with the `byfn` network of the `fabric-samples`.

- In the [marble](samples/marble) sample folder, execute `make create` and `make build` to generate and build chaincode source code from the flogo model [`marble_app.json`](samples/marble/marble_app.json).
- Execute `make deploy` to build and deploy the chaincode to the `fabric-samples` chaincode folder.  Note that you may need to set `FAB_PATH` env or edit the [`Makefile`](samples/marble/Makefile) to match the installation folder of `fabric-samples` if it is not downloaded to the default location under `$GOPATH`.
- Execute `make metadata` to generate `metadata` for client apps.

The detailed commands of the above steps are as follows:

```bash
cd $HOME/work/DovetailDemo/dovetail-contrib/hyperledger-fabric/samples/marble
make create
make build
make deploy
make metadata
```

## Install and test chaincode using fabric sample byfn network

Start Hyperledger Fabric sample `byfn` network with CouchDB:

```bash
cd $HOME/work/DovetailDemo/fabric-samples/first-network
./byfn.sh up -n -s couchdb
```

Use `cli` docker container to install and instantiate the `marble_cc` chaincode.

```bash
cd $HOME/work/DovetailDemo/dovetail-contrib/hyperledger-fabric/samples/marble
make cli-init
```

Test the instantiated chaincode from `cli` docker container, i.e.,

```bash
cd $HOME/work/DovetailDemo/dovetail-contrib/hyperledger-fabric/samples/marble
make cli-test
```

Note that developers can also use Fabric dev-mode to test chaincode (refer [dev](samples/marble/dev.md) for more details).  For issues regarding how to work with the Fabric network, please refer the [Hyperledger Fabric docs](https://hyperledger-fabric.readthedocs.io/en/latest/build_network.html).
