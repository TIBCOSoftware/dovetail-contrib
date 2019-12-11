// Copyright © 2018. TIBCO Software Inc.
//
// This file is subject to the license terms contained
// in the license file that is distributed with this file.

package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/protobuf/proto"
//	pb "github.com/hyperledger/fabric-protos-go/peer"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/core/chaincode/platforms/golang"
	"github.com/pkg/errors"
)

func packageCDS(path, name, version, outFile string) error {
	// verify path of format $GOPATH/src/cc
	var gopath, srcpath string
	abspath, err := filepath.Abs(path)
	if err != nil {
		return errors.Wrapf(err, "failed to get absolute path for %s", path)
	}
	srcindex := strings.LastIndex(filepath.ToSlash(abspath), "/src/")
	if srcindex > 0 {
		gopath = abspath[0:srcindex]
		srcpath = abspath[srcindex+5:]
	} else {
		return errors.Errorf("path '%s' does not contain folder 'src'", abspath)
	}
	os.Setenv("GOPATH", gopath)
	fmt.Printf("gopath %s srcpath %s\n", gopath, srcpath)

	// generate cds content
	input := &pb.ChaincodeInput{}
	spec := &pb.ChaincodeSpec{
		Type:        pb.ChaincodeSpec_GOLANG,
		ChaincodeId: &pb.ChaincodeID{Path: srcpath, Name: name, Version: version},
		Input:       input,
	}
	pr := &golang.Platform{}
	codePackageBytes, err := pr.GetDeploymentPayload(spec.ChaincodeId.Path)
	if err != nil {
		return errors.Wrapf(err, "failed to generate deployment payload from source %s", spec.ChaincodeId.Path)
	}

	// write cds file
	chaincodeDeploymentSpec := &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec, CodePackage: codePackageBytes}
	data, err := proto.Marshal(chaincodeDeploymentSpec)
	if err != nil {
		return errors.Wrap(err, "failed to marshal cds data")
	}
	fmt.Println("write cds file", outFile)
	return ioutil.WriteFile(outFile, data, 0700)
}
