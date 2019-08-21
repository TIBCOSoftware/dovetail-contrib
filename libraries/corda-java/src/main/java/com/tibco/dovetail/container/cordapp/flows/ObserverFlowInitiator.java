package com.tibco.dovetail.container.cordapp.flows;

import java.util.List;

import com.tibco.dovetail.container.cordapp.ProgressTrackerSteps;

import co.paralleluniverse.fibers.Suspendable;
import net.corda.core.flows.FlowException;
import net.corda.core.flows.FlowLogic;
import net.corda.core.flows.FlowSession;
import net.corda.core.flows.InitiatingFlow;
import net.corda.core.identity.Party;
import net.corda.core.transactions.SignedTransaction;
import net.corda.core.utilities.ProgressTracker;

@InitiatingFlow
public class ObserverFlowInitiator extends FlowLogic<Void>{

	 List<Party> observers;
	 List<Party> signingParties;
	 SignedTransaction txn;
	 boolean isConfidential;
	 ProgressTracker tracker ;
	 
	 public ObserverFlowInitiator(SignedTransaction tx, List<Party> obs, boolean isConfidential, List<Party> signParty, ProgressTracker atracker) {
		 observers = obs;
		 txn = tx;
		 this.isConfidential = isConfidential;
		 this.signingParties = signParty;
		 this.tracker = atracker;
	 }
	 
       
	@Override
	@Suspendable
	public Void call() throws FlowException {
		
		for(Party p : observers) {
			FlowSession session = initiateFlow(p);
			tracker.setCurrentStep(ProgressTrackerSteps.SEND_SIGNED_TRANSACTION);
			session.send(txn);
			tracker.setCurrentStep(ProgressTrackerSteps.SYNC_IDENTITIY_REQ);
			session.send(this.isConfidential);
			if(this.isConfidential) {
				tracker.setCurrentStep(ProgressTrackerSteps.SYNC_OTHER_IDENTITIES_OBSERVER);
			 	session.send(this.signingParties);
			}
		}
		
		return null;
	}
	
	@Override
	public ProgressTracker getProgressTracker() {
		return this.tracker;
	}
	
	public static ProgressTracker getTracker() {
		return new ProgressTracker(ProgressTrackerSteps.SEND_SIGNED_TRANSACTION, ProgressTrackerSteps.SYNC_IDENTITIY_REQ, ProgressTrackerSteps.SYNC_OTHER_IDENTITIES_OBSERVER);
	}
 }
