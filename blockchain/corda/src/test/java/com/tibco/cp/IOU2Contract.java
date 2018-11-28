/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package com.tibco.cp;

import com.tibco.dovetail.container.corda.CordaCommandDataWithData;
import com.tibco.dovetail.container.corda.CordaFlowContract;
import net.corda.core.contracts.Contract;
import net.corda.core.transactions.LedgerTransaction;

import java.io.InputStream;

public class IOU2Contract extends  CordaFlowContract implements Contract {
    public static final String IOU_CONTRACT_ID = "com.tibco.cp.IOUContract";

    public static class IssueIOU extends CordaCommandDataWithData {
        public IssueIOU(IOU3 iou , String transactionId , String timestamp ){

            putData("iou", iou);
            putData("transactionId", transactionId);
            putData("timestamp", timestamp);
        }
    }

    public static class TransferIOU extends CordaCommandDataWithData {
        public TransferIOU(IOU3 iou , net.corda.core.identity.Party newLender , String transactionId , String timestamp ){

            putData("iou", iou);
            putData("newLender", newLender);
            putData("transactionId", transactionId);
            putData("timestamp", timestamp);
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
}
