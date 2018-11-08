package runtime

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/TIBCOSoftware/flogo-lib/app"
	"github.com/stretchr/testify/assert"
)

func TestExplicitReply(t *testing.T) {
	model := `{
		"name": "IOUDemo",
		"description": " ",
		"version": "1.0.0",
		"type": "flogo:app",
		"appModel": "1.0.0",
		"resources": [
		 {
		  "id": "flow:IssueIOU",
		  "data": {
		   "name": "IssueIOU",
		   "description": "",
		   "tasks": [
			{
			 "id": "Logger",
			 "name": "Logger",
			 "activity": {
			  "ref": "github.com/TIBCOSoftware/dovetail-contrib/SmartContract/activity/logger",
			  "input": {
			   "logLevel": "INFO"
			  },
			  "output": {},
			  "mappings": {
			   "input": [
				{
				 "mapTo": "$INPUT['message']",
				 "type": "expression",
				 "value": "issue IOU flow started..."
				},
				{
				 "mapTo": "$INPUT['containerServiceStub']",
				 "type": "expression",
				 "value": "$flow.containerServiceStub"
				}
			   ]
			  }
			 }
			},
			{
			 "id": "LedgerOperation",
			 "name": "LedgerOperation",
			 "activity": {
			  "ref": "github.com/TIBCOSoftware/dovetail-contrib/SmartContract/activity/ledger",
			  "input": {
			   "model": "6c06b8e0-dbba-11e8-9620-03a6d4bbf53c",
			   "assetName": "com.tibco.cp.IOU",
			   "operation": "PUT",
			   "input": {
				"metadata": "",
				"value": ""
			   },
			   "identifier": "linearId"
			  },
			  "output": {
			   "output": {
				"metadata": "{\"$schema\":\"http://json-schema.org/draft-04/schema#\",\"title\":\"IOU\",\"type\":\"object\",\"properties\":{\"lender\":{\"type\":\"string\"},\"borrower\":{\"type\":\"string\"},\"amt\":{\"type\":\"object\",\"properties\":{\"currency\":{\"type\":\"string\"},\"quantity\":{\"type\":\"integer\",\"default\":\"0\"}},\"required\":[\"currency\",\"quantity\"]},\"paid\":{\"type\":\"object\",\"properties\":{\"currency\":{\"type\":\"string\"},\"quantity\":{\"type\":\"integer\",\"default\":\"0\"}},\"required\":[\"currency\",\"quantity\"]},\"linearId\":{\"type\":\"string\"}},\"required\":[\"lender\",\"borrower\",\"amt\",\"paid\",\"linearId\"],\"description\":\"{\\\"metadata\\\":{\\\"type\\\":\\\"Asset\\\",\\\"parent\\\":\\\"com.tibco.dovetail.system.LinearState\\\",\\\"isAbstract\\\":false,\\\"identifiedBy\\\":\\\"linearId\\\",\\\"decorators\\\":[]},\\\"attributes\\\":[{\\\"name\\\":\\\"lender\\\",\\\"type\\\":\\\"com.tibco.dovetail.system.Party\\\",\\\"isRef\\\":true},{\\\"name\\\":\\\"borrower\\\",\\\"type\\\":\\\"com.tibco.dovetail.system.Party\\\",\\\"isRef\\\":true},{\\\"name\\\":\\\"amt\\\",\\\"isOptional\\\":false,\\\"type\\\":\\\"com.tibco.dovetail.system.Amount\\\"},{\\\"name\\\":\\\"paid\\\",\\\"isOptional\\\":false,\\\"type\\\":\\\"com.tibco.dovetail.system.Amount\\\"},{\\\"name\\\":\\\"linearId\\\",\\\"isOptional\\\":false,\\\"type\\\":\\\"String\\\"}]}\"}",
				"value": ""
			   }
			  },
			  "mappings": {
			   "input": [
				{
				 "mapTo": "$INPUT['input']",
				 "type": "expression",
				 "value": "$flow.transactionInput.iou"
				},
				{
				 "mapTo": "$INPUT['containerServiceStub']",
				 "type": "expression",
				 "value": "$flow.containerServiceStub"
				}
			   ]
			  }
			 }
			},
			{
			 "id": "TransactionResponse",
			 "name": "TransactionResponse",
			 "activity": {
			  "ref": "github.com/TIBCOSoftware/dovetail-contrib/SmartContract/activity/txnreply",
			  "input": {
			   "status": "Success"
			  },
			  "output": {}
			 }
			}
		   ],
		   "links": [
			{
			 "id": 1,
			 "from": "Logger",
			 "to": "LedgerOperation",
			 "type": "default"
			},
			{
			 "id": 2,
			 "from": "LedgerOperation",
			 "to": "TransactionResponse",
			 "type": "default"
			}
		   ],
		   "metadata": {
			"input": [],
			"output": []
		   }
		  }
		 }
		],
		"triggers": [
		 {
		  "ref": "github.com/TIBCOSoftware/dovetail-contrib/SmartContract/trigger/transaction",
		  "name": "SmartContractTXNTrigger",
		  "description": "",
		  "settings": {
		   "model": "6c06b8e0-dbba-11e8-9620-03a6d4bbf53c",
		   "createAll": false,
		   "assets": "[\"com.tibco.dovetail.system.Cash\",\"com.tibco.cp.IOU\"]",
		   "transactions": "[\"com.tibco.cp.IssueIOU\",\"com.tibco.cp.TransferIOU\",\"com.tibco.cp.SettleIOU\"]",
		   "concepts": "[\"com.tibco.dovetail.system.Amount\",\"com.tibco.dovetail.system.IssueAmount\",\"com.tibco.dovetail.system.TimeWindow\"]",
		   "schemas": "[[\"com.tibco.dovetail.system.LinearState\",\"{\\\"$schema\\\":\\\"http://json-schema.org/draft-04/schema#\\\",\\\"title\\\":\\\"LinearState\\\",\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"linearId\\\":{\\\"type\\\":\\\"string\\\"}},\\\"required\\\":[\\\"linearId\\\"],\\\"description\\\":\\\"{\\\\\\\"metadata\\\\\\\":{\\\\\\\"type\\\\\\\":\\\\\\\"Asset\\\\\\\",\\\\\\\"parent\\\\\\\":\\\\\\\"org.hyperledger.composer.system.Asset\\\\\\\",\\\\\\\"isAbstract\\\\\\\":true,\\\\\\\"decorators\\\\\\\":[{\\\\\\\"name\\\\\\\":\\\\\\\"CordaClass\\\\\\\",\\\\\\\"args\\\\\\\":[\\\\\\\"net.corda.core.contracts.LinearState\\\\\\\"]}],\\\\\\\"cordaClass\\\\\\\":\\\\\\\"net.corda.core.contracts.LinearState\\\\\\\"},\\\\\\\"attributes\\\\\\\":[{\\\\\\\"name\\\\\\\":\\\\\\\"linearId\\\\\\\",\\\\\\\"isOptional\\\\\\\":false,\\\\\\\"type\\\\\\\":\\\\\\\"String\\\\\\\"}]}\\\"}\"],[\"com.tibco.dovetail.system.OwnableState\",\"{\\\"$schema\\\":\\\"http://json-schema.org/draft-04/schema#\\\",\\\"title\\\":\\\"OwnableState\\\",\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"linearId\\\":{\\\"type\\\":\\\"string\\\"},\\\"owner\\\":{\\\"type\\\":\\\"string\\\"}},\\\"required\\\":[\\\"linearId\\\",\\\"owner\\\"],\\\"description\\\":\\\"{\\\\\\\"metadata\\\\\\\":{\\\\\\\"type\\\\\\\":\\\\\\\"Asset\\\\\\\",\\\\\\\"parent\\\\\\\":\\\\\\\"org.hyperledger.composer.system.Asset\\\\\\\",\\\\\\\"isAbstract\\\\\\\":true,\\\\\\\"decorators\\\\\\\":[{\\\\\\\"name\\\\\\\":\\\\\\\"CordaClass\\\\\\\",\\\\\\\"args\\\\\\\":[\\\\\\\"net.corda.core.contracts.OwnableState\\\\\\\"]}],\\\\\\\"cordaClass\\\\\\\":\\\\\\\"net.corda.core.contracts.OwnableState\\\\\\\"},\\\\\\\"attributes\\\\\\\":[{\\\\\\\"name\\\\\\\":\\\\\\\"linearId\\\\\\\",\\\\\\\"isOptional\\\\\\\":false,\\\\\\\"type\\\\\\\":\\\\\\\"String\\\\\\\"},{\\\\\\\"name\\\\\\\":\\\\\\\"owner\\\\\\\",\\\\\\\"type\\\\\\\":\\\\\\\"com.tibco.dovetail.system.Party\\\\\\\",\\\\\\\"isRef\\\\\\\":true}]}\\\"}\"],[\"com.tibco.dovetail.system.FungibleAsset\",\"{\\\"$schema\\\":\\\"http://json-schema.org/draft-04/schema#\\\",\\\"title\\\":\\\"FungibleAsset\\\",\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"owner\\\":{\\\"type\\\":\\\"string\\\"},\\\"amt\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"currency\\\":{\\\"type\\\":\\\"string\\\"},\\\"quantity\\\":{\\\"type\\\":\\\"integer\\\",\\\"default\\\":\\\"0\\\"}},\\\"required\\\":[\\\"currency\\\",\\\"quantity\\\"]}},\\\"required\\\":[\\\"owner\\\",\\\"amt\\\"],\\\"description\\\":\\\"{\\\\\\\"metadata\\\\\\\":{\\\\\\\"type\\\\\\\":\\\\\\\"Asset\\\\\\\",\\\\\\\"parent\\\\\\\":\\\\\\\"org.hyperledger.composer.system.Asset\\\\\\\",\\\\\\\"isAbstract\\\\\\\":true,\\\\\\\"decorators\\\\\\\":[{\\\\\\\"name\\\\\\\":\\\\\\\"CordaClass\\\\\\\",\\\\\\\"args\\\\\\\":[\\\\\\\"net.corda.core.contracts.FungibleAsset\\\\\\\"]}],\\\\\\\"cordaClass\\\\\\\":\\\\\\\"net.corda.core.contracts.FungibleAsset\\\\\\\"},\\\\\\\"attributes\\\\\\\":[{\\\\\\\"name\\\\\\\":\\\\\\\"owner\\\\\\\",\\\\\\\"type\\\\\\\":\\\\\\\"com.tibco.dovetail.system.Party\\\\\\\",\\\\\\\"isRef\\\\\\\":true},{\\\\\\\"name\\\\\\\":\\\\\\\"amt\\\\\\\",\\\\\\\"isOptional\\\\\\\":false,\\\\\\\"type\\\\\\\":\\\\\\\"com.tibco.dovetail.system.Amount\\\\\\\"}]}\\\"}\"],[\"com.tibco.dovetail.system.Cash\",\"{\\\"$schema\\\":\\\"http://json-schema.org/draft-04/schema#\\\",\\\"title\\\":\\\"Cash\\\",\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"assetId\\\":{\\\"type\\\":\\\"string\\\"},\\\"owner\\\":{\\\"type\\\":\\\"string\\\"},\\\"amt\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"currency\\\":{\\\"type\\\":\\\"string\\\"},\\\"quantity\\\":{\\\"type\\\":\\\"integer\\\",\\\"default\\\":\\\"0\\\"}},\\\"required\\\":[\\\"currency\\\",\\\"quantity\\\"]}},\\\"required\\\":[\\\"assetId\\\",\\\"owner\\\",\\\"amt\\\"],\\\"description\\\":\\\"{\\\\\\\"metadata\\\\\\\":{\\\\\\\"type\\\\\\\":\\\\\\\"Asset\\\\\\\",\\\\\\\"parent\\\\\\\":\\\\\\\"com.tibco.dovetail.system.FungibleAsset\\\\\\\",\\\\\\\"isAbstract\\\\\\\":false,\\\\\\\"identifiedBy\\\\\\\":\\\\\\\"assetId\\\\\\\",\\\\\\\"decorators\\\\\\\":[{\\\\\\\"name\\\\\\\":\\\\\\\"CordaClass\\\\\\\",\\\\\\\"args\\\\\\\":[\\\\\\\"net.corda.finance.contracts.asset.Cash.State\\\\\\\"]}],\\\\\\\"cordaClass\\\\\\\":\\\\\\\"net.corda.finance.contracts.asset.Cash.State\\\\\\\"},\\\\\\\"attributes\\\\\\\":[{\\\\\\\"name\\\\\\\":\\\\\\\"assetId\\\\\\\",\\\\\\\"isOptional\\\\\\\":false,\\\\\\\"type\\\\\\\":\\\\\\\"String\\\\\\\"},{\\\\\\\"name\\\\\\\":\\\\\\\"owner\\\\\\\",\\\\\\\"type\\\\\\\":\\\\\\\"com.tibco.dovetail.system.Party\\\\\\\",\\\\\\\"isRef\\\\\\\":true},{\\\\\\\"name\\\\\\\":\\\\\\\"amt\\\\\\\",\\\\\\\"isOptional\\\\\\\":false,\\\\\\\"type\\\\\\\":\\\\\\\"com.tibco.dovetail.system.Amount\\\\\\\"}]}\\\"}\"],[\"com.tibco.dovetail.system.Amount\",\"{\\\"$schema\\\":\\\"http://json-schema.org/draft-04/schema#\\\",\\\"title\\\":\\\"Amount\\\",\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"currency\\\":{\\\"type\\\":\\\"string\\\"},\\\"quantity\\\":{\\\"type\\\":\\\"integer\\\",\\\"default\\\":\\\"0\\\"}},\\\"required\\\":[\\\"currency\\\",\\\"quantity\\\"],\\\"description\\\":\\\"{\\\\\\\"metadata\\\\\\\":{\\\\\\\"type\\\\\\\":\\\\\\\"Concept\\\\\\\",\\\\\\\"isAbstract\\\\\\\":false,\\\\\\\"decorators\\\\\\\":[{\\\\\\\"name\\\\\\\":\\\\\\\"CordaClass\\\\\\\",\\\\\\\"args\\\\\\\":[\\\\\\\"net.corda.core.contracts.Amount<Currency>\\\\\\\"]}],\\\\\\\"cordaClass\\\\\\\":\\\\\\\"net.corda.core.contracts.Amount<Currency>\\\\\\\"},\\\\\\\"attributes\\\\\\\":[{\\\\\\\"name\\\\\\\":\\\\\\\"currency\\\\\\\",\\\\\\\"isOptional\\\\\\\":false,\\\\\\\"type\\\\\\\":\\\\\\\"String\\\\\\\"},{\\\\\\\"name\\\\\\\":\\\\\\\"quantity\\\\\\\",\\\\\\\"isOptional\\\\\\\":false,\\\\\\\"type\\\\\\\":\\\\\\\"Long\\\\\\\",\\\\\\\"defaultValue\\\\\\\":\\\\\\\"0\\\\\\\"}]}\\\"}\"],[\"com.tibco.dovetail.system.IssueAmount\",\"{\\\"$schema\\\":\\\"http://json-schema.org/draft-04/schema#\\\",\\\"title\\\":\\\"IssueAmount\\\",\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"issuer\\\":{\\\"type\\\":\\\"string\\\"},\\\"currency\\\":{\\\"type\\\":\\\"string\\\"},\\\"quantity\\\":{\\\"type\\\":\\\"integer\\\",\\\"default\\\":\\\"0\\\"}},\\\"required\\\":[\\\"issuer\\\",\\\"currency\\\",\\\"quantity\\\"],\\\"description\\\":\\\"{\\\\\\\"metadata\\\\\\\":{\\\\\\\"type\\\\\\\":\\\\\\\"Concept\\\\\\\",\\\\\\\"isAbstract\\\\\\\":false,\\\\\\\"decorators\\\\\\\":[{\\\\\\\"name\\\\\\\":\\\\\\\"CordaClass\\\\\\\",\\\\\\\"args\\\\\\\":[\\\\\\\"net.corda.core.contracts.Amount<Issue<Currency>>\\\\\\\"]}],\\\\\\\"cordaClass\\\\\\\":\\\\\\\"net.corda.core.contracts.Amount<Issue<Currency>>\\\\\\\"},\\\\\\\"attributes\\\\\\\":[{\\\\\\\"name\\\\\\\":\\\\\\\"issuer\\\\\\\",\\\\\\\"type\\\\\\\":\\\\\\\"com.tibco.dovetail.system.Party\\\\\\\",\\\\\\\"isRef\\\\\\\":true},{\\\\\\\"name\\\\\\\":\\\\\\\"currency\\\\\\\",\\\\\\\"isOptional\\\\\\\":false,\\\\\\\"type\\\\\\\":\\\\\\\"String\\\\\\\"},{\\\\\\\"name\\\\\\\":\\\\\\\"quantity\\\\\\\",\\\\\\\"isOptional\\\\\\\":false,\\\\\\\"type\\\\\\\":\\\\\\\"Long\\\\\\\",\\\\\\\"defaultValue\\\\\\\":\\\\\\\"0\\\\\\\"}]}\\\"}\"],[\"com.tibco.dovetail.system.TimeWindow\",\"{\\\"$schema\\\":\\\"http://json-schema.org/draft-04/schema#\\\",\\\"title\\\":\\\"TimeWindow\\\",\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"from\\\":{\\\"format\\\":\\\"date-time\\\",\\\"type\\\":\\\"string\\\"},\\\"until\\\":{\\\"format\\\":\\\"date-time\\\",\\\"type\\\":\\\"string\\\"}},\\\"required\\\":[],\\\"description\\\":\\\"{\\\\\\\"metadata\\\\\\\":{\\\\\\\"type\\\\\\\":\\\\\\\"Concept\\\\\\\",\\\\\\\"isAbstract\\\\\\\":false,\\\\\\\"decorators\\\\\\\":[]},\\\\\\\"attributes\\\\\\\":[{\\\\\\\"name\\\\\\\":\\\\\\\"from\\\\\\\",\\\\\\\"isOptional\\\\\\\":true,\\\\\\\"type\\\\\\\":\\\\\\\"DateTime\\\\\\\"},{\\\\\\\"name\\\\\\\":\\\\\\\"until\\\\\\\",\\\\\\\"isOptional\\\\\\\":true,\\\\\\\"type\\\\\\\":\\\\\\\"DateTime\\\\\\\"}]}\\\"}\"],[\"com.tibco.dovetail.system.Party\",\"{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"id\\\":{\\\"type\\\":\\\"string\\\"}},\\\"required\\\":[\\\"id\\\"],\\\"description\\\":\\\"{\\\\\\\"metadata\\\\\\\":{\\\\\\\"type\\\\\\\":\\\\\\\"Participant\\\\\\\",\\\\\\\"parent\\\\\\\":\\\\\\\"org.hyperledger.composer.system.Participant\\\\\\\",\\\\\\\"isAbstract\\\\\\\":false,\\\\\\\"identifiedBy\\\\\\\":\\\\\\\"id\\\\\\\",\\\\\\\"decorators\\\\\\\":[]},\\\\\\\"attributes\\\\\\\":[{\\\\\\\"name\\\\\\\":\\\\\\\"id\\\\\\\",\\\\\\\"isOptional\\\\\\\":false,\\\\\\\"type\\\\\\\":\\\\\\\"String\\\\\\\"}]}\\\"}\"],[\"com.tibco.cp.IOU\",\"{\\\"$schema\\\":\\\"http://json-schema.org/draft-04/schema#\\\",\\\"title\\\":\\\"IOU\\\",\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"lender\\\":{\\\"type\\\":\\\"string\\\"},\\\"borrower\\\":{\\\"type\\\":\\\"string\\\"},\\\"amt\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"currency\\\":{\\\"type\\\":\\\"string\\\"},\\\"quantity\\\":{\\\"type\\\":\\\"integer\\\",\\\"default\\\":\\\"0\\\"}},\\\"required\\\":[\\\"currency\\\",\\\"quantity\\\"]},\\\"paid\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"currency\\\":{\\\"type\\\":\\\"string\\\"},\\\"quantity\\\":{\\\"type\\\":\\\"integer\\\",\\\"default\\\":\\\"0\\\"}},\\\"required\\\":[\\\"currency\\\",\\\"quantity\\\"]},\\\"linearId\\\":{\\\"type\\\":\\\"string\\\"}},\\\"required\\\":[\\\"lender\\\",\\\"borrower\\\",\\\"amt\\\",\\\"paid\\\",\\\"linearId\\\"],\\\"description\\\":\\\"{\\\\\\\"metadata\\\\\\\":{\\\\\\\"type\\\\\\\":\\\\\\\"Asset\\\\\\\",\\\\\\\"parent\\\\\\\":\\\\\\\"com.tibco.dovetail.system.LinearState\\\\\\\",\\\\\\\"isAbstract\\\\\\\":false,\\\\\\\"identifiedBy\\\\\\\":\\\\\\\"linearId\\\\\\\",\\\\\\\"decorators\\\\\\\":[]},\\\\\\\"attributes\\\\\\\":[{\\\\\\\"name\\\\\\\":\\\\\\\"lender\\\\\\\",\\\\\\\"type\\\\\\\":\\\\\\\"com.tibco.dovetail.system.Party\\\\\\\",\\\\\\\"isRef\\\\\\\":true},{\\\\\\\"name\\\\\\\":\\\\\\\"borrower\\\\\\\",\\\\\\\"type\\\\\\\":\\\\\\\"com.tibco.dovetail.system.Party\\\\\\\",\\\\\\\"isRef\\\\\\\":true},{\\\\\\\"name\\\\\\\":\\\\\\\"amt\\\\\\\",\\\\\\\"isOptional\\\\\\\":false,\\\\\\\"type\\\\\\\":\\\\\\\"com.tibco.dovetail.system.Amount\\\\\\\"},{\\\\\\\"name\\\\\\\":\\\\\\\"paid\\\\\\\",\\\\\\\"isOptional\\\\\\\":false,\\\\\\\"type\\\\\\\":\\\\\\\"com.tibco.dovetail.system.Amount\\\\\\\"},{\\\\\\\"name\\\\\\\":\\\\\\\"linearId\\\\\\\",\\\\\\\"isOptional\\\\\\\":false,\\\\\\\"type\\\\\\\":\\\\\\\"String\\\\\\\"}]}\\\"}\"],[\"com.tibco.cp.IssueIOU\",\"{\\\"$schema\\\":\\\"http://json-schema.org/draft-04/schema#\\\",\\\"title\\\":\\\"IssueIOU\\\",\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"iou\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"lender\\\":{\\\"type\\\":\\\"string\\\"},\\\"borrower\\\":{\\\"type\\\":\\\"string\\\"},\\\"amt\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"currency\\\":{\\\"type\\\":\\\"string\\\"},\\\"quantity\\\":{\\\"type\\\":\\\"integer\\\",\\\"default\\\":\\\"0\\\"}},\\\"required\\\":[\\\"currency\\\",\\\"quantity\\\"]},\\\"paid\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"currency\\\":{\\\"type\\\":\\\"string\\\"},\\\"quantity\\\":{\\\"type\\\":\\\"integer\\\",\\\"default\\\":\\\"0\\\"}},\\\"required\\\":[\\\"currency\\\",\\\"quantity\\\"]},\\\"linearId\\\":{\\\"type\\\":\\\"string\\\"}},\\\"required\\\":[\\\"lender\\\",\\\"borrower\\\",\\\"amt\\\",\\\"paid\\\",\\\"linearId\\\"]},\\\"transactionId\\\":{\\\"type\\\":\\\"string\\\"},\\\"timestamp\\\":{\\\"format\\\":\\\"date-time\\\",\\\"type\\\":\\\"string\\\"}},\\\"required\\\":[\\\"iou\\\",\\\"transactionId\\\",\\\"timestamp\\\"],\\\"description\\\":\\\"{\\\\\\\"metadata\\\\\\\":{\\\\\\\"type\\\\\\\":\\\\\\\"Transaction\\\\\\\",\\\\\\\"parent\\\\\\\":\\\\\\\"org.hyperledger.composer.system.Transaction\\\\\\\",\\\\\\\"isAbstract\\\\\\\":false,\\\\\\\"identifiedBy\\\\\\\":\\\\\\\"transactionId\\\\\\\",\\\\\\\"decorators\\\\\\\":[{\\\\\\\"name\\\\\\\":\\\\\\\"InitiatedBy\\\\\\\",\\\\\\\"args\\\\\\\":[\\\\\\\"$tx.iou.lender\\\\\\\"]}]},\\\\\\\"attributes\\\\\\\":[{\\\\\\\"name\\\\\\\":\\\\\\\"iou\\\\\\\",\\\\\\\"isOptional\\\\\\\":false,\\\\\\\"type\\\\\\\":\\\\\\\"com.tibco.cp.IOU\\\\\\\"},{\\\\\\\"name\\\\\\\":\\\\\\\"transactionId\\\\\\\",\\\\\\\"isOptional\\\\\\\":false,\\\\\\\"type\\\\\\\":\\\\\\\"String\\\\\\\"},{\\\\\\\"name\\\\\\\":\\\\\\\"timestamp\\\\\\\",\\\\\\\"isOptional\\\\\\\":false,\\\\\\\"type\\\\\\\":\\\\\\\"DateTime\\\\\\\"}]}\\\"}\"],[\"com.tibco.cp.TransferIOU\",\"{\\\"$schema\\\":\\\"http://json-schema.org/draft-04/schema#\\\",\\\"title\\\":\\\"TransferIOU\\\",\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"iou\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"lender\\\":{\\\"type\\\":\\\"string\\\"},\\\"borrower\\\":{\\\"type\\\":\\\"string\\\"},\\\"amt\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"currency\\\":{\\\"type\\\":\\\"string\\\"},\\\"quantity\\\":{\\\"type\\\":\\\"integer\\\",\\\"default\\\":\\\"0\\\"}},\\\"required\\\":[\\\"currency\\\",\\\"quantity\\\"]},\\\"paid\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"currency\\\":{\\\"type\\\":\\\"string\\\"},\\\"quantity\\\":{\\\"type\\\":\\\"integer\\\",\\\"default\\\":\\\"0\\\"}},\\\"required\\\":[\\\"currency\\\",\\\"quantity\\\"]},\\\"linearId\\\":{\\\"type\\\":\\\"string\\\"}},\\\"required\\\":[\\\"lender\\\",\\\"borrower\\\",\\\"amt\\\",\\\"paid\\\",\\\"linearId\\\"]},\\\"newLender\\\":{\\\"type\\\":\\\"string\\\"},\\\"transactionId\\\":{\\\"type\\\":\\\"string\\\"},\\\"timestamp\\\":{\\\"format\\\":\\\"date-time\\\",\\\"type\\\":\\\"string\\\"}},\\\"required\\\":[\\\"iou\\\",\\\"newLender\\\",\\\"transactionId\\\",\\\"timestamp\\\"],\\\"description\\\":\\\"{\\\\\\\"metadata\\\\\\\":{\\\\\\\"type\\\\\\\":\\\\\\\"Transaction\\\\\\\",\\\\\\\"parent\\\\\\\":\\\\\\\"org.hyperledger.composer.system.Transaction\\\\\\\",\\\\\\\"isAbstract\\\\\\\":false,\\\\\\\"identifiedBy\\\\\\\":\\\\\\\"transactionId\\\\\\\",\\\\\\\"decorators\\\\\\\":[{\\\\\\\"name\\\\\\\":\\\\\\\"InitiatedBy\\\\\\\",\\\\\\\"args\\\\\\\":[\\\\\\\"$tx.iou.lender\\\\\\\"]}]},\\\\\\\"attributes\\\\\\\":[{\\\\\\\"name\\\\\\\":\\\\\\\"iou\\\\\\\",\\\\\\\"type\\\\\\\":\\\\\\\"com.tibco.cp.IOU\\\\\\\",\\\\\\\"isRef\\\\\\\":true},{\\\\\\\"name\\\\\\\":\\\\\\\"newLender\\\\\\\",\\\\\\\"type\\\\\\\":\\\\\\\"com.tibco.dovetail.system.Party\\\\\\\",\\\\\\\"isRef\\\\\\\":true},{\\\\\\\"name\\\\\\\":\\\\\\\"transactionId\\\\\\\",\\\\\\\"isOptional\\\\\\\":false,\\\\\\\"type\\\\\\\":\\\\\\\"String\\\\\\\"},{\\\\\\\"name\\\\\\\":\\\\\\\"timestamp\\\\\\\",\\\\\\\"isOptional\\\\\\\":false,\\\\\\\"type\\\\\\\":\\\\\\\"DateTime\\\\\\\"}]}\\\"}\"],[\"com.tibco.cp.SettleIOU\",\"{\\\"$schema\\\":\\\"http://json-schema.org/draft-04/schema#\\\",\\\"title\\\":\\\"SettleIOU\\\",\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"iou\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"lender\\\":{\\\"type\\\":\\\"string\\\"},\\\"borrower\\\":{\\\"type\\\":\\\"string\\\"},\\\"amt\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"currency\\\":{\\\"type\\\":\\\"string\\\"},\\\"quantity\\\":{\\\"type\\\":\\\"integer\\\",\\\"default\\\":\\\"0\\\"}},\\\"required\\\":[\\\"currency\\\",\\\"quantity\\\"]},\\\"paid\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"currency\\\":{\\\"type\\\":\\\"string\\\"},\\\"quantity\\\":{\\\"type\\\":\\\"integer\\\",\\\"default\\\":\\\"0\\\"}},\\\"required\\\":[\\\"currency\\\",\\\"quantity\\\"]},\\\"linearId\\\":{\\\"type\\\":\\\"string\\\"}},\\\"required\\\":[\\\"lender\\\",\\\"borrower\\\",\\\"amt\\\",\\\"paid\\\",\\\"linearId\\\"]},\\\"payments\\\":{\\\"type\\\":\\\"array\\\",\\\"items\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"assetId\\\":{\\\"type\\\":\\\"string\\\"},\\\"owner\\\":{\\\"type\\\":\\\"string\\\"},\\\"amt\\\":{\\\"type\\\":\\\"object\\\",\\\"properties\\\":{\\\"currency\\\":{\\\"type\\\":\\\"string\\\"},\\\"quantity\\\":{\\\"type\\\":\\\"integer\\\",\\\"default\\\":\\\"0\\\"}},\\\"required\\\":[\\\"currency\\\",\\\"quantity\\\"]}},\\\"required\\\":[\\\"assetId\\\",\\\"owner\\\",\\\"amt\\\"]}},\\\"transactionId\\\":{\\\"type\\\":\\\"string\\\"},\\\"timestamp\\\":{\\\"format\\\":\\\"date-time\\\",\\\"type\\\":\\\"string\\\"}},\\\"required\\\":[\\\"iou\\\",\\\"payments\\\",\\\"transactionId\\\",\\\"timestamp\\\"],\\\"description\\\":\\\"{\\\\\\\"metadata\\\\\\\":{\\\\\\\"type\\\\\\\":\\\\\\\"Transaction\\\\\\\",\\\\\\\"parent\\\\\\\":\\\\\\\"org.hyperledger.composer.system.Transaction\\\\\\\",\\\\\\\"isAbstract\\\\\\\":false,\\\\\\\"identifiedBy\\\\\\\":\\\\\\\"transactionId\\\\\\\",\\\\\\\"decorators\\\\\\\":[{\\\\\\\"name\\\\\\\":\\\\\\\"InitiatedBy\\\\\\\",\\\\\\\"args\\\\\\\":[\\\\\\\"$tx.iou.borrower\\\\\\\"]}]},\\\\\\\"attributes\\\\\\\":[{\\\\\\\"name\\\\\\\":\\\\\\\"iou\\\\\\\",\\\\\\\"type\\\\\\\":\\\\\\\"com.tibco.cp.IOU\\\\\\\",\\\\\\\"isRef\\\\\\\":true},{\\\\\\\"name\\\\\\\":\\\\\\\"payments\\\\\\\",\\\\\\\"type\\\\\\\":\\\\\\\"com.tibco.dovetail.system.Cash\\\\\\\",\\\\\\\"isRef\\\\\\\":true,\\\\\\\"isArray\\\\\\\":true},{\\\\\\\"name\\\\\\\":\\\\\\\"transactionId\\\\\\\",\\\\\\\"isOptional\\\\\\\":false,\\\\\\\"type\\\\\\\":\\\\\\\"String\\\\\\\"},{\\\\\\\"name\\\\\\\":\\\\\\\"timestamp\\\\\\\",\\\\\\\"isOptional\\\\\\\":false,\\\\\\\"type\\\\\\\":\\\\\\\"DateTime\\\\\\\"}]}\\\"}\"]]"
		  },
		  "id": "SmartContractTXNTrigger",
		  "handlers": [
		   {
			"description": "",
			"settings": {
			 "transaction": "com.tibco.cp.IssueIOU"
			},
			"outputs": {
			 "transactionInput": {
			  "metadata": "{\"$schema\":\"http://json-schema.org/draft-04/schema#\",\"title\":\"IssueIOU\",\"type\":\"object\",\"properties\":{\"iou\":{\"type\":\"object\",\"properties\":{\"lender\":{\"type\":\"string\"},\"borrower\":{\"type\":\"string\"},\"amt\":{\"type\":\"object\",\"properties\":{\"currency\":{\"type\":\"string\"},\"quantity\":{\"type\":\"integer\",\"default\":\"0\"}},\"required\":[\"currency\",\"quantity\"]},\"paid\":{\"type\":\"object\",\"properties\":{\"currency\":{\"type\":\"string\"},\"quantity\":{\"type\":\"integer\",\"default\":\"0\"}},\"required\":[\"currency\",\"quantity\"]},\"linearId\":{\"type\":\"string\"}},\"required\":[\"lender\",\"borrower\",\"amt\",\"paid\",\"linearId\"]},\"transactionId\":{\"type\":\"string\"},\"timestamp\":{\"format\":\"date-time\",\"type\":\"string\"}},\"required\":[\"iou\",\"transactionId\",\"timestamp\"],\"description\":\"{\\\"metadata\\\":{\\\"type\\\":\\\"Transaction\\\",\\\"parent\\\":\\\"org.hyperledger.composer.system.Transaction\\\",\\\"isAbstract\\\":false,\\\"identifiedBy\\\":\\\"transactionId\\\",\\\"decorators\\\":[{\\\"name\\\":\\\"InitiatedBy\\\",\\\"args\\\":[\\\"$tx.iou.lender\\\"]}]},\\\"attributes\\\":[{\\\"name\\\":\\\"iou\\\",\\\"isOptional\\\":false,\\\"type\\\":\\\"com.tibco.cp.IOU\\\"},{\\\"name\\\":\\\"transactionId\\\",\\\"isOptional\\\":false,\\\"type\\\":\\\"String\\\"},{\\\"name\\\":\\\"timestamp\\\",\\\"isOptional\\\":false,\\\"type\\\":\\\"DateTime\\\"}]}\"}",
			  "value": ""
			 }
			},
			"action": {
			 "ref": "github.com/TIBCOSoftware/flogo-contrib/action/flow",
			 "data": {
			  "flowURI": "res://flow:IssueIOU"
			 },
			 "mappings": {
			  "input": [],
			  "output": []
			 }
			}
		   }
		  ]
		 }
		],
		"ui": "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",
		"contrib": "W3sicmVmIjoiZ2l0aHViLmNvbS9USUJDT1NvZnR3YXJlL2RvdmV0YWlsLWNvbnRyaWIvU21hcnRDb250cmFjdCIsInMzbG9jYXRpb24iOiJ7VVNFUklEfS9TbWFydENvbnRyYWN0In1d"
	   }`

	appCfg := &app.Config{}

	jsonParser := json.NewDecoder(strings.NewReader(model))
	err := jsonParser.Decode(&appCfg)
	assert.Nil(t, err)

	e, err := NewEngine(appCfg)
	assert.Nil(t, err)

	newm, err := json.Marshal(e.app)
	assert.Nil(t, err)

	fmt.Print(string(newm))
}
