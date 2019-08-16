# illustration
This example illustrates the extensions of the [project Dovetail](https://tibcosoftware.github.io/dovetail/) for developing chaincode and client apps of the [Hyperledger Fabric](https://www.hyperledger.org/projects/fabric).  The Dovetail extension contains triggers and activities used to construct [Flogo®](https://www.flogo.io/) models by visual programming with zero-code.  Flogo® models can be created, imported, edited, and/or exported by using [TIBCO Flogo® Enterprise](https://docs.tibco.com/products/tibco-flogo-enterprise-2-6-1) or [Dovetail](https://github.com/TIBCOSoftware/dovetail).

## Prerequisite
- Download [TIBCO Flogo® Enterprise 2.6](https://edelivery.tibco.com/storefront/eval/tibco-flogo-enterprise/prod11810.html). If you do not have access to `Flogo Enterprise`, you may sign up a trial on [TIBCO CLOUD Integration (TCI)](https://cloud.tibco.com/), or download Dovetail v0.2.0.  This sample uses `TIBCO Flogo® Enterprise`, but the models can be imported and edited by using Dovetail v0.2.0 and above.
- Clone dovetail-contrib with this Flogo extension

## Upload and view the illustration
- Start TIBCO Flogo® Enterprise as described in [User's Guide](https://docs.tibco.com/pub/flogo/2.6.1/doc/pdf/TIB_flogo_2.6_users_guide.pdf?id=2)
- Upload [`fabticExtension.zip`](../../fabricExtension.zip) to TIBCO Flogo® Enterprise [Extensions](http://localhost:8090/wistudio/extensions).  Note that you can generate this `zip` by using the script [`zip-fabric.sh`](../../zip-fabric.sh).
- Create new Flogo App of name `all_fabric` and choose `Import app` to import the model [`all_fabric.json`](all_fabric.json)
- You can then view or edit the trigger and activities in graphical modeler of the TIBCO Flogo® Enterprise.

## Build and deploy chaincode to Hyperledger Fabric

This model is for illustration only.  For building and deploying real chaincode, refer other samples, e.g., [marble](../marble).
