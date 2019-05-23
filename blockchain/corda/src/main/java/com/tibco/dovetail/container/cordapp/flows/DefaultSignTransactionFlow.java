package com.tibco.dovetail.container.cordapp.flows;

import co.paralleluniverse.fibers.Suspendable;
import net.corda.core.flows.FlowException;
import net.corda.core.flows.FlowSession;
import net.corda.core.flows.SignTransactionFlow;
import net.corda.core.transactions.SignedTransaction;

public class DefaultSignTransactionFlow extends SignTransactionFlow {

	public DefaultSignTransactionFlow(FlowSession otherSideSession) {
		super(otherSideSession);
	}

	@Override
	@Suspendable
	protected void checkTransaction(SignedTransaction arg0) throws FlowException {
	
		return;
	}

}
