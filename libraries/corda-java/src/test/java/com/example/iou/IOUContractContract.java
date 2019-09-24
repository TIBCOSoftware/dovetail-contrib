/*
	* Copyright © 2018. TIBCO Software Inc.
	* This file is subject to the license terms contained
	* in the license file that is distributed with this file.
	 */
	
package com.example.iou;
import com.tibco.dovetail.container.corda.CordaCommandDataWithData;
import com.tibco.dovetail.container.corda.CordaFlowContract;
import com.tibco.dovetail.core.runtime.trigger.ITrigger;

import net.corda.core.contracts.Contract;
import net.corda.core.identity.AbstractParty;
import net.corda.core.serialization.CordaSerializable;
import net.corda.core.transactions.LedgerTransaction;

import java.io.InputStream;
import java.util.Currency;
import java.util.List;

public class IOUContractContract extends  CordaFlowContract implements Contract{
	
	    public static final String IOUContractContract_CONTRACT_ID = "com.example.iou.IOUContractContract";

	    public static class IssueIOU extends CordaCommandDataWithData {
	        public IssueIOU(IOU iou , String transactionId , String timestamp ){

	            putData("iou", iou);
	            putData("transactionId", transactionId);
	            putData("timestamp", timestamp);
	        }
	    }

	    public static class TransferIOU extends CordaCommandDataWithData {
	        public TransferIOU(IOU iou , net.corda.core.identity.Party newLender , String transactionId , String timestamp ){

	            putData("iou", iou);
	            putData("newLender", newLender);
	            putData("transactionId", transactionId);
	            putData("timestamp", timestamp);
	        }
	    }
	    
	    @CordaSerializable
	    public static class SettleIOU extends CordaCommandDataWithData{
	       public SettleIOU(com.example.iou.IOU iou , List<net.corda.finance.contracts.asset.Cash.State> funds , net.corda.core.contracts.Amount<Currency> payAmt  , AbstractParty sendChangeTo, AbstractParty sendPaymentTo , String transactionId , java.time.Instant  timestamp) {
	            putData("iou", iou);
	            putData("funds", funds);
	            putData("payAmt", payAmt);
	            putData("sendChangeTo", sendChangeTo);
	            putData("sendPaymentTo", sendPaymentTo);
	            putData("transactionId", transactionId);
	            putData("timestamp", timestamp);
	            putData("command", "com.example.iou.SettleIOU");
	       }
	    }

	    @Override
	    protected String getResourceHash() {
	        return "48846cb669a097b229d841b3b8e55e8d56785b197f74f132f39d61c29c25e01e";
	    }

	    @Override
	    protected InputStream getTransactionJson() {

	        return this.getClass().getResourceAsStream("transactions.json");

	    }

		@Override
		public void verify(LedgerTransaction tx) throws IllegalArgumentException {
			verifyTransaction(tx);
		}

		@Override
		protected ITrigger getTrigger(String name) {
			 return IOUContractContractImpl.getTrigger(name);
		}

}
