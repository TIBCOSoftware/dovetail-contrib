# dovetail-tools

This package provides script to build Dovetail flow models into chaincode CDS file or client app executables.  The script uses a pre-configured docker container to build Dovetail flow models, and thus developers do not have to configure a local dev environment to compile and build Dovetail artifacts.

## Build Hyperledger Fabric chaincode

For Hyperledger Fabric release 1.4, chaincode is packaged as a CDS file and then installed on peer nodes, although it will change for release 2.0.  After you complete a chaincode model in Flogo Enterprise or TIBCO Cloud, you can download the model as a `json` file, e.g., [iou.json](../samples/iou/iou.json), and then build the chaincode CDS package as follows:

```bash
cd /path/to/dovetail-contrib/hyperledger-fabric/dovetail-tools
./build.sh cds -f ../samples/iou/iou.json
```

This command will write the resulting CDS file as `/path/to/dovetail-contrib/hyperledger-fabric/dovetail-tools/work/iou_cc_1.0.cds`.  You can use other command options to override the default chaincode name and version number, which is described in the following section.

This command will start a `dovetail-tools` container to execute the build commands if the container is not running already.

## Build application executable

After you complete a client app model in Flogo Enterprise or TIBCO Cloud, you can download the model as a `json` file, e.g., [iou_client.json](../samples/iou/iou_client.json), and then build the executable for a specified operating system as follows:

```bash
cd /path/to/dovetail-contrib/hyperledger-fabric/dovetail-tools
./build.sh client -f ../samples/iou/iou_client.json -s darwin
```

This command will write the resulting executable for MacOS as `/path/to/dovetail-contrib/hyperledger-fabric/dovetail-tools/work/iou_client_darwin_amd64`.  You can use other command options to override the default executable name and hardware architecture, which is described in the following section.

This command will start a `dovetail-tools` container to execute the build commands if the container is not running already.

## Other build options

Following command prints out all build options:

```bash
cd /path/to/dovetail-contrib/hyperledger-fabric/dovetail-tools
./build.sh -h

Usage:
  build.sh <cmd> [ args ]
    <cmd> - one of the following
      - 'cds' - build chaincode cds with args: -f model-json [ -n cc-name -v cc-version ]
      - 'client' - build client executable with args: -f model-json [ -n exe-name -s GOOS -a GOARCH ]
      - 'start' - start docker-tools docker container
      - 'shutdown' - shutdown docker-tools docker container
    -f <model json> - path of the flogo model json file
    -n <name> - name of the chaincode or client exe file, default <model>_cc_<version>.cds or <model>_<goos>_<goarch>
    -v <cc-version> - version of the chaincode, e.g., '1.0' (default)
    -s <GOOS> - GOOS platform for the client exe, e.g., linux (default), darwin, or windows
    -a <GOARCH> - GOARCH for the the client exe, e.g., amd64 (default), or 386
  build.sh -h (print this message)
```

By using this script, you can also start the `dovetail-tools` container in advance, and use it to execute multiple build commands, i.e.,

```bash
cd /path/to/dovetail-contrib/hyperledger-fabric/dovetail-tools
./build.sh start
```

## Build docker image for dovetail-tools

By default, the build script uses a docker image of `dovetail-tools` in Docker Hub, as specified in [dovetail-tools.yaml](./dovetail-tools.yaml).  You can build your own docker image using the following script if you want to include more advanced features of the TIBCO Flogo Enterprise or TIBCO Cloud.

```bash
cd /path/to/dovetail-contrib/hyperledger-fabric/dovetail-tools
./docker-image.sh build -e ~/felib/flogo.zip
./docker-image.sh upload -u dhuser -d
```

This command builds a `dovetail-tools` docker image to support extension libs of Flogo Enterprise included in `flogo.zip`, and then uploads the docker image to Docker Hub.  Note that you will need TIBCO license to use the Flogo Enterprise extensions.  Without a license, you can use only Flogo open-source extensions and Dovetail extensions, which are also open-source.

Other options of the `docker-image` command are as follows:

```bash
cd /path/to/dovetail-contrib/hyperledger-fabric/dovetail-tools
./docker-image.sh -h

Usage:
  docker-image.sh <cmd> [ args ]
    <cmd> - one of the following
      - 'build' - build dovetail-tools image with optional args: [ -n name -v version -e flogo-zip ]
      - 'upload' - upload image to docker hub with optional args: -u user -p passwd [ -n name -v version -d ]
    -n <image name> - name of the docker image, e.g., dovetail-tools (default)
    -v <image version> - version of the docker image, e.g., 'v1.0.0' (default)
    -e <flogo lib> - path of the zip file for Flogo Enterprise library
    -u <user> - user name for a docker hub account
    -p <passwd> - password for a docker hub account
    -d - flag to cleanup local docker images
  docker-image.sh -h (print this message)
```