# Setup Development Environment
Dovetail fabric extensions can be used in one of the following 3 modeling environments:
- [Dovetail release v0.2.0](https://github.com/TIBCOSoftware/dovetail/releases)
- [TIBCO Flogo® Enterprise v2.6.1](https://docs.tibco.com/products/tibco-flogo-enterprise-2-6-1)
- [TIBCO Cloud Integration (TCI)](https://cloud.tibco.com/)

## Prerequisite for local development
Following are packages required for setting up development evironment locally on Mac or Linux.
- Download [TIBCO Flogo® Enterprise 2.6.1](https://edelivery.tibco.com/storefront/eval/tibco-flogo-enterprise/prod11810.html), or
- Download [Dovetail v0.2.0](https://github.com/TIBCOSoftware/dovetail/releases)
- [Install Go](https://golang.org/doc/install)
- Clone [Hyperledger Fabric](https://github.com/hyperledger/fabric)
- Install [Fabric CA binaries](https://hyperledger-fabric-ca.readthedocs.io/en/release-1.4/users-guide.html)
- Download Hyperledger Fabric samples and executables of latest production release as described [here](https://github.com/hyperledger/fabric-samples/tree/release-1.4)
- Download and install [flogo-cli](https://github.com/project-flogo/cli)
- Clone [dovetail-contrib](https://github.com/TIBCOSoftware/dovetail-contrib) with Flogo extension for Hyperledger Fabric

There are different ways to clone these packages.  This document assumes that you [install Go](https://golang.org/doc/install) first, and then install other packages under $GOPATH, i.e.,
```
go get -u github.com/hyperledger/fabric
go get -u github.com/hyperledger/fabric-ca/cmd/...
cd $GOPATH/src/github.com/hyperledger
curl -sSL http://bit.ly/2ysbOFE | bash -s
export PATH=$GOPATH/src/github.com/hyperledger/fabric-samples/bin:$PATH
go get -u github.com/project-flogo/cli/...
go get -u github.com/TIBCOSoftware/dovetail-contrib
```
Download the latest version of the Dovetail extension for Hyperledger Fabric from the [`develop` branch of the `dovetail-contrib`](https://github.com/TIBCOSoftware/dovetail-contrib/tree/develop), and rebuild `fabric-tools`, i.e.,
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib
git checkout develop
cd hyperledger-fabric/fabric-tools
go install
```
For Mac users, update the `docker-compose-cli.yaml` to speed up chaincode installation as described [here](https://docs.docker.com/compose/compose-file/#caching-options-for-volume-mounts-docker-for-mac), i.e.,
```
cd $GOPATH/src/github.com/hyperledger/fabric-samples/first-network
sed -i -e "s/github.com\/chaincode.*/github.com\/chaincode:cached/" ./docker-compose-cli.yaml
```

## Configure Dovetail
We use the [marble](samples/marble) sample to describe the initial setup of the Dovetail UI when you start the first app model.

- Download `TIB_dovetail-fabric_0.2.0_macosx_x86_64.zip` or the `Linux` version of the latest release from https://github.com/TIBCOSoftware/dovetail/releases. Unzip it to an installation folder.
- Start the `Dovetail` container using the installed script `dovetail/0.2/bin/run-studio.sh`.  Note that you'll need to install docker if it is not already done. Mac users can download `Docker Desktop for Mac` [here](https://hub.docker.com/editions/community/docker-ce-desktop-mac).
- Open http://localhost:8090 in Chrome web browser.
- Open [Extensions](http://localhost:8090/wistudio/extensions) link, and upload [`fabricExtension.zip`](fabricExtension.zip).  Note that you can generate this `zip` by using the script [`zip-fabric.sh`](zip-fabric.sh).
- Upload [`fabclientExtension.zip`](fabclientExtension.zip).  Note that you can generate this `zip` by using the script [`zip-fabclient.sh`](zip-fabclient.sh).
- Upload other Flogo open-source extensions from `github`.  The following 4 extensions are used by Dovetail samples, and are packaged in folder [tci](tci):
  - [GraphQL Trigger](tci/trigger-graphql.zip), which is an open-source project, [project-flogo/graphql](https://github.com/project-flogo/graphql)
  - [REST Trigger](tci/trigger-rest.zip), which is from the open-source project, [project-flogo/contrib](https://github.com/project-flogo/contrib)
  - [Log Activity](tci/activity-log.zip), which is from the open-source project, [project-flogo/contrib](https://github.com/project-flogo/contrib)
  - [Error Activity](tci/activity-error.zip), which is from the open-source project, [project-flogo/contrib](https://github.com/project-flogo/contrib)
- Create new Flogo App of name `marble_app` and choose `Import app` to import the model [`marble_app.json`](samples/marble/marble_app.json)
- Optionally, you can then add or update the flow models in the browser.
- After you are done editing, export the Flogo App, and copy the downloaded model file, i.e., [`marble_app.json`](samples/marble/marble_app.json) to the [marble](samples/marble) sample folder.

Note that when a Flogo model is imported to `Dovetail v0.2.0`, a `return` activity is automatically added to the end of all branches, which could be an issue if the `return` activity is not at the end of a flow.  Thus, you need to carefully remove the mistakenly added `return` activities after a model is imported.  Besides, the trigger config properties in HTTP client app, such as [audit_client](samples/audit/audit_client.json), are not imported correctly.  These issues will be fixed in a later release of `Dovetail`.

## Configure TIBCO Flogo® Enterprise
If you have the license for the `TIBCO Flogo® Enterprise`, you can use it to edit models exported from `Dovetail` or vice versa.  We use the [marble](samples/marble) sample to describe the initial setup of Flogo Enterprise UI when you start the first app model.

- Start TIBCO Flogo® Enterprise as described in [User's Guide](https://docs.tibco.com/pub/flogo/2.6.1/doc/pdf/TIB_flogo_2.6_users_guide.pdf?id=2)
- Open http://localhost:8090 in Chrome web browser.
- Open [Extensions](http://localhost:8090/wistudio/extensions) link, and upload [`fabricExtension.zip`](fabricExtension.zip).  Note that you can generate this `zip` by using the script [`zip-fabric.sh`](zip-fabric.sh).
- Upload [`fabclientExtension.zip`](fabclientExtension.zip).  Note that you can generate this `zip` by using the script [`zip-fabclient.sh`](zip-fabclient.sh).
- Create new Flogo App of name `marble_app` and choose `Import app` to import the model [`marble_app.json`](samples/marble/marble_app.json)
- Optionally, you can then add or update the flow models in the browser.
- After you are done editing, export the Flogo App, and copy the downloaded model file, i.e., [`marble_app.json`](marble_app.json) to the [marble](samples/marble) sample folder.

Note that when a Flogo model is imported to `Flogo® Enterprise v2.6.1`, a `return` activity is automatically added to the end of all branches, which could be an issue if the `return` activity is not at the end of a flow.  Thus, you need to carefully remove the mistakenly added `return` activities after a model is imported.  This issue will be fixed in a later release of the `Flogo® Enterprise`.

## Modeling with TIBCO Cloud Integration (TCI)
If you are already a subscriber of [TIBCO Cloud Integration (TCI)](https://cloud.tibco.com/), or you plan to sign-up for a TCI trial, you can use TCI to edit app models exported from `Dovetail` or `TIBCO Flogo Enterprise`.  Refer to [Modeling with TCI](tci) for more detailed instructions.

## Build and deploy chaincode to Hyperledger Fabric
We use the [marble](samples/marble) sample to describe the steps to deploy and invoke chaincode with the `byfn` network of the `fabric-samples`.

- In the [marble](samples/marble) sample folder, execute `make create` to generate chaincode source code from the flogo model [`marble_app.json`](samples/marble/marble_app.json).
- Execute `make deploy` to build and deploy the chaincode to the `fabric-samples` chaincode folder.  Note that you may need to edit the [`Makefile`](samples/marble/Makefile) and set `CC_DEPLOY` to match the installation folder of `fabric-samples` if it is not downloaded to the default location under `$GOPATH`.
- Execute `make package` to generate `cds` package for cloud deployment, and `metadata` for client apps.

The detailed commands of the above steps are as follows:
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/marble
make create
make deploy
make package
```
Note that if `make package` failed due to missing the `fabric-tools` executable, you can rebuild the tool as follows:
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric-tools
go install
```
If the command `fabric-tools package` failed with the following error:
```
panic: /debug/requests is already registered. You may have two independent copies of golang.org/x/net/trace in your binary, trying to maintain separate state. This may involve a vendored copy of golang.org/x/net/trace.

goroutine 1 [running]:
github.com/hyperledger/fabric/vendor/golang.org/x/net/trace.init.0()
	$GOPATH/src/github.com/hyperledger/fabric/vendor/golang.org/x/net/trace/trace.go:116 +0x1a4
```
you can delete the `trace` folder under `fabric/vendor` and rebuild the `fabric-tools`, and then retry `make package`, i.e.,
```
rm -R $GOPATH/go/src/github.com/hyperledger/fabric/vendor/golang.org/x/net/trace
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/fabric-tools
go install
```
## Install and test chaincode using fabric sample byfn network
Start Hyperledger Fabric sample `byfn` network with CouchDB:
```
cd $GOPATH/src/github.com/hyperledger/fabric-samples/first-network
./byfn.sh up -n -s couchdb
```
Use `cli` docker container to install and instantiate the `marble_cc` chaincode.
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/marble
make cli-init
```
Note that this script installs chaincode on 4 peer nodes using the `cli` container.  The default configuration is very slow on Mac due to slow volume mounts in the docker desktop for Mac.  The following [solution](https://docs.docker.com/compose/compose-file/#caching-options-for-volume-mounts-docker-for-mac) will speed up the chaincode installation by more than 4 times.  Thus, make the edit as follows if you have not done it already.
```
cd $GOPATH/src/github.com/hyperledger/fabric-samples/first-network
sed -i -e "s/github.com\/chaincode.*/github.com\/chaincode:cached/" ./docker-compose-cli.yaml
```

Test the instantiated chaincode from `cli` docker container, i.e.,
```
cd $GOPATH/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/samples/marble
make cli-test
```

Note that developers can also use Fabric dev-mode to test chaincode (refer [dev](samples/marble/dev.md) for more details).  For issues regarding how to work with the Fabric network, please refer the [Hyperledger Fabric docs](https://hyperledger-fabric.readthedocs.io/en/latest/build_network.html).
