# illustration
This example illustrates the extensions of the [project Dovetail](https://tibcosoftware.github.io/dovetail/) for developing chaincode and client apps of the [Hyperledger Fabric](https://www.hyperledger.org/projects/fabric).  The Dovetail extension contains triggers and activities used to construct [Flogo®](https://www.flogo.io/) models by visual programming with zero-code.  Flogo® models can be created, imported, edited, and/or exported by using [TIBCO Flogo® Enterprise](https://docs.tibco.com/products/tibco-flogo-enterprise-2-8-0) or [Dovetail](https://github.com/TIBCOSoftware/dovetail).

## Prerequisite
Follow the instructions [here](../../development.md) to setup the Dovetail development environment on Mac or Linux.

## Upload and view the illustration
- Start TIBCO Flogo® Enterprise or Dovetail.
- Open http://localhost:8090 in Chrome web browser.
- Create new Flogo App of name `all_fabric` and choose `Import app` to import the model [`all_fabric.json`](all_fabric.json)
- You can then view or edit the trigger and activities in graphical modeler of the TIBCO Flogo® Enterprise.

## Build and deploy chaincode to Hyperledger Fabric

This model is for illustration only.  For building and deploying real chaincode, refer other samples, e.g., [marble](../marble).
