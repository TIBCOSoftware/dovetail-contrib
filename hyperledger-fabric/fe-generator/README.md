# fe-generator

As of Flogo Enterprise 2.8, it does not support build with Go modules. This script here generates necessary files in a Flogo Enterprise installation, so it can support Go modules.

Execute the following script:

```bash
FE_HOME=/path/to/flogo/2.8
cd ./fe-generator
./init-gomod.sh ${FE_HOME}
```

Only the following folders are needed to build Flogo models that use Flogo Enterprise components, so you can keep only these folders on the build server/container:

- flogo/2.8/lib/core/
- flogo/data/localstack/wicontributions/Tibco/
