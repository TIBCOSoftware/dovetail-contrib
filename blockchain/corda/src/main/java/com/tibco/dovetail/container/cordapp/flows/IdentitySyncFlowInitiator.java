package com.tibco.dovetail.container.cordapp.flows;

import java.util.HashSet;
import java.util.List;
import java.util.Set;

import com.tibco.dovetail.container.cordapp.ProgressTrackerSteps;

import co.paralleluniverse.fibers.Suspendable;
import net.corda.confidential.IdentitySyncFlow;
import net.corda.core.flows.FlowException;
import net.corda.core.flows.FlowLogic;
import net.corda.core.flows.FlowSession;
import net.corda.core.flows.InitiatingFlow;
import net.corda.core.identity.Party;
import net.corda.core.transactions.WireTransaction;
import net.corda.core.utilities.ProgressTracker;

@InitiatingFlow
public class IdentitySyncFlowInitiator extends FlowLogic<Boolean>{
	
	 List<Party> otherParties;
	 WireTransaction txn;
	 ProgressTracker tracker;
	 public IdentitySyncFlowInitiator(List<Party> other, WireTransaction tx, ProgressTracker atracker) {
		 otherParties = other;
		 txn = tx;
		 this.tracker = atracker;
	 }
	 
       
	@Override
	@Suspendable
	public java.lang.Boolean call() throws FlowException {
		Set<FlowSession> sessions = new HashSet<FlowSession>();
		otherParties.forEach( p -> sessions.add(initiateFlow(p)));
		
		tracker.setCurrentStep(ProgressTrackerSteps.SYNC_IDENTITIES);
		subFlow(new IdentitySyncFlow.Send(sessions, txn, ProgressTrackerSteps.SYNC_IDENTITIES.childProgressTracker()));
		return true;
	}
	
	public static ProgressTracker getTracker() {
		return new ProgressTracker(ProgressTrackerSteps.SYNC_IDENTITIES);
	}
 }
