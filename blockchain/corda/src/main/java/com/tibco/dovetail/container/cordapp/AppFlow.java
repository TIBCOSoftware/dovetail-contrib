package com.tibco.dovetail.container.cordapp;

import java.io.InputStream;
import java.security.PublicKey;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.stream.Collectors;

import com.tibco.dovetail.container.corda.CordaUtil;
import com.tibco.dovetail.core.model.flow.FlowAppConfig;
import com.tibco.dovetail.core.runtime.engine.DovetailEngine;
import com.tibco.dovetail.core.runtime.trigger.ITrigger;

import co.paralleluniverse.fibers.Suspendable;
import net.corda.confidential.IdentitySyncFlow;
import net.corda.confidential.SwapIdentitiesFlow;
import net.corda.core.contracts.CommandData;
import net.corda.core.contracts.ContractState;
import net.corda.core.contracts.StateAndRef;
import net.corda.core.flows.CollectSignaturesFlow;
import net.corda.core.flows.FinalityFlow;
import net.corda.core.flows.FlowException;
import net.corda.core.flows.FlowLogic;
import net.corda.core.flows.FlowSession;
import net.corda.core.flows.ReceiveFinalityFlow;
import net.corda.core.flows.SignTransactionFlow;
import net.corda.core.identity.AbstractParty;
import net.corda.core.identity.AnonymousParty;
import net.corda.core.identity.CordaX500Name;
import net.corda.core.identity.Party;
import net.corda.core.transactions.SignedTransaction;
import net.corda.core.transactions.TransactionBuilder;

public abstract class AppFlow extends FlowLogic<SignedTransaction>{
	private TransactionBuilder builder = new TransactionBuilder();
	private ArrayList<StateAndRef<?>> inputStates = new ArrayList<StateAndRef<?>>();
	private ArrayList<ContractState> outputStates = new ArrayList<ContractState>();
	private Set<CommandData> commands = new HashSet<CommandData>();
	
	private Set<Party> counterParties = new HashSet<Party>();
	private boolean requireIdentitySync = false;
	
	private Party notary;
	private boolean isInitiator;
	private Map<PublicKey, FlowSession> opensessions = new LinkedHashMap<PublicKey, FlowSession>();
	private List<Party> swappedIdentityParties = new ArrayList<Party>();
	private AbstractParty selfIdentity;
	private Set<PublicKey> ourSignKeys = new HashSet<PublicKey>();
    private boolean isConfidential = false;
   
    
	public AppFlow(boolean initiating, boolean useAnnon) {
		this.isInitiator = initiating;
		this.isConfidential = useAnnon;
	}
	
	public void setOurIdentity() {
		if(this.isConfidential)
			selfIdentity = getOurIdentity().anonymise();
		else
			selfIdentity = getOurIdentity();
	}
	
	@Suspendable
	public void swapIdentitiesInitiator(Map<String, Object> flowIn) throws FlowException {
		for(String k : flowIn.keySet()) {
			Object v = flowIn.get(k);
			if(v instanceof Party) {
				Party p = (Party)v;
				if(getServiceHub().getNetworkMapCache().isNotary(p) == false) {
					FlowSession session = getPartyFlowSession(p);
					session.send(true);
					LinkedHashMap<Party, AnonymousParty> anonymousIdentitiesResult = subFlow(new SwapIdentitiesFlow(session));
					selfIdentity = anonymousIdentitiesResult.get(getOurIdentity());
					AnonymousParty anony = anonymousIdentitiesResult.get(p);
					flowIn.put(k, anony);
					swappedIdentityParties.add(p);	
				}
			}
		}
	}
	
	@Suspendable
	public void swapIdentitiesReceiver(FlowSession counterpartySession) throws FlowException {
		boolean exchangeIdentities = counterpartySession.receive(Boolean.class).getFromUntrustedWorld();
        if (exchangeIdentities) {
            subFlow(new SwapIdentitiesFlow(counterpartySession));
        }


        boolean syncIdentities = counterpartySession.receive(Boolean.class).getFromUntrustedWorld();
        if (syncIdentities) {
            subFlow(new IdentitySyncFlow.Receive(counterpartySession));
        }
	}
	
	public void runFlow(String flowName, ITrigger trigger, LinkedHashMap<String, Object> args) {
       try {
             System.out.println("****** run flow " + flowName + "... ******");
             AppContainer ctnr = new AppContainer(this);
             AppTransactionService txnSvc = new AppTransactionService(args, flowName, selfIdentity==null?getOurIdentity():selfIdentity);
            
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
		Set<PublicKey> signkeys= new HashSet<PublicKey>();
		Set<PublicKey> inannonkeys= new HashSet<PublicKey>();
		Set<PublicKey> outannonkeys= new HashSet<PublicKey>();
		
		if(notary == null) {
			notary = this.getServiceHub().getNetworkMapCache().getNotaryIdentities().get(0);
		}
		builder.setNotary(notary);
		
		inputStates.forEach(in -> {
			builder.addInputState(in);
			in.getState().getData().getParticipants().forEach(p -> {
				signkeys.add(p.getOwningKey());
				Party signparty = null;
				if(p instanceof AnonymousParty) {
					signparty = this.getServiceHub().getIdentityService().requireWellKnownPartyFromAnonymous(p);
					inannonkeys.add(p.getOwningKey());
				} else {
					signparty = (Party) p;
				}
				
				if(signparty.equals(getOurIdentity()))
					ourSignKeys.add(p.getOwningKey());
				else
					counterParties.add(signparty);
					
			});
		});
		
		outputStates.forEach(out -> {
			builder.addOutputState(out);
			out.getParticipants().forEach(p -> {
				signkeys.add(p.getOwningKey());
				Party signparty = null;
				if(p instanceof AnonymousParty) {
					signparty = this.getServiceHub().getIdentityService().requireWellKnownPartyFromAnonymous(p);
					outannonkeys.add(p.getOwningKey());
				} else {
					signparty = (Party) p;
				}
				
				if(signparty.equals(getOurIdentity()))
					ourSignKeys.add(p.getOwningKey());
				else
					counterParties.add(signparty);
			});
		});
		
		commands.forEach(cmd -> builder.addCommand(cmd, new ArrayList<PublicKey>(signkeys)));
		
		if(!inputStates.isEmpty() && !inannonkeys.containsAll(outannonkeys)) {
			this.requireIdentitySync = true;
		} else {
			this.requireIdentitySync = false;
		}

		return this.getServiceHub().signInitialTransaction(builder, new ArrayList(this.ourSignKeys));
	}
	
	
	@Suspendable
	public SignedTransaction initiatorCommit(SignedTransaction txn) throws FlowException {
		
		Set<FlowSession> sessions  = getCounterPartyFlowSessions();
		
		if(this.isConfidential) {
			//send false to parties that do not require swap identity
			for(FlowSession s : sessions) {
				if(!this.swappedIdentityParties.contains(s.getCounterparty()))
					s.send(false);
			}
			
			//identity sync
			if(this.requireIdentitySync) {
				for(FlowSession s: sessions) {
					s.send(true);
				}
				subFlow(new IdentitySyncFlow.Send(sessions, txn.getTx()));
			} else {
				for(FlowSession s : sessions) {
					s.send(false);
				}
			}
		}

		SignedTransaction fullysigned = (SignedTransaction) subFlow(new CollectSignaturesFlow(txn, sessions, new ArrayList<PublicKey>(this.ourSignKeys)));
		
		if(isConfidential) {
			for(FlowSession s : sessions) {
				List<Party> otherParties = sessions.stream().map(in -> in.getCounterparty()).filter(f -> !f.equals(s.getCounterparty())).collect(Collectors.toList());
				s.send(otherParties);
			}	
		}
		
		return subFlow(new FinalityFlow(fullysigned, sessions));	
		
	}
	
	private Set<FlowSession> getCounterPartyFlowSessions() {
		Set<FlowSession> sessions = new HashSet<FlowSession>();

		counterParties.forEach(signer -> sessions.add(getPartyFlowSession(signer)));
		
		return sessions;
	}
	
	public FlowSession getPartyFlowSession(Party party) {
		FlowSession session = opensessions.get(party.getOwningKey());
		if(session == null) {
			session = initiateFlow(party);
			opensessions.put(party.getOwningKey(), session);
		}
		return session;
	}
	
	@Suspendable
	public SignedTransaction receiverSignAndCommit(SignTransactionFlow signTransactionFlow, FlowSession otherParty) throws FlowException {
		SignedTransaction txn = subFlow(signTransactionFlow);
		
		if(isConfidential) {
			List syncParties = otherParty.receive(List.class).getFromUntrustedWorld() ;
			if(!syncParties.isEmpty()) {
				subFlow(new IdentitySyncFlowInitiator(syncParties,txn.getTx()));
			}
		}
		
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

	public void addCommand(CommandData cmd) {
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
