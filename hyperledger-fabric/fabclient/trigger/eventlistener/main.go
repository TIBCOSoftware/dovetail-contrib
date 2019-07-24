package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/event"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/core"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric/protos/ledger/rwset"
	pb "github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/protos/utils"
)

const (
	netConfigPath     = "${GOPATH}/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/testdata/config_min.yaml"
	entityMatcherPath = "${GOPATH}/src/github.com/TIBCOSoftware/dovetail-contrib/hyperledger-fabric/testdata/local_entity_matchers.yaml"
	user              = "User1"
	org               = "org1"
	channelID         = "mychannel"
	cryptoPath        = "/Users/yxu/go/src/github.com/hyperledger/fabric-samples/first-network/crypto-config"
)

func main() {
	os.Setenv("CRYPTO_PATH", cryptoPath)
	netConfig, _ := ReadFile(netConfigPath)
	entityMatcher, _ := ReadFile(entityMatcherPath)
	sdk, err := fabsdk.New(networkConfigProvider(netConfig, entityMatcher))
	if err != nil {
		panic(errors.Wrapf(err, "Failed to create new SDK"))
	}

	opts := []fabsdk.ContextOption{fabsdk.WithUser(user)}
	if org != "" {
		opts = append(opts, fabsdk.WithOrg(org))
	}

	client, err := event.New(sdk.ChannelContext(channelID, opts...), event.WithBlockEvents())
	if err != nil {
		panic(errors.Wrapf(err, "Failed to create event client"))
	}

	registration, blkChan, err := client.RegisterBlockEvent()
	if err != nil {
		panic(errors.Wrapf(err, "Failed to register block event"))
	}
	defer client.Unregister(registration)

	fmt.Println("block event registered successfully")

	var blkEvent *fab.BlockEvent
	select {
	case blkEvent = <-blkChan:
		// block number
		fmt.Printf("Received block peer-URL %s, number %d\n", blkEvent.SourceURL, blkEvent.Block.Header.GetNumber())
		for _, d := range blkEvent.Block.Data.Data {
			envelope, err := utils.GetEnvelopeFromBlock(d)
			if err != nil {
				panic(errors.Wrapf(err, "failed to get envelope"))
			}
			payload, err := utils.GetPayload(envelope)
			if err != nil {
				panic(errors.Wrapf(err, "failed to get payload"))
			}

			// channel header
			if payload.Header == nil {
				panic(errors.Errorf("payload header is empty"))
			}
			chdr, err := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
			if err != nil {
				panic(errors.Wrapf(err, "failed to unmarshal channel header"))
			}
			cheader, err := json.Marshal(chdr)
			if err != nil {
				fmt.Printf("failed to marshal channel header to json: %+v\n", err)
			}
			fmt.Printf("channel header: %s\n", string(cheader))
			fmt.Printf("channel id %s, txID: %s, timestamp: %s\n", chdr.ChannelId, chdr.TxId, chdr.Timestamp.String())

			txn, err := utils.GetTransaction(payload.Data)
			if err != nil {
				panic(errors.Wrapf(err, "failed to get transaction"))
			}
			for _, t := range txn.Actions {
				// transaction payload
				ccAction, err := utils.GetChaincodeActionPayload(t.Payload)
				if err != nil {
					panic(errors.Wrapf(err, "failed to get action payload"))
				}
				proposalPayload, err := utils.GetChaincodeProposalPayload(ccAction.ChaincodeProposalPayload)
				if err != nil {
					panic(errors.Wrapf(err, "failed to get proposal payload"))
				}
				cis := &pb.ChaincodeInvocationSpec{}
				err = proto.Unmarshal(proposalPayload.Input, cis)
				if err != nil {
					fmt.Printf("failed to unmarshal chaincode input: %+v\n", err)
				}
				ccjson, err := json.Marshal(cis.ChaincodeSpec)
				if err != nil {
					fmt.Printf("failed to marshal chaincode spec to json: %+v\n", err)
				}
				fmt.Printf("chaincode spec: %s\n", string(ccjson))
				fmt.Print("input args: ")
				for _, arg := range cis.ChaincodeSpec.Input.Args {
					fmt.Print(string(arg) + ", ")
				}
				fmt.Println("")
				// transient map: proposalPayload.TransientMap

				prespPayload, err := utils.GetProposalResponsePayload(ccAction.Action.ProposalResponsePayload)
				if err != nil {
					panic(errors.Wrapf(err, "failed to get proposal response payload"))
				}
				cact, err := utils.GetChaincodeAction(prespPayload.Extension)
				if err != nil {
					panic(errors.Wrapf(err, "failed to get chaincode action"))
				}
				if cact.Events != nil {
					ccEvt, err := utils.GetChaincodeEvents(cact.Events)
					if err != nil {
						fmt.Printf("failed to get chaincode event: %+v\n", err)
					}
					if ccevtjson, err := json.Marshal(ccEvt); err != nil {
						fmt.Printf("failed to marshal chaincode event to json: %+v\n", err)
					} else {
						fmt.Printf("chaincode event: %s\n", ccevtjson)
						fmt.Printf("event payload: %s\n", string(ccEvt.Payload))
					}
				}
				if cact.Response != nil {
					fmt.Printf("response status: %d message: %s payload %s\n", cact.Response.Status, cact.Response.Message, string(cact.Response.Payload))
				}
				if cact.Results != nil {
					txrw := &rwset.TxReadWriteSet{}
					err := proto.Unmarshal(cact.Results, txrw)
					if err != nil {
						fmt.Printf("failed to unmarshal tx rwset: %+v\n", err)
					} else {
						fmt.Printf("Number of rwsets: %d\n", len(txrw.NsRwset))
					}
				}
				if ccid, err := json.Marshal(cact.ChaincodeId); err != nil {
					fmt.Printf("failed to marshal chaincode id to json: %+v\n", err)
				} else {
					fmt.Printf("chaincode id: %s\n", ccid)
				}
			}
		}
	case <-time.After(time.Second * 3600):
		fmt.Println("Timeout waiting for block event")
	}
}

func networkConfigProvider(networkConfig []byte, entityMatcherOverride []byte) core.ConfigProvider {
	configProvider := config.FromRaw(networkConfig, "yaml")

	if entityMatcherOverride != nil {
		return func() ([]core.ConfigBackend, error) {
			matcherProvider := config.FromRaw(entityMatcherOverride, "yaml")
			matcherBackends, err := matcherProvider()
			if err != nil {
				fmt.Printf("failed to parse entity matchers: %+v\n", err)
				// return the original config provider defined by configPath
				return configProvider()
			}

			currentBackends, err := configProvider()
			if err != nil {
				fmt.Printf("failed to parse network config: %+v\n", err)
				return nil, err
			}

			// return the combined config with matcher precedency
			return append(matcherBackends, currentBackends...), nil
		}
	}
	return configProvider
}

// ReadFile returns content of a specified file
func ReadFile(filePath string) ([]byte, error) {
	f, err := os.Open(Subst(filePath))
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to open file: %s", filePath)
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read file stat: %s", filePath)
	}
	s := fi.Size()
	cBytes := make([]byte, s)
	n, err := f.Read(cBytes)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read file: %s", filePath)
	}
	if n == 0 {
		fmt.Printf("file %s is empty\n", filePath)
	}
	return cBytes, err
}

// Subst replaces instances of '${VARNAME}' (eg ${GOPATH}) with the variable.
// Variables names that are not set by the SDK are replaced with the environment variable.
func Subst(path string) string {
	const (
		sepPrefix = "${"
		sepSuffix = "}"
	)

	splits := strings.Split(path, sepPrefix)

	var buffer bytes.Buffer

	// first split precedes the first sepPrefix so should always be written
	buffer.WriteString(splits[0]) // nolint: gas

	for _, s := range splits[1:] {
		subst, rest := substVar(s, sepPrefix, sepSuffix)
		buffer.WriteString(subst) // nolint: gas
		buffer.WriteString(rest)  // nolint: gas
	}

	return buffer.String()
}

// substVar searches for an instance of a variables name and replaces them with their value.
// The first return value is substituted portion of the string or noMatch if no replacement occurred.
// The second return value is the unconsumed portion of s.
func substVar(s string, noMatch string, sep string) (string, string) {
	endPos := strings.Index(s, sep)
	if endPos == -1 {
		return noMatch, s
	}

	v, ok := os.LookupEnv(s[:endPos])
	if !ok {
		return noMatch, s
	}

	return v, s[endPos+1:]
}
