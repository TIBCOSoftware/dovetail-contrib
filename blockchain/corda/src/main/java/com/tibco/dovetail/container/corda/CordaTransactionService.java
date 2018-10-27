package com.tibco.dovetail.container.corda;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;

import com.tibco.dovetail.core.runtime.transaction.ITransactionService;
import com.tibco.dovetail.core.runtime.transaction.TxnInputAttribute;

import net.corda.core.contracts.ContractState;
import net.corda.core.contracts.ContractsDSL;
import net.corda.core.identity.AbstractParty;
import net.corda.core.transactions.LedgerTransaction;

public class CordaTransactionService implements ITransactionService {

	private LedgerTransaction tx;
	private CordaCommandDataWithData cmd;
	
	public CordaTransactionService(LedgerTransaction tx, CordaCommandDataWithData cmd ) {
		this.tx = tx;
		this.cmd = cmd;
	}
	
	@Override
	public Map<String, Object> resolveTransactionInput(List<TxnInputAttribute> txnInputs) {
        List<ContractState> inputStates = new ArrayList<ContractState>();
        Map<String, Object> flowInputs = new LinkedHashMap<String, Object>();
        
        for(TxnInputAttribute k : txnInputs){
        		String attr = k.getName();
            Object value = cmd.getData(attr);
            if(value == null)
                throw new IllegalArgumentException("flow input " + k.getName() + " is not found in command " + cmd.getClass().getName());

        		if(value instanceof AbstractParty) {
        			flowInputs.put(attr, CordaUtil.toString(value));
        		}
        		else if(value instanceof ContractState){
        			flowInputs.put(attr, CordaUtil.toJsonObject((ContractState)value));
        			if(k.isAssetRef()) {
        				inputStates.add((ContractState)value);
        			}
    			} else if(value instanceof List){
                List<Object> objs = (List<Object>)value;
                if(objs.size() > 0) {
	                if(objs.get(0) instanceof ContractState) {
	                	    List<ContractState> states = (List<ContractState>)value;
	                	    flowInputs.put(attr, CordaUtil.toJsonObject(states));
		                if(k.isAssetRef()) {
		                    inputStates.addAll(states);
		                }
	                } else if(objs.get(0) instanceof AbstractParty) {
	                		List<String> parties = objs.stream().map(p -> CordaUtil.toString(p)).collect(Collectors.toList());
	                		flowInputs.put(attr, parties);
	                } else if (objs.get(0) instanceof String || objs.get(0)  instanceof Long || objs.get(0)  instanceof Integer || objs.get(0)  instanceof Boolean || objs.get(0)  instanceof Double) {
	                		flowInputs.put(attr, objs);
	                } else {
	                		flowInputs.put(attr, CordaUtil.toJsonObject(objs));
	                }
                }
            } else if (value instanceof String || value instanceof Long || value instanceof Integer || value instanceof Boolean || value instanceof Double) {
            		flowInputs.put(attr, value);
            } else {
            		flowInputs.put(attr, CordaUtil.toJsonObject(value));
            }
        }

        ContractsDSL.requireThat(check -> {
            List<ContractState> txIn = tx.getInputStates();
            String exp = "";
            String act = "";
            boolean success = inputStates.size() == txIn.size() && inputStates.containsAll(txIn);
            if(!success) {
            		exp = txIn.stream().map(s -> s.toString()).collect(Collectors.joining(","));
            		act = inputStates.stream().map(s -> s.toString()).collect(Collectors.joining(","));
            }
           
            check.using("inputs in command data must match transaction inputs, actural="+act + ", exp=" + exp, inputStates.size() == txIn.size() && inputStates.containsAll(txIn));
            return null;
        });

        return flowInputs;

	}

	@Override
	public boolean isTransactionSecuritySupported() {
		return false;
	}

	@Override
	public String getTransactionName() {
		return (String)this.cmd.getData("command");
	}

	@Override
	public String getTransactionInitiator() {
		return null;
	}

	@Override
	public String getInitiatorCertAttribute(String attr) {
		return null;
	}
}
