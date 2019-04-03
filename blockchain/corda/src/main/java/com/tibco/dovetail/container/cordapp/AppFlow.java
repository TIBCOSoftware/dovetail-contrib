package com.tibco.dovetail.container.cordapp;

import java.io.InputStream;
import java.security.PublicKey;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;

import com.tibco.dovetail.container.corda.CordaCommandDataWithData;
import com.tibco.dovetail.core.model.flow.FlowAppConfig;
import com.tibco.dovetail.core.runtime.engine.DovetailEngine;
import com.tibco.dovetail.core.runtime.trigger.ITrigger;

import co.paralleluniverse.fibers.Suspendable;
import net.corda.core.contracts.ContractState;
import net.corda.core.contracts.StateAndRef;
import net.corda.core.flows.CollectSignaturesFlow;
import net.corda.core.flows.FinalityFlow;
import net.corda.core.flows.FlowException;
import net.corda.core.flows.FlowLogic;
import net.corda.core.flows.FlowSession;
import net.corda.core.flows.ReceiveFinalityFlow;
import net.corda.core.flows.SignTransactionFlow;
import net.corda.core.identity.CordaX500Name;
import net.corda.core.identity.Party;
import net.corda.core.node.ServiceHub;
import net.corda.core.transactions.SignedTransaction;
import net.corda.core.transactions.TransactionBuilder;

public abstract class AppFlow extends FlowLogic<SignedTransaction>{
	private TransactionBuilder builder = new TransactionBuilder();
	private ArrayList<StateAndRef<?>> inputStates = new ArrayList<StateAndRef<?>>();
	private ArrayList<ContractState> outputStates = new ArrayList<ContractState>();
	private ArrayList<CordaCommandDataWithData> commands = new ArrayList<CordaCommandDataWithData>();
	private Set<PublicKey> signers = new HashSet<PublicKey>();
	private Party notary;
	private boolean isInitiator;

    
	public AppFlow(boolean initiating) {
		this.isInitiator = initiating;
	}
	
	public void runFlow(String flowName, ITrigger trigger, LinkedHashMap<String, Object> args) {
       try {
             System.out.println("****** run flow " + flowName + "... ******");
             AppContainer ctnr = new AppContainer(this);
             
             AppTransactionService txnSvc = new AppTransactionService(args, flowName, getOurIdentity());
            
             trigger.invoke(ctnr, txnSvc);
             
             System.out.println("****** flow " + flowName + " done ********");
         }catch (Exception e){
             throw new RuntimeException(e);
         }
	}
	
	protected synchronized LinkedHashMap<String, ITrigger> compileAndCacheTrigger(InputStream txJson) {
   	 try {
	        //compile flow app and cache the trigger object
        	 	FlowAppConfig app = FlowAppConfig.parseModel(txJson);
        	 	DovetailEngine engine = new DovetailEngine(app);
        	 	return engine.getTriggers();
	        
        }catch(Exception e) {
        		throw new RuntimeException(e);
        }
   }
	
	
	public TransactionBuilder getTransactionBuilder() {
		return this.builder;
	}

	public SignedTransaction initiatorSignTxn() {
		
		if(notary == null) {
			notary = this.getServiceHub().getNetworkMapCache().getNotaryIdentities().get(0);
		}
		builder.setNotary(notary);
		
		inputStates.forEach(in -> {
			builder.addInputState(in);
			in.getState().getData().getParticipants().forEach(p -> signers.add(p.getOwningKey()));
		});
		
		outputStates.forEach(out -> {
			builder.addOutputState(out);
			out.getParticipants().forEach(p -> signers.add(p.getOwningKey()));
		});
		
		commands.forEach(cmd -> builder.addCommand(cmd, new ArrayList<PublicKey>(signers)));
		
		return this.getServiceHub().signInitialTransaction(builder);
	}
	
	
	@Suspendable
	public SignedTransaction initiatorCommit(SignedTransaction txn) throws FlowException {
		
		List<FlowSession> sessions  = getSignerFlowSessions();
			
		SignedTransaction fullysigned = (SignedTransaction) subFlow(new CollectSignaturesFlow(txn, sessions));
		return subFlow(new FinalityFlow(fullysigned, sessions));	
		
	}
	
	private List<FlowSession> getSignerFlowSessions() {
		List<FlowSession> sessions = new ArrayList<FlowSession>();
		Map<String, Party> parties = new LinkedHashMap<String, Party>();
		ServiceHub servicehub = this.getServiceHub();
		signers.forEach(signer -> {
			Party p = servicehub.getIdentityService().partyFromKey(signer);
			if(!getOurIdentity().toString().equals(p.toString()))
				parties.put(p.toString(), p);
		});
		
		parties.values().forEach(p -> sessions.add(initiateFlow(p)));
		return sessions;
	}
	
	
	
	@Suspendable
	public SignedTransaction receiverSignAndCommit(SignTransactionFlow signTransactionFlow, FlowSession otherParty) throws FlowException {
		
		SignedTransaction txn = subFlow(signTransactionFlow);
		return subFlow(new ReceiveFinalityFlow(otherParty, txn.getId()));
	}
	
	@Suspendable
	public SignedTransaction initiatorSignAndCommit() throws FlowException {
		return initiatorCommit(initiatorSignTxn());
	}
	
	public void addInputStates(List<StateAndRef<?>> inputs) {
		inputs.addAll(inputs);
	}
	
	public void addInputState(StateAndRef<?> input) {
		inputStates.add(input);
	}

	public void addOutputStates(List<ContractState> outputs) {
		outputStates.addAll(outputs);
	}
	
	public List<ContractState> getOutputStates() {
		return this.outputStates;
	}
	
	public void addOutputState(ContractState output) {
		outputStates.add(output);
	}

	public void addCommand(CordaCommandDataWithData cmd) {
		commands.add(cmd);
	}
	
	@Suspendable
	public void setTransactionNotory(String notaryNm) {
		this.notary = this.getServiceHub().getIdentityService().wellKnownPartyFromX500Name(CordaX500Name.parse(notaryNm));
		if (this.notary == null) {
			throw new RuntimeException("notary party " + notaryNm + " is not found");
		}
	}
	
	
	public boolean isInitiatingFlow() {
		return this.isInitiator;
	}

	
}
