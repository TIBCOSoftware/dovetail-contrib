package com.example.cp;

import com.tibco.dovetail.container.corda.CordaCommandDataWithData;
import com.tibco.dovetail.container.corda.CordaFlowContract;
import net.corda.core.contracts.Contract;
import net.corda.core.transactions.LedgerTransaction;
import java.io.InputStream;
import net.corda.core.serialization.CordaSerializable;
import java.util.Currency;
import java.util.List;

import net.corda.core.identity.AbstractParty;
import com.tibco.dovetail.core.runtime.trigger.ITrigger;

@CordaSerializable
public class IOUContract extends CordaFlowContract implements Contract {
    
    public static String IOUContract_CONTRACT_ID = "com.example.cp.IOUContract";
    

    @CordaSerializable
    public static class IssueIOU extends CordaCommandDataWithData {
    	 	public IssueIOU(IOU iou ){

            putData("iou", iou);
            putData("command", "com.example.cp.IssueIOU");
        }
    }

    @CordaSerializable
    public static class TransferIOU extends CordaCommandDataWithData {
    		public TransferIOU(IOU iou , net.corda.core.identity.Party newOwner  ){

            putData("iou", iou);
            putData("newOwner", newOwner);
            putData("command", "com.example.cp.TransferIOU");
        }
    }

    @CordaSerializable
    public static class SettleIOU extends CordaCommandDataWithData {
    	 	public SettleIOU(com.example.cp.IOU iou , List<net.corda.finance.contracts.asset.Cash.State> funds , net.corda.core.contracts.Amount<Currency> payAmt  , AbstractParty sendChangeTo, AbstractParty sendPaymentTo) {

            putData("iou", iou);
            putData("funds", funds);
            putData("payAmt", payAmt);
            putData("sendChangeTo", sendChangeTo);
            putData("sendPaymentTo", sendPaymentTo);
            putData("command", "com.example.cp.SettleIOU");
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
		 return IOUContractImpl.getTrigger(name);
	}

}