package com.tibco.dovetail.container.cordapp;

import com.tibco.dovetail.container.cordapp.flows.IdentitySyncFlowInitiator;
import com.tibco.dovetail.container.cordapp.flows.ObserverFlowInitiator;

import net.corda.confidential.IdentitySyncFlow;
import net.corda.core.flows.CollectSignaturesFlow;
import net.corda.core.flows.FinalityFlow;
import net.corda.core.utilities.ProgressTracker;
import net.corda.core.utilities.ProgressTracker.Step;

public class ProgressTrackerSteps  {
	public static Step CHECK_INITIATOR = new Step("Checking current lender is initiating flow.");
	public static Step PREPARATION = new Step("Obtaining Obligation from vault.");
	public static Step BUILD_TRANSACTION = new Step("Building transaction.");
	public static Step SIGN_TRANSACTION = new Step("Signing transaction.");
	public static Step RECEIVE_TRANSACTION = new Step("Receiving signed transaction.");
	
	public static Step SYNC_OUR_IDENTITY = new Step("Syncing our identity with the counterparties.") {
            @Override
			public ProgressTracker childProgressTracker() {
            		return IdentitySyncFlow.Send.Companion.tracker();
            }
	};
	
	public static Step COLLECT_SIGS = new Step("Collecting counterparty signatures.") {
		@Override
		public ProgressTracker childProgressTracker() {
        		return CollectSignaturesFlow.Companion.tracker();
        }
	};
	
	public static Step SYNC_OTHER_IDENTITIES = new Step("Making other parties sync identities with each other.");
	public static Step RECEIVE_SYNC_OTHER_IDENTITIES = new Step("Receiving other parties to sync identities.") {
		@Override
		public ProgressTracker childProgressTracker() {
			return IdentitySyncFlowInitiator.getTracker();
        }
	};
	
	public static Step SYNC_IDENTITIES = new Step("Syncing identities.") {
		@Override
		public ProgressTracker childProgressTracker() {
			return IdentitySyncFlow.Send.Companion.tracker();
        }
	};
	
	public static Step FINALISE = new Step("Finalising transaction.") {
		@Override
		public ProgressTracker childProgressTracker() {
        		return 	FinalityFlow.Companion.tracker();
        }
	};
	
	public static Step SWAP_IDENTITY = new Step("Swapping annonymous identity with the counterparty.");
	public static Step SWAP_IDENTITIY_FALSE = new Step("Letting other parties know swap annonymous identity will not be performed");
	public static Step SYNC_IDENTITIY_REQ = new Step("Letting counter parties know if sync identity will be performed.");
	public static Step RECEIVE_SYNC_IDENTITIY_REQ = new Step("Receiving sync identity request");
	public static Step RECEIVE_SWAP_IDENTITIY_REQ = new Step("Receiving swap identity request");
	public static Step RECEIVE_TXN_INITIATOR_SYNC_IDENTITIY = new Step("Receiving transaction initiator identity.");
	public static Step OBSERVER_IDENTITY_SYC_REQ = new Step("Letting observer partites know if sync identity will be performed.");
	
	public static Step SYNC_OUR_IDENTITY_WITH_OBSERVER = new Step("Syncing our identity with the observers.") {
        @Override
		public ProgressTracker childProgressTracker() {
        		return IdentitySyncFlow.Send.Companion.tracker();
        }
	};
	
	public static Step SYNC_OTHER_IDENTITIES_OBSERVER = new Step("Making observers sync identities with transaction's participants.");
	public static Step START_MANUAL_OBSERVER_FLOW = new Step("Start manual flow to send transactions to observers.") {
		@Override
		public ProgressTracker childProgressTracker() {
        		return ObserverFlowInitiator.getTracker();
        }
	};
	
	public static Step SEND_SIGNED_TRANSACTION = new Step("Send signed transactions.");
	public static Step OBSERVER_PHASE = new Step("Send signed transactions to observers.");
	public static Step RECORD_TRANSACTION = new Step("Recording transaction.");
	public static Step RUN_APP_FLOW = new Step("Running application flow");
}

