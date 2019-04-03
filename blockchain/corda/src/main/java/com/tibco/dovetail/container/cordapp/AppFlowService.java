package com.tibco.dovetail.container.cordapp;

import java.security.PublicKey;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.HashSet;
import java.util.List;
import java.util.Map;
import java.util.Set;
import com.tibco.dovetail.container.corda.CordaCommandDataWithData;

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
import co.paralleluniverse.fibers.Suspendable;

public class AppFlowService {
	private TransactionBuilder builder;
	private ArrayList<StateAndRef<?>> inputStates = new ArrayList<StateAndRef<?>>();
	private ArrayList<ContractState> outputStates = new ArrayList<ContractState>();
	private ArrayList<CordaCommandDataWithData> commands = new ArrayList<CordaCommandDataWithData>();
	private Set<PublicKey> signers = new HashSet<PublicKey>();
	private Party notary;
	private AppContainer container;
	private FlowLogic<?> flow;
	private FlowSession otherParty;
	private boolean isInitiator;
	private SignedTransaction signedTxn = null;
	
	public AppFlowService(AppContainer ctnr, FlowLogic<?> flow, FlowSession otherParty, boolean isInitiating) {
		builder = new TransactionBuilder();
		this.container = ctnr;
		this.flow = flow;
		this.otherParty = otherParty;
		this.isInitiator = isInitiating;
	}
	
	public TransactionBuilder getTransactionBuilder() {
		return this.builder;
	}

	public SignedTransaction signTxn() {
		if(isInitiator) {
			if(notary == null) {
				notary = container.getServiceHub().getNetworkMapCache().getNotaryIdentities().get(0);
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
			
			this.signedTxn = this.container.getServiceHub().signInitialTransaction(builder);
			return this.signedTxn;
		} else {
			return null;
		}
	}
	
	public SignedTransaction getSignedTransaction() {
		return this.signedTxn;
	}
	
	@Suspendable
	public void commit(SignedTransaction txn) throws FlowException {
		if(isInitiator) {
			
			List<FlowSession> sessions = new ArrayList<FlowSession>();
			Map<String, Party> parties = new LinkedHashMap<String, Party>();
			ServiceHub servicehub = container.getServiceHub();
			signers.forEach(signer -> {
				Party p = servicehub.getIdentityService().partyFromKey(signer);
				if(!flow.getOurIdentity().toString().equals(p.toString()))
					parties.put(p.toString(), p);
			});
			
			parties.values().forEach(p -> sessions.add(flow.initiateFlow(p)));
			
			SignedTransaction fullysigned = (SignedTransaction) flow.subFlow(new CollectSignaturesFlow(txn, sessions));
			flow.subFlow(new FinalityFlow(fullysigned, sessions));	
		} else {
			flow.subFlow(new ReceiveFinalityFlow(otherParty, txn.getId()));
		}
	}
	
	@Suspendable
	public void commit() throws FlowException {
		if(this.signedTxn == null)
			throw new FlowException("Transaction must be signed first before committing to ledger");
		
		if(isInitiator) {
			
			List<FlowSession> sessions = new ArrayList<FlowSession>();
			Map<String, Party> parties = new LinkedHashMap<String, Party>();
			ServiceHub servicehub = container.getServiceHub();
			signers.forEach(signer -> {
				Party p = servicehub.getIdentityService().partyFromKey(signer);
				if(!flow.getOurIdentity().toString().equals(p.toString()))
					parties.put(p.toString(), p);
			});
			
			parties.values().forEach(p -> sessions.add(flow.initiateFlow(p)));
			
			SignedTransaction fullysigned = (SignedTransaction) flow.subFlow(new CollectSignaturesFlow(this.signedTxn, sessions));
			flow.subFlow(new FinalityFlow(fullysigned, sessions));	
		} else {
			flow.subFlow(new ReceiveFinalityFlow(otherParty, this.signedTxn.getId()));
		}
	}
	
	@Suspendable
	public SignedTransaction receiverVerifyAndSign(SignTransactionFlow signTransactionFlow) throws FlowException {
		if(!isInitiator) {
			this.signedTxn = flow.subFlow(signTransactionFlow);
			return this.signedTxn;
		} else {
			return null;
		}
	}
	
	@Suspendable
	public void signAndCommit() throws FlowException {
		commit(signTxn());
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
		this.notary = this.flow.getServiceHub().getIdentityService().wellKnownPartyFromX500Name(CordaX500Name.parse(notaryNm));
		if (this.notary == null) {
			throw new RuntimeException("notary party " + notaryNm + " is not found");
		}
	}
	
	public FlowSession getFlowOtherParty() {
		return this.otherParty;
	}
	
	public boolean isInitiatingFlow() {
		return this.isInitiator;
	}

}
