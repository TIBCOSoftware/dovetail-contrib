package com.tibco.dovetail.container.cordapp;

import co.paralleluniverse.fibers.Suspendable;
import net.corda.confidential.IdentitySyncFlow;
import net.corda.core.flows.FlowException;
import net.corda.core.flows.FlowLogic;
import net.corda.core.flows.FlowSession;
import net.corda.core.flows.InitiatedBy;

@InitiatedBy(IdentitySyncFlowInitiator.class)
public class IdentitySyncFlowReceiver extends FlowLogic<Boolean>{
	
    FlowSession otherSideSession;
    
	public IdentitySyncFlowReceiver(FlowSession other) {
    		otherSideSession = other;
	}
		
		
	@Override
	@Suspendable
	public java.lang.Boolean call() throws FlowException {
		subFlow(new IdentitySyncFlow.Receive(otherSideSession));
        return true;
	}
			
}
