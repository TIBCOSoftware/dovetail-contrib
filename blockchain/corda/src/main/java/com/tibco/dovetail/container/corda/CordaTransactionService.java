/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package com.tibco.dovetail.container.corda;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.container.corda.CordaUtil;
import com.tibco.dovetail.core.runtime.transaction.ITransactionService;
import com.tibco.dovetail.core.runtime.transaction.TxnInputAttribute;
import com.tibco.dovetail.core.runtime.util.JsonUtil;

import net.corda.core.contracts.ContractState;
import net.corda.core.contracts.ContractsDSL;
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
		List<DocumentContext> inputStates = new ArrayList<DocumentContext>();
        Map<String, Object> flowInputs = new LinkedHashMap<String, Object>();
        DocumentContext doc = JsonUtil.getJsonParser().parse("{}");
        
        try {
	        for(TxnInputAttribute k : txnInputs){
	        		String attr = k.getName();
	            Object value = cmd.getData(attr);
	            if(value == null && !attr.equalsIgnoreCase("transactionId") && !attr.equalsIgnoreCase("timestamp"))
	                throw new IllegalArgumentException("flow input " + attr + " is not found in command " + cmd.getClass().getName());

	            if(value == null)
	            		continue;
	            
	            DocumentContext valdoc = CordaUtil.getInstance().toJsonObject(value);
	            doc.put("$", attr, valdoc.json());

	            if(k.isAssetRef()) {
	            		if (value instanceof List) {
	            			List<?> objs = (List<?>)value;
	    	                 objs.forEach(o -> {
	    	                	 	DocumentContext val = CordaUtil.getInstance().toJsonObject(o);
	    	                	 	inputStates.add(val);
	    	                 }); 
	            		} else {
	            			inputStates.add(valdoc);
	            		}
	            }
	        }
	        flowInputs.put("transactionInput", doc);
        }catch(Exception e) {
        		throw new RuntimeException(e);
        }
        
        if(tx != null) {
	        ContractsDSL.requireThat(check -> {
	            List<ContractState> txIn = tx.getInputStates();
	            List<DocumentContext> txInDocs = new ArrayList<DocumentContext>();
	            txIn.forEach(in -> txInDocs.add(CordaUtil.getInstance().toJsonObject(in)));
	            
	            CordaUtil.getInstance().compare(txInDocs, inputStates);
	            
	            return null;
	        });
        }

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
