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

import com.tibco.dovetail.container.cordapp.flows.IdentitySyncFlowInitiator;
import com.tibco.dovetail.container.cordapp.flows.ObserverFlowInitiator;
import com.tibco.dovetail.core.model.flow.FlowAppConfig;
import com.tibco.dovetail.core.runtime.engine.DovetailEngine;
import com.tibco.dovetail.core.runtime.trigger.ITrigger;

import co.paralleluniverse.fibers.Suspendable;
import net.corda.confidential.IdentitySyncFlow;
import net.corda.confidential.SwapIdentitiesFlow;
import net.corda.core.contracts.CommandData;
import net.corda.core.contracts.ContractState;
import net.corda.core.contracts.StateAndRef;
import net.corda.core.contracts.TimeWindow;
import net.corda.core.flows.CollectSignaturesFlow;
import net.corda.core.flows.FinalityFlow;
import net.corda.core.flows.FlowException;
import net.corda.core.flows.FlowLogic;
import net.corda.core.flows.FlowSession;
import net.corda.core.flows.ReceiveFinalityFlow;
import net.corda.core.flows.SignTransactionFlow;
import net.corda.core.identity.AbstractParty;
import net.corda.core.identity.AnonymousParty;
import net.corda.core.identity.Party;
import net.corda.core.node.StatesToRecord;
import net.corda.core.transactions.SignedTransaction;
import net.corda.core.transactions.TransactionBuilder;
import net.corda.core.utilities.ProgressTracker;

public abstract class AppFlow extends FlowLogic<SignedTransaction>{
	private TransactionBuilder builder = new TransactionBuilder();
	private ArrayList<StateAndRef<?>> inputStates = new ArrayList<StateAndRef<?>>();
	private ArrayList<ContractState> outputStates = new ArrayList<ContractState>();
	private Set<CommandData> commands = new HashSet<CommandData>();
	
	private Set<Party> counterParties = new HashSet<Party>();
	private boolean requireIdentitySync = false;
	
	private List<Party> notaries = new ArrayList<Party>();
	private List<Party> observers = new ArrayList<Party>();
	private boolean isInitiator;
	private Map<PublicKey, FlowSession> opensessions = new LinkedHashMap<PublicKey, FlowSession>();
	private List<Party> swappedIdentityParties = new ArrayList<Party>();
	private AbstractParty selfIdentity;
	private Set<PublicKey> ourSignKeys = new HashSet<PublicKey>();
    private boolean isConfidential = false;
    private boolean observerSendManual;
    
	public AppFlow(boolean initiating, boolean useAnnon) {
		this.isInitiator = initiating;
		this.isConfidential = useAnnon;
	}
	
	protected void setOurIdentity() {
		if(this.isConfidential)
			selfIdentity = getOurIdentity().anonymise();
		else
			selfIdentity = getOurIdentity();
	}
	
	@Suspendable
	public void swapIdentitiesInitiator(Map<String, Object> flowIn, List<String> participants, ProgressTracker tracker) throws FlowException {
		tracker.setCurrentStep(ProgressTrackerSteps.SWAP_IDENTITY);
		for(String k : participants) {
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
	protected void swapIdentitiesReceiver(FlowSession counterpartySession, ProgressTracker tracker) throws FlowException {
		
		tracker.setCurrentStep(ProgressTrackerSteps.RECEIVE_SWAP_IDENTITIY_REQ);
		boolean exchangeIdentities = counterpartySession.receive(Boolean.class).getFromUntrustedWorld();
        if (exchangeIdentities) {
        		tracker.setCurrentStep(ProgressTrackerSteps.SWAP_IDENTITY);
            subFlow(new SwapIdentitiesFlow(counterpartySession));
        }

        tracker.setCurrentStep(ProgressTrackerSteps.RECEIVE_SYNC_IDENTITIY_REQ);
        boolean syncIdentities = counterpartySession.receive(Boolean.class).getFromUntrustedWorld();
        if (syncIdentities) {
        	tracker.setCurrentStep(ProgressTrackerSteps.SYNC_IDENTITIES);
            subFlow(new IdentitySyncFlow.Receive(counterpartySession));
        }
	}
	
	protected void runFlow(String flowName, ITrigger trigger, LinkedHashMap<String, Object> args) {
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

	protected SignedTransaction initiatorSignTxn(ProgressTracker initiatorTracker) {
		initiatorTracker.setCurrentStep(ProgressTrackerSteps.BUILD_TRANSACTION);
		Set<PublicKey> signkeys= new HashSet<PublicKey>();
		Set<PublicKey> inannonkeys= new HashSet<PublicKey>();
		Set<PublicKey> outannonkeys= new HashSet<PublicKey>();
		
		if(notaries.isEmpty()) {
			notaries.add(this.getServiceHub().getNetworkMapCache().getNotaryIdentities().get(0));
		}
		builder.setNotary(notaries.get(0));
		
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

		initiatorTracker.setCurrentStep(ProgressTrackerSteps.SIGN_TRANSACTION);
		return this.getServiceHub().signInitialTransaction(builder, new ArrayList<PublicKey>(this.ourSignKeys));
	}
	
	
	@Suspendable
	protected SignedTransaction initiatorCommit(SignedTransaction txn, ProgressTracker initiatorTracker) throws FlowException {
		
		Set<FlowSession> sessions  = getCounterPartyFlowSessions();
		
		if(this.isConfidential) {
			initiatorTracker.setCurrentStep(ProgressTrackerSteps.SWAP_IDENTITIY_FALSE);
			//send false to parties that do not require swap identity
			for(FlowSession s : sessions) {
				if(!this.swappedIdentityParties.contains(s.getCounterparty()))
					s.send(false);
			}
			
			//identity sync
			initiatorTracker.setCurrentStep(ProgressTrackerSteps.SYNC_IDENTITIY_REQ);
			for(FlowSession s: sessions) {
				s.send(requireIdentitySync);
			}
			if(this.requireIdentitySync) {
				initiatorTracker.setCurrentStep(ProgressTrackerSteps.SYNC_OUR_IDENTITY);
				subFlow(new IdentitySyncFlow.Send(sessions, txn.getTx(), ProgressTrackerSteps.SYNC_OUR_IDENTITY.childProgressTracker()));
			} 
		}

		initiatorTracker.setCurrentStep(ProgressTrackerSteps.COLLECT_SIGS);
		SignedTransaction fullysigned = (SignedTransaction) subFlow(new CollectSignaturesFlow(txn, sessions, new ArrayList<PublicKey>(this.ourSignKeys), ProgressTrackerSteps.COLLECT_SIGS.childProgressTracker()));
		
		if(isConfidential) {
			//let involved participant parties know each other to sync identities
			initiatorTracker.setCurrentStep(ProgressTrackerSteps.SYNC_OTHER_IDENTITIES);
			for(FlowSession s : sessions) {
				List<Party> otherParties = this.counterParties.stream().filter(f -> !f.equals(s.getCounterparty())).collect(Collectors.toList());
				s.send(otherParties);
			}	
		}
		
		//observers, send as part of finality flow
		Set<FlowSession> obssessions  = new HashSet<FlowSession>();
		if(this.observers.size() > 0 && !this.observerSendManual) {
			this.observers.forEach(o -> obssessions.add(initiateFlow(o)));	
			sessions.addAll(obssessions);
		}
		
		initiatorTracker.setCurrentStep(ProgressTrackerSteps.FINALISE);
		SignedTransaction tx = subFlow(new FinalityFlow(fullysigned, sessions, ProgressTrackerSteps.FINALISE.childProgressTracker()));	
		
		if(this.observers.size() > 0) {
			initiatorTracker.setCurrentStep(ProgressTrackerSteps.OBSERVER_PHASE);
			if( !this.observerSendManual) {
				initiatorTracker.setCurrentStep(ProgressTrackerSteps.OBSERVER_IDENTITY_SYC_REQ);
				//sync identities
				//require sync
				for(FlowSession s : obssessions) {
					s.send(isConfidential);
				}
				if(isConfidential) {
					initiatorTracker.setCurrentStep(ProgressTrackerSteps.SYNC_OUR_IDENTITY_WITH_OBSERVER);
					subFlow(new IdentitySyncFlow.Send(obssessions, txn.getTx(), ProgressTrackerSteps.SYNC_OUR_IDENTITY_WITH_OBSERVER.childProgressTracker()));
					for(FlowSession s: obssessions) {
						initiatorTracker.setCurrentStep(ProgressTrackerSteps.SYNC_OTHER_IDENTITIES_OBSERVER);
						s.send(new ArrayList<Party>(this.counterParties));
					}
				} 
			} else{
				try {
						List<Party> allsigners = new ArrayList<Party>(this.counterParties);
						allsigners.add(getOurIdentity());
						initiatorTracker.setCurrentStep(ProgressTrackerSteps.START_MANUAL_OBSERVER_FLOW);
						subFlow(new ObserverFlowInitiator(tx, this.observers, this.isConfidential, allsigners, ProgressTrackerSteps.START_MANUAL_OBSERVER_FLOW.childProgressTracker() ));
				}catch(Exception e) {
					throw new FlowException(e);
				}
			}
		}
		
		return tx;
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
	protected SignedTransaction receiverSignAndCommit(SignTransactionFlow signTransactionFlow, FlowSession otherParty, ProgressTracker tracker) throws FlowException {
		tracker.setCurrentStep(ProgressTrackerSteps.SIGN_TRANSACTION);
		SignedTransaction txn = subFlow(signTransactionFlow);
		
		if(isConfidential) {
			tracker.setCurrentStep(ProgressTrackerSteps.RECEIVE_SYNC_OTHER_IDENTITIES);
			List<Party> syncParties = otherParty.receive(List.class).getFromUntrustedWorld() ;
			if(!syncParties.isEmpty()) {
				subFlow(new IdentitySyncFlowInitiator(syncParties,txn.getTx(), ProgressTrackerSteps.RECEIVE_SYNC_OTHER_IDENTITIES.childProgressTracker()));
			}
		}
		
		tracker.setCurrentStep(ProgressTrackerSteps.RECORD_TRANSACTION);
		return subFlow(new ReceiveFinalityFlow(otherParty, txn.getId()));
	}
	
	@Suspendable
	protected SignedTransaction observerRecordTxn(FlowSession otherParty, ProgressTracker tracker) throws FlowException {
		tracker.setCurrentStep(ProgressTrackerSteps.RECORD_TRANSACTION);
		SignedTransaction txn = subFlow(new ReceiveFinalityFlow(otherParty, null, StatesToRecord.ALL_VISIBLE));
		
		tracker.setCurrentStep(ProgressTrackerSteps.RECEIVE_SYNC_IDENTITIY_REQ);
		boolean syncIdentities = otherParty.receive(Boolean.class).getFromUntrustedWorld();
        if (syncIdentities) {
        		tracker.setCurrentStep(ProgressTrackerSteps.RECEIVE_TXN_INITIATOR_SYNC_IDENTITIY);
            subFlow(new IdentitySyncFlow.Receive(otherParty));
            
            tracker.setCurrentStep(ProgressTrackerSteps.RECEIVE_SYNC_OTHER_IDENTITIES);
            List<Party> syncParties = otherParty.receive(List.class).getFromUntrustedWorld() ;
			if(!syncParties.isEmpty()) {
				subFlow(new IdentitySyncFlowInitiator(syncParties,txn.getTx(), ProgressTrackerSteps.RECEIVE_SYNC_OTHER_IDENTITIES.childProgressTracker()));
			}
        }
		
		return txn;
	}
	
	@Suspendable
	protected SignedTransaction initiatorSignAndCommit(ProgressTracker initiatorTracker) throws FlowException {
		return initiatorCommit(initiatorSignTxn(initiatorTracker), initiatorTracker);
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
	
	public boolean isInitiatingFlow() {
		return this.isInitiator;
	}

	protected void addNotary(Party n) throws FlowException{
		if (this.getServiceHub().getNetworkMapCache().isNotary(n))
			this.notaries.add(n);
		else
			throw new FlowException("Party " + n.getName().getCommonName() + " is not a notary");
	}
	
	protected void addNotary(List<Party> n) throws FlowException{
		for(Party p : n) {
			addNotary(p);
		}
	}
	
	
	protected void addObserver(Party p) {
		this.observers.add(p);
	}
	
	protected void addObserver(List<Party> p) {
		this.observers.addAll(p);
	}
	
	protected void setObserverConfig(boolean isManual) {
		this.observerSendManual = isManual;
	}
	
	public void setTimeWindow(TimeWindow tw) {
		this.builder.setTimeWindow(tw);
	}
	
	public static ProgressTracker getInitiatorProgressTracker(boolean isConfidential, boolean hasObservers) {
		if(isConfidential && hasObservers)
			return new ProgressTracker(
					ProgressTrackerSteps.SWAP_IDENTITY,
					ProgressTrackerSteps.RUN_APP_FLOW,
					ProgressTrackerSteps.BUILD_TRANSACTION,
					ProgressTrackerSteps.SIGN_TRANSACTION,
					ProgressTrackerSteps.SWAP_IDENTITIY_FALSE,
					ProgressTrackerSteps.SYNC_IDENTITIY_REQ,
					ProgressTrackerSteps.SYNC_OUR_IDENTITY,
					ProgressTrackerSteps.COLLECT_SIGS,
					ProgressTrackerSteps.SYNC_OTHER_IDENTITIES,
					ProgressTrackerSteps.FINALISE,
					ProgressTrackerSteps.OBSERVER_PHASE,
					ProgressTrackerSteps.OBSERVER_IDENTITY_SYC_REQ,
					ProgressTrackerSteps.SYNC_OUR_IDENTITY_WITH_OBSERVER,
					ProgressTrackerSteps.SYNC_OTHER_IDENTITIES_OBSERVER,
					ProgressTrackerSteps.START_MANUAL_OBSERVER_FLOW);
		else if(isConfidential && !hasObservers)
			return new ProgressTracker(
					ProgressTrackerSteps.SWAP_IDENTITY,
					ProgressTrackerSteps.RUN_APP_FLOW,
					ProgressTrackerSteps.BUILD_TRANSACTION,
					ProgressTrackerSteps.SIGN_TRANSACTION,
					ProgressTrackerSteps.SWAP_IDENTITIY_FALSE,
					ProgressTrackerSteps.SYNC_IDENTITIY_REQ,
					ProgressTrackerSteps.SYNC_OUR_IDENTITY,
					ProgressTrackerSteps.COLLECT_SIGS,
					ProgressTrackerSteps.SYNC_OTHER_IDENTITIES,
					ProgressTrackerSteps.FINALISE);
		else if(!isConfidential && hasObservers)
			return new ProgressTracker(
					ProgressTrackerSteps.RUN_APP_FLOW,
					ProgressTrackerSteps.BUILD_TRANSACTION,
					ProgressTrackerSteps.SIGN_TRANSACTION,
					ProgressTrackerSteps.COLLECT_SIGS,
					ProgressTrackerSteps.FINALISE,
					ProgressTrackerSteps.OBSERVER_PHASE,
					ProgressTrackerSteps.OBSERVER_IDENTITY_SYC_REQ,
					ProgressTrackerSteps.SYNC_OUR_IDENTITY_WITH_OBSERVER,
					ProgressTrackerSteps.SYNC_OTHER_IDENTITIES_OBSERVER,
					ProgressTrackerSteps.START_MANUAL_OBSERVER_FLOW);
		else
			return new ProgressTracker(
					ProgressTrackerSteps.RUN_APP_FLOW,
					ProgressTrackerSteps.BUILD_TRANSACTION,
					ProgressTrackerSteps.SIGN_TRANSACTION,
					ProgressTrackerSteps.COLLECT_SIGS,
					ProgressTrackerSteps.FINALISE);
			
	}
	
	public static ProgressTracker getResponderProgressTracker(boolean isConfidential) {
		if(isConfidential)
			return new ProgressTracker(
					ProgressTrackerSteps.RECEIVE_SWAP_IDENTITIY_REQ,
					ProgressTrackerSteps.SWAP_IDENTITY,
					ProgressTrackerSteps.RECEIVE_SYNC_IDENTITIY_REQ,
					ProgressTrackerSteps.SYNC_IDENTITIES,
					ProgressTrackerSteps.RUN_APP_FLOW,
					ProgressTrackerSteps.SIGN_TRANSACTION,
					ProgressTrackerSteps.RECORD_TRANSACTION,
					ProgressTrackerSteps.RECEIVE_SYNC_OTHER_IDENTITIES);
		else
			return new ProgressTracker(
					ProgressTrackerSteps.RUN_APP_FLOW,
					ProgressTrackerSteps.SIGN_TRANSACTION,
					ProgressTrackerSteps.RECORD_TRANSACTION);
	}
	
	public static ProgressTracker getObserverProgressTracker() {
	
		return new ProgressTracker(
					ProgressTrackerSteps.RUN_APP_FLOW,
					ProgressTrackerSteps.RECORD_TRANSACTION,
					ProgressTrackerSteps.RECEIVE_SYNC_IDENTITIY_REQ,
					ProgressTrackerSteps.RECEIVE_TXN_INITIATOR_SYNC_IDENTITIY,
					ProgressTrackerSteps.RECEIVE_SYNC_OTHER_IDENTITIES);
	}
}
