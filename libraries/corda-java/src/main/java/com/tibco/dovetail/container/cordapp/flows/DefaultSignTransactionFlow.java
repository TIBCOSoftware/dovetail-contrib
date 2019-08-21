package com.tibco.dovetail.container.cordapp.flows;

import java.util.LinkedHashMap;

import com.tibco.dovetail.container.cordapp.AppFlow;
import com.tibco.dovetail.core.runtime.trigger.ITrigger;

import co.paralleluniverse.fibers.Suspendable;
import net.corda.core.flows.FlowException;
import net.corda.core.flows.FlowSession;
import net.corda.core.flows.SignTransactionFlow;
import net.corda.core.transactions.SignedTransaction;

public class DefaultSignTransactionFlow extends SignTransactionFlow {
	private AppFlow appflow;
	ITrigger trigger;
	String flowName;
	public DefaultSignTransactionFlow(FlowSession otherSideSession, AppFlow appf, String flowName, ITrigger trigger) {
		super(otherSideSession);
		this.appflow = appf;
		this.trigger = trigger;
		this.flowName = flowName;
	}

	@Override
	@Suspendable
	protected void checkTransaction(SignedTransaction arg0) throws FlowException {
		if(this.appflow != null) {
			LinkedHashMap<String, Object> args = new LinkedHashMap<String, Object>();
			args.put("ledgerTxn", arg0);
			appflow.runFlow(this.flowName, trigger, args);
		}
		return;
	}

}
