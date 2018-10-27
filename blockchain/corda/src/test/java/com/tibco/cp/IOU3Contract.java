package com.tibco.cp;

import com.tibco.dovetail.container.corda.CordaCommandDataWithData;
import com.tibco.dovetail.container.corda.CordaFlowContract;

import net.corda.core.contracts.Contract;
import net.corda.core.identity.Party;
import net.corda.core.transactions.LedgerTransaction;
import net.corda.finance.contracts.asset.Cash;

import java.io.InputStream;
import java.util.List;

public class IOU3Contract extends CordaFlowContract implements Contract {
    public static final String IOU_CONTRACT_ID = "com.tibco.cp.IOUContract";

    @Override
    protected String getResourceHash() {
        return "abcd";
    }

    @Override
    protected InputStream getTransactionJson() {

        return this.getClass().getResourceAsStream("transactions.json");

    }

    @Override
    protected InputStream getSchemasJson() {
        return this.getClass().getResourceAsStream("schemas.json");
    }

    public static class IssueIOU extends CordaCommandDataWithData {
        public IssueIOU(IOU3 iou){
            putData("iou", iou);
            putData("command", "com.tibco.cp.IssueIOU");
        }
    }
    public static class TransferIOU extends CordaCommandDataWithData {
        public TransferIOU(IOU3 iou, Party newLender){
            putData("iou", iou);
            putData("newLender" , newLender);
            putData("command", "com.tibco.cp.TransferIOU");
        }
    }
    public static class SettleIOU extends CordaCommandDataWithData {
        public SettleIOU(IOU3 iou, List<Cash.State> payments){
            putData("iou", iou);
            putData("payments", payments);
            putData("command", "com.tibco.cp.SettleIOU");
        }
    }
	@Override
	public void verify(LedgerTransaction arg0) throws IllegalArgumentException {
		// TODO Auto-generated method stub
		
	}
}
