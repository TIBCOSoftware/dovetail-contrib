package com.tibco.dovetail.container.corda;

import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.core.runtime.services.IDataService;

import net.corda.core.contracts.ContractState;
import net.corda.finance.contracts.asset.Cash;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;

import javax.persistence.criteria.CriteriaBuilder.In;

public class CordaDataService implements IDataService {
	 private List<DocumentContext> outputStates = new ArrayList<DocumentContext>();
	 private List<DocumentContext> inputStates = new ArrayList<DocumentContext>();
	 private List<String> inputStatesClassName = new ArrayList<String>();
	 
	public CordaDataService(List<ContractState> inputs) {
		inputs.forEach(state -> {
			inputStates.add(CordaUtil.toJsonObject(state));
			if(state instanceof Cash.State)
				inputStatesClassName.add("com.tibco.dovetail.system.Cash");
			else
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

	@Override
	public DocumentContext putState(String assetName, String assetKey, DocumentContext assetValue) {
		String key = assetValue.read("$." + assetKey).toString();
		if(assetName.equals("com.tibco.dovetail.system.Cash")) {
			LinkedHashMap<String, Object> mapV = assetValue.json();
			mapV.remove(assetKey);
		}
		
		outputStates.add(assetValue);
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
	
}
