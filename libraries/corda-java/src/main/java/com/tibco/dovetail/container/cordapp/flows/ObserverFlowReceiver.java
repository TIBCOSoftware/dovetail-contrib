package com.tibco.dovetail.container.cordapp.flows;

import java.util.Arrays;
import java.util.List;

import com.tibco.dovetail.container.cordapp.ProgressTrackerSteps;

import co.paralleluniverse.fibers.Suspendable;
import net.corda.confidential.IdentitySyncFlow;
import net.corda.core.flows.FlowException;
import net.corda.core.flows.FlowLogic;
import net.corda.core.flows.FlowSession;
import net.corda.core.flows.InitiatedBy;
import net.corda.core.node.StatesToRecord;
import net.corda.core.transactions.SignedTransaction;
import net.corda.core.utilities.ProgressTracker;

@InitiatedBy(ObserverFlowInitiator.class)
public class ObserverFlowReceiver extends FlowLogic<Void>{
	
    FlowSession otherParty;
    ProgressTracker tracker;
    
	public ObserverFlowReceiver(FlowSession other) {
		otherParty = other;
		tracker = getProgressTracker();
	}
		
		
	@Override
	@Suspendable
	public Void call() throws FlowException {
		tracker.setCurrentStep(ProgressTrackerSteps.RECEIVE_TRANSACTION);
		SignedTransaction txn = otherParty.receive(SignedTransaction.class).getFromUntrustedWorld();
		
		tracker.setCurrentStep(ProgressTrackerSteps.RECORD_TRANSACTION);
        getServiceHub().recordTransactions(StatesToRecord.ALL_VISIBLE, Arrays.asList(txn));
			
        tracker.setCurrentStep(ProgressTrackerSteps.RECEIVE_SYNC_IDENTITIY_REQ);
		boolean syncIdentities = otherParty.receive(Boolean.class).getFromUntrustedWorld();
        if (syncIdentities) {
        		tracker.setCurrentStep(ProgressTrackerSteps.RECEIVE_SYNC_OTHER_IDENTITIES);
            List syncParties = otherParty.receive(List.class).getFromUntrustedWorld() ;
			if(!syncParties.isEmpty()) {
				subFlow(new IdentitySyncFlowInitiator(syncParties,txn.getTx(), ProgressTrackerSteps.RECEIVE_SYNC_OTHER_IDENTITIES.childProgressTracker()));
			}
        }
		
        return null;
	}
	
	@Override
	public ProgressTracker getProgressTracker() {
		return new ProgressTracker(ProgressTrackerSteps.RECEIVE_TRANSACTION,
						ProgressTrackerSteps.RECORD_TRANSACTION,
						ProgressTrackerSteps.RECEIVE_SYNC_IDENTITIY_REQ,
						ProgressTrackerSteps.RECEIVE_SYNC_OTHER_IDENTITIES
						);
	}
			
}
