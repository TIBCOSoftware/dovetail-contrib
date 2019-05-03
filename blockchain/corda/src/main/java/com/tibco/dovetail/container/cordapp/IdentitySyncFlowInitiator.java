package com.tibco.dovetail.container.cordapp;

import java.util.HashSet;
import java.util.List;
import java.util.Set;

import co.paralleluniverse.fibers.Suspendable;
import net.corda.confidential.IdentitySyncFlow;
import net.corda.core.flows.FlowException;
import net.corda.core.flows.FlowLogic;
import net.corda.core.flows.FlowSession;
import net.corda.core.flows.InitiatingFlow;
import net.corda.core.identity.Party;
import net.corda.core.transactions.WireTransaction;

@InitiatingFlow
public class IdentitySyncFlowInitiator extends FlowLogic<Boolean>{
	 	 
	 List<Party> otherParties;
	 WireTransaction txn;
	 
	 public IdentitySyncFlowInitiator(List<Party> other, WireTransaction tx) {
		 otherParties = other;
		 txn = tx;
	 }
	 
       
	@Override
	@Suspendable
	public java.lang.Boolean call() throws FlowException {
		Set<FlowSession> sessions = new HashSet<FlowSession>();
		otherParties.forEach( p -> sessions.add(initiateFlow(p)));
		
		subFlow(new IdentitySyncFlow.Send(sessions, txn));
		return true;
	}
 }
