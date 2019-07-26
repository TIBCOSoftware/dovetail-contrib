package main

import (
	"bytes"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
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
	"github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/ledger/rwset"
	"github.com/hyperledger/fabric/protos/msp"
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
	subEvent          = "CHAINCODE" // BLOCK, FILTERED, CHAINCODE
	eventFilter       = "PO|add"
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

	if subEvent == "BLOCK" {
		// register and wait for one block event
		registration, blkChan, err := client.RegisterBlockEvent()
		if err != nil {
			panic(errors.Wrapf(err, "Failed to register block event"))
		}
		defer client.Unregister(registration)

		fmt.Println("block event registered successfully")
		receiveBlockEvent(blkChan)
	} else if subEvent == "FILTERED" {
		// register and wait for one filtered block event
		registration, blkChan, err := client.RegisterFilteredBlockEvent()
		if err != nil {
			panic(errors.Wrapf(err, "Failed to register filtered block event"))
		}
		defer client.Unregister(registration)

		fmt.Println("filtered block event registered successfully")
		receiveFilteredBlockEvent(blkChan)
	} else if subEvent == "CHAINCODE" {
		// register and wait for one chaincode event
		registration, ccChan, err := client.RegisterChaincodeEvent("equinix_cc", eventFilter)
		if err != nil {
			panic(errors.Wrapf(err, "Failed to register chaincode event"))
		}
		defer client.Unregister(registration)

		fmt.Println("chaincode event registered successfully")
		receiveChaincodeEvent(ccChan)
	}
}

func receiveChaincodeEvent(ccChan <-chan *fab.CCEvent) {
	var ccEvent *fab.CCEvent
	select {
	case ccEvent = <-ccChan:
		cce := unmarshalChaincodeEvent(ccEvent)
		ccejson, err := json.Marshal(cce)
		if err != nil {
			fmt.Printf("got chaincode event %+v\n", cce)
		} else {
			fmt.Printf("got chaincode event %s\n", ccejson)
		}
	case <-time.After(time.Second * 3600):
		fmt.Println("Timeout waiting for chaincode event")
	}
}

func unmarshalChaincodeEvent(ccEvent *fab.CCEvent) *CCEventDetail {
	ced := CCEventDetail{
		BlockNumber: ccEvent.BlockNumber,
		SourceURL:   ccEvent.SourceURL,
		TxID:        ccEvent.TxID,
		ChaincodeID: ccEvent.ChaincodeID,
		EventName:   ccEvent.EventName,
		Payload:     string(ccEvent.Payload),
	}
	return &ced
}

func receiveFilteredBlockEvent(blkChan <-chan *fab.FilteredBlockEvent) {
	var blkEvent *fab.FilteredBlockEvent
	select {
	case blkEvent = <-blkChan:
		bed, err := unmarshalFilteredBlockEvent(blkEvent)
		if err != nil {
			panic(err)
		}
		blk, err := json.Marshal(bed)
		if err != nil {
			fmt.Printf("got filtered block %+v\n", bed)
		} else {
			fmt.Printf("got filtered block %s\n", blk)
		}
	case <-time.After(time.Second * 3600):
		fmt.Println("Timeout waiting for filtered block event")
	}
}

func unmarshalFilteredBlockEvent(blkEvent *fab.FilteredBlockEvent) (*BlockEventDetail, error) {
	blk := blkEvent.FilteredBlock
	//	blkjson, _ := json.Marshal(blk)
	//	fmt.Println(string(blkjson))

	bed := BlockEventDetail{
		SourceURL:    blkEvent.SourceURL,
		Number:       blk.Number,
		Transactions: []*TransactionDetail{},
	}

	for _, d := range blk.FilteredTransactions {
		td := TransactionDetail{
			TxType:    common.HeaderType_name[int32(d.Type)],
			TxID:      d.Txid,
			ChannelID: blk.ChannelId,
			Actions:   []*ActionDetail{},
		}
		bed.Transactions = append(bed.Transactions, &td)
		actions := d.GetTransactionActions()
		if actions != nil {
			for _, ta := range actions.ChaincodeActions {
				ce := ta.GetChaincodeEvent()
				if ce != nil && ce.ChaincodeId != "" {
					action := ActionDetail{
						Chaincode: &ChaincodeID{Name: ce.ChaincodeId},
						Result: &ChaincodeResult{
							Event: &ChaincodeEvent{
								Name:    ce.EventName,
								Payload: string(ce.Payload),
							},
						},
					}
					td.Actions = append(td.Actions, &action)
				}
			}
		}
	}
	return &bed, nil
}

func receiveBlockEvent(blkChan <-chan *fab.BlockEvent) {
	var blkEvent *fab.BlockEvent
	select {
	case blkEvent = <-blkChan:
		bed, err := unmarshalBlockEvent(blkEvent)
		if err != nil {
			panic(err)
		}
		blk, err := json.Marshal(bed)
		if err != nil {
			fmt.Printf("got block %+v\n", bed)
		} else {
			fmt.Printf("got block %s\n", blk)
		}
	case <-time.After(time.Second * 3600):
		fmt.Println("Timeout waiting for block event")
	}
}

func unmarshalBlockEvent(blkEvent *fab.BlockEvent) (*BlockEventDetail, error) {
	bed := BlockEventDetail{
		SourceURL:    blkEvent.SourceURL,
		Number:       blkEvent.Block.Header.Number,
		Transactions: []*TransactionDetail{},
	}
	for _, d := range blkEvent.Block.Data.Data {
		txn, err := unmarshalTransaction(d)
		if err != nil {
			fmt.Printf("Error unmarshalling transaction: %+v\n", err)
			continue
		} else {
			bed.Transactions = append(bed.Transactions, txn)
		}
	}
	return &bed, nil
}

func unmarshalTransaction(data []byte) (*TransactionDetail, error) {
	envelope, err := utils.GetEnvelopeFromBlock(data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get envelope")
	}
	payload, err := utils.GetPayload(envelope)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get payload")
	}
	if payload.Header == nil {
		return nil, errors.Errorf("payload header is empty")
	}

	// channel header
	chdr, err := utils.UnmarshalChannelHeader(payload.Header.ChannelHeader)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal channel header")
	}
	td := TransactionDetail{
		TxType:    common.HeaderType_name[chdr.Type],
		TxID:      chdr.TxId,
		TxTime:    time.Unix(chdr.Timestamp.Seconds, int64(chdr.Timestamp.Nanos)).UTC().String(),
		ChannelID: chdr.ChannelId,
		Actions:   []*ActionDetail{},
	}

	// signature header
	shdr := &common.SignatureHeader{}
	if err = proto.Unmarshal(payload.Header.SignatureHeader, shdr); err != nil {
		fmt.Printf("failed to unmarshal signature header: %+v\n", err)
	} else {
		cid, err := unmarshalIdentity(shdr.Creator)
		if err != nil {
			fmt.Printf("failed to unmarshal creator identity: %+v\n", err)
		} else {
			td.CreatorIdentity = cid
		}
	}

	txn, err := utils.GetTransaction(payload.Data)
	if err != nil {
		return &td, errors.Wrapf(err, "failed to get transaction")
	}
	for _, ta := range txn.Actions {
		act, err := unmarshalAction(ta.Payload)
		if err != nil {
			fmt.Printf("Error unmarshalling action: %+v\n", err)
			continue
		} else {
			td.Actions = append(td.Actions, act)
		}
	}
	return &td, nil
}

func unmarshalIdentity(data []byte) (*Identity, error) {
	cid := &msp.SerializedIdentity{}
	if err := proto.Unmarshal(data, cid); err != nil {
		return nil, err
	}
	id := Identity{Mspid: cid.Mspid, Cert: string(cid.IdBytes)}

	// extract info from x509 certificate
	block, _ := pem.Decode(cid.IdBytes)
	if block == nil {
		fmt.Println("creator certificate is empty")
		return &id, nil
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		fmt.Printf("failed to parse creator certificate: %+v\n", err)
		return &id, nil
	}
	id.Subject = cert.Subject.CommonName
	id.Issuer = cert.Issuer.CommonName
	return &id, nil
}

func unmarshalAction(data []byte) (*ActionDetail, error) {
	ccAction, err := utils.GetChaincodeActionPayload(data)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get action payload")
	}

	// proposal payload
	proposalPayload, err := utils.GetChaincodeProposalPayload(ccAction.ChaincodeProposalPayload)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get proposal payload")
	}
	cis := &pb.ChaincodeInvocationSpec{}
	err = proto.Unmarshal(proposalPayload.Input, cis)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal chaincode input")
	}

	// chaincode spec
	ccID := ChaincodeID{
		Type: pb.ChaincodeSpec_Type_name[int32(cis.ChaincodeSpec.Type)],
		Name: cis.ChaincodeSpec.ChaincodeId.Name,
	}

	// chaincode input
	args := cis.ChaincodeSpec.Input.Args
	ccInput := ChaincodeInput{
		Function: string(args[0]),
		Args:     []string{},
	}
	if len(args) > 1 {
		for _, arg := range args[1:] {
			ccInput.Args = append(ccInput.Args, string(arg))
		}
	}
	if proposalPayload.TransientMap != nil {
		tm := make(map[string]string)
		for k, v := range proposalPayload.TransientMap {
			tm[k] = string(v)
		}
		if tb, err := json.Marshal(tm); err != nil {
			fmt.Printf("failed to marshal transient map to JSON: %+v\n", err)
		} else {
			ccInput.TransientMap = string(tb)
		}
	}

	// action response payload
	prespPayload, err := utils.GetProposalResponsePayload(ccAction.Action.ProposalResponsePayload)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get proposal response payload")
	}
	cact, err := utils.GetChaincodeAction(prespPayload.Extension)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get chaincode action")
	}
	if cact.Response == nil {
		return nil, errors.New("chaincode response is empty")
	}
	if cact.ChaincodeId != nil {
		ccID.Version = cact.ChaincodeId.Version
	}

	// chaincode response
	ccResult := ChaincodeResult{
		Response: &ChaincodeResponse{
			Status:  cact.Response.Status,
			Message: cact.Response.Message,
			Payload: string(cact.Response.Payload),
		},
	}
	action := ActionDetail{
		Chaincode:     &ccID,
		Input:         &ccInput,
		Result:        &ccResult,
		EndorserCount: len(ccAction.Action.Endorsements),
	}

	// chaincode result
	if cact.Results != nil {
		txrw := &rwset.TxReadWriteSet{}
		err := proto.Unmarshal(cact.Results, txrw)
		if err != nil {
			fmt.Printf("failed to unmarshal tx rwset: %+v\n", err)
		} else {
			ccResult.ReadWriteCount = len(txrw.NsRwset)
		}
	}

	// chaincode event
	if cact.Events != nil {
		if ccEvt, err := utils.GetChaincodeEvents(cact.Events); err != nil {
			fmt.Printf("failed to get chaincode event: %+v\n", err)
		} else {
			ccResult.Event = &ChaincodeEvent{
				Name:    ccEvt.EventName,
				Payload: string(ccEvt.Payload),
			}
		}
	}
	return &action, nil
}

type CCEventDetail struct {
	BlockNumber uint64 `json:"block"`
	SourceURL   string `json:"source,omitempty"`
	TxID        string `json:"txId"`
	ChaincodeID string `json:"chaincode"`
	EventName   string `json:"name"`
	Payload     string `json:"payload"`
}

// BlockEventDetail contains data in a block event
type BlockEventDetail struct {
	Number       uint64               `json:"block"`
	SourceURL    string               `json:"source,omitempty"`
	Transactions []*TransactionDetail `json:"transactions"`
}

// TransactionDetail contains data in a transaction
type TransactionDetail struct {
	TxType          string          `json:"type"`
	TxID            string          `json:"txId"`
	TxTime          string          `json:"txTime,omitempty"`
	ChannelID       string          `json:"channel"`
	CreatorIdentity *Identity       `json:"creator,omitempty"`
	Actions         []*ActionDetail `json:"actions,omitempty"`
}

// Identity contains creator's mspid and certificate
type Identity struct {
	Mspid   string `json:"mspid"`
	Subject string `json:"subject,omitempty"`
	Issuer  string `json:"issuer,omitempty"`
	Cert    string `json:"cert,omitempty"`
}

// ActionDetail contains data in a chaincode invocation
type ActionDetail struct {
	Chaincode     *ChaincodeID     `json:"chaincode,omitempty"`
	Input         *ChaincodeInput  `json:"input,omitempty"`
	Result        *ChaincodeResult `json:"result,omitempty"`
	EndorserCount int              `json:"endorsers,omitempty"`
}

// ChaincodeID defines chaincode identity
type ChaincodeID struct {
	Type    string `json:"type,omitempty"`
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

// ChaincodeInput defines input parameters of a chaincode invocation
type ChaincodeInput struct {
	Function     string   `json:"function"`
	Args         []string `json:"args,omitempty"`
	TransientMap string   `json:"transient,omitempty"`
}

// ChaincodeResult defines result of a chaincode invocation
type ChaincodeResult struct {
	ReadWriteCount int                `json:"rwset,omitempty"`
	Response       *ChaincodeResponse `json:"response,omitempty"`
	Event          *ChaincodeEvent    `json:"event,omitempty"`
}

// ChaincodeResponse defines response from a chaincode invocation
type ChaincodeResponse struct {
	Status  int32  `json:"status"`
	Message string `json:"message,omitempty"`
	Payload string `json:"payload,omitempty"`
}

// ChaincodeEvent defines event created by a chaincode invocation
type ChaincodeEvent struct {
	Name    string `json:"name"`
	Payload string `json:"payload,omitempty"`
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
