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

import co.paralleluniverse.fibers.Suspendable;
import kotlin.jvm.functions.Function0;
import net.corda.core.contracts.Amount;
import net.corda.core.contracts.ContractState;
import net.corda.core.contracts.StateAndRef;
import net.corda.core.contracts.StateRef;
import net.corda.core.contracts.TransactionResolutionException;
import net.corda.core.identity.AbstractParty;
import net.corda.core.identity.Party;
import net.corda.core.node.ServiceHub;
import net.corda.core.node.services.Vault.Page;
import net.corda.core.transactions.TransactionBuilder;
import net.corda.core.utilities.OpaqueBytes;
import net.corda.finance.contracts.asset.Cash;
import net.corda.finance.workflows.asset.selection.AbstractCashSelection;

import java.sql.DatabaseMetaData;
import java.sql.SQLException;
import java.util.ArrayList;
import java.util.Currency;
import java.util.HashSet;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.UUID;
import java.util.concurrent.atomic.AtomicReference;


public class AppDataService implements IDataService<StateRef,ContractState> {
	private ServiceHub serviceHub;
	private TransactionBuilder builder;
	private UUID flowRunId;
	private Map<String, StateAndRef> states = new LinkedHashMap<String, StateAndRef>();
	
	public AppDataService(ServiceHub hub, TransactionBuilder builder, UUID runId) {
		this.serviceHub = hub;
		this.builder = builder;
		this.flowRunId = runId;
	}

	@Override
	public ContractState putState(String assetName, String assetKey, ContractState assetValue) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public ContractState getState(String assetName, String assetKey, StateRef keyValue) {
		String name = getCordaAssetName(assetName);
		try {
			ContractState s =  serviceHub.toStateAndRef(keyValue).getState().getData();
			if(s != null && name != null) {
				if(s.getClass().getName().equals(name))
					return s;
				else
					return null;
			}
			
			return s;
		}catch(Exception e) {
			throw new RuntimeException("AppDataService::getState", e);
		}
	}
	
	@Override
	public List<ContractState> getStates(String assetName, String assetKey, Set<StateRef> keyValue) {
		List<ContractState> states = new ArrayList<ContractState>();
		String name = getCordaAssetName(assetName);
		try {
			serviceHub.loadStates(keyValue).forEach(s -> {
				if(s != null && name != null) {
					if(s.getClass().getName().equals(name))
						states.add(s.getState().getData());
				}
			});
			
			return states;
		}catch(Exception e) {
			throw new RuntimeException("AppDataService::getState", e);
		}
	}

	private String getCordaAssetName(String asset) {
		if(asset.equals("com.tibco.dovetail.system.Cash"))
			return "import net.corda.finance.contracts.asset.Cash$State";
		else 
			return asset;
	}
	@Override
	public ContractState deleteState(String assetName, String assetKey, StateRef keyValue) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public List<ContractState> lookupState(String assetName, String assetKey, StateRef keyValue) {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public List<ContractState> getHistory(String assetName, String assetKey, StateRef keyValue) {
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
	
	@Suspendable
	public List<StateAndRef<Cash.State>> getFunds(Set<AbstractParty> issuers, Amount<Currency> amt) {
	
		try {
			AbstractCashSelection db = AbstractCashSelection.Companion.getInstance( () -> {
					try {
						return this.serviceHub.jdbcSession().getMetaData();
					}catch(Exception e) {
						throw new RuntimeException("getFunds error", e);
					}
				});
			
			List<StateAndRef<Cash.State>> funds = db.unconsumedCashStatesForSpending(this.serviceHub, amt, issuers, null, this.flowRunId, new HashSet<OpaqueBytes>());
			funds.forEach(s -> {
				states.put(StateAndRefSerializer.getRef(s), s);
			});
			return funds;
		}catch(Exception e) {
			throw new RuntimeException("getFunds error", e);
		}
	}
	
	public Amount<Currency> getAccountBalance(Currency c) {
		return net.corda.finance.workflows.GetBalances.getCashBalance(this.serviceHub, c);
	}

	@Override
	public boolean processPayment(DocumentContext assetValue) {
		throw new RuntimeException("not implememted");
	}
	
}
