/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package com.tibco.dovetail.container.cordapp;

import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.container.corda.CordaUtil;
import com.tibco.dovetail.corda.json.StateAndRefSerializer;
import com.tibco.dovetail.core.runtime.services.IDataService;

import net.corda.core.contracts.ContractState;
import net.corda.core.contracts.StateAndRef;
import net.corda.core.node.ServiceHub;
import net.corda.core.node.services.Vault.Page;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;


public class AppDataService implements IDataService {
	private ServiceHub serviceHub;
	private Map<String, StateAndRef> states = new LinkedHashMap<String, StateAndRef>();
	
	public AppDataService(ServiceHub hub) {
		this.serviceHub = hub;
	}

	@Override
	public DocumentContext putState(String assetName, String assetKey, DocumentContext assetValue) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public DocumentContext getState(String assetName, String assetKey, DocumentContext keyValue) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public DocumentContext deleteState(String assetName, String assetKey, DocumentContext keyValue) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public List<DocumentContext> lookupState(String assetName, String assetKey, DocumentContext keyValue) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public List<DocumentContext> getHistory(String assetName, String assetKey, DocumentContext keyValue) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public List<DocumentContext> queryState(Object query){
		try {
			VaultQuery queryObj = (VaultQuery)query;
			Page<ContractState> results = serviceHub.getVaultService().queryBy((Class<ContractState>) Class.forName(queryObj.getState()), queryObj.getCriteria());
			
			List<DocumentContext> docs = new ArrayList<DocumentContext>();
			results.getStates().forEach(s -> {
				states.put(StateAndRefSerializer.getRef(s), s);
				docs.add(CordaUtil.toJsonObject(s));
			});
			return docs;
		} catch(Exception e) {
			throw new RuntimeException("AppDataService::queryState", e);
		}
	}
	
	public StateAndRef getStateRef(String ref) {
		return this.states.get(ref);
	}
	
}
