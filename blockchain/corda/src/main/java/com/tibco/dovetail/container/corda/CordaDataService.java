/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package com.tibco.dovetail.container.corda;

import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.core.runtime.services.IDataService;

import net.corda.core.contracts.CommandData;
import net.corda.core.contracts.ContractState;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;


public class CordaDataService implements IDataService {
	 private List<DocumentContext> outputStates = new ArrayList< DocumentContext>();
	 private List<DocumentContext> inputStates = new ArrayList<DocumentContext>();
	 private List<String> inputStatesClassName = new ArrayList<String>();
	 private ContractCommandOutput cmdOutput = new ContractCommandOutput();
	 
	public CordaDataService(List<ContractState> inputs) {
		inputs.forEach(state -> {
			inputStates.add(CordaUtil.toJsonObject(state));
			
			inputStatesClassName.add(state.getClass().getName());
		});
	}
    
    @Override
    public List<DocumentContext> queryState(Object query) {
        throw new IllegalArgumentException("query is not supported");
    }

    public List<DocumentContext> getModifiedStates() {
        return outputStates;
    }
    
    public ContractCommandOutput getContractCommandOutput() {
        return this.cmdOutput;
    }

	@Override
	public DocumentContext putState(String assetName, String assetKey, DocumentContext assetValue) {
			
		outputStates.add(assetValue);
		cmdOutput.addOutputState(assetName, assetValue, null);
		return assetValue;
	}

	@Override
	public DocumentContext getState(String assetName, String assetKey, DocumentContext keyValue) {
		for (int i=0; i<inputStatesClassName.size(); i++) {
			if(inputStatesClassName.get(i).equals(assetName)) {
				DocumentContext state = inputStates.get(i);
				String[] keys = assetKey.split(",");
				boolean found = true;
				for(String k : keys) {
					if(!state.read("$." + k).equals(keyValue.read("$." + k))) {
						found = false;
						break;
					}
				}
				
				if(found) {
					return state;
				}
			}
		}
		
		return null;
	}

	@Override
	public DocumentContext deleteState(String assetName, String assetKey, DocumentContext keyValue) {
		return getState(assetName, assetKey, keyValue);
	}

	@Override
	public List<DocumentContext> getHistory(String assetName, String assetKey, DocumentContext keyValue) {
		throw new IllegalArgumentException("history is not supported");
	}

	@Override
	public List<DocumentContext> lookupState(String assetName, String assetKey, DocumentContext keyValue) {
		List<DocumentContext> result = new ArrayList<DocumentContext>();
		for (int i=0; i<inputStatesClassName.size(); i++) {
			if(inputStatesClassName.get(i).equals(assetName)) {
				DocumentContext state = inputStates.get(i);
				String[] keys = assetKey.split(assetKey);
				boolean found = true;
				for(String k : keys) {
					if(!state.read("$." + k).equals(keyValue.read("$." + k))) {
						found = false;
						break;
					}
				}
				
				if(found) {
					result.add(state);
				}
			}
		}
	
		return result;
	}

	@Override
	public boolean processPayment(DocumentContext assetValue) {
		LinkedHashMap inputval = assetValue.json();
		
		String payTo = inputval.get("sendPaymentTo").toString();
		String changeTo = inputval.get("sendChangeTo").toString();
		LinkedHashMap payAmt =(LinkedHashMap)inputval.get("paymentAmt");
		List<LinkedHashMap> funds = (List<LinkedHashMap>) inputval.get("funds");
		
		long remaining = Long.valueOf(payAmt.get("quantity").toString());

		LinkedHashMap<String, Long> payoutputs = new LinkedHashMap<String, Long>();
		LinkedHashMap<String, Long> changeoutputs = new LinkedHashMap<String, Long>();
		net.corda.finance.contracts.asset.Cash c = new net.corda.finance.contracts.asset.Cash();
		CommandData move = c.generateMoveCommand();
		
		Map<String, List<LinkedHashMap>> groupbyIssuer = funds.stream().collect(Collectors.groupingBy(f -> f.get("issuer").toString()));
		for(String issuer : groupbyIssuer.keySet()) {
			List<LinkedHashMap> v = groupbyIssuer.get(issuer);
			long payByIssuer = payoutputs.get(issuer) == null? 0 :  payoutputs.get(issuer) ;
			long chgByIssuer = changeoutputs.get(issuer) == null? 0 :  changeoutputs.get(issuer) ;
			
			for(LinkedHashMap m : v) {
				this.cmdOutput.addCommand(move, CordaUtil.decodeKey(m.get("owner").toString()));
				long amt = Long.valueOf(((LinkedHashMap)m.get("amt")).get("quantity").toString());
				if (remaining > 0) {
					if (amt >= remaining) {
						chgByIssuer += amt - remaining;
						payByIssuer += remaining;
						remaining = 0;
					} else {
						remaining = remaining - amt;
						payByIssuer += amt;
					}
				} else {
					chgByIssuer += amt;
				}
			}
			
			payoutputs.put(issuer, payByIssuer);
			changeoutputs.put(issuer, chgByIssuer);
		}
		
		if (remaining > 0)
			throw new RuntimeException("payment::not enough funds");	
		
		payoutputs.forEach((k,v) -> {
			if(v > 0) {
				DocumentContext doc = CordaUtil.toJsonObject(groupbyIssuer.get(k).get(0));
				doc.put("$", "owner", payTo);
				doc.put("$.amt", "quantity", v);
				outputStates.add(doc);
				this.cmdOutput.addOutputState("net.corda.finance.contracts.asset.Cash$State", doc, move);
				this.cmdOutput.addCommand(move, CordaUtil.decodeKey(payTo));
			}
			
		});
		
		changeoutputs.forEach((k,v) -> {
			if(v > 0) {
				DocumentContext doc = CordaUtil.toJsonObject(groupbyIssuer.get(k).get(0));
				doc.put("$", "owner", changeTo);
				doc.put("$.amt", "quantity", v);
				outputStates.add(doc);
				this.cmdOutput.addOutputState("net.corda.finance.contracts.asset.Cash$State", doc, move);
				this.cmdOutput.addCommand(move, CordaUtil.decodeKey(changeTo));
			}
			
		});
		
		return true;
	}
	
	
	
}
