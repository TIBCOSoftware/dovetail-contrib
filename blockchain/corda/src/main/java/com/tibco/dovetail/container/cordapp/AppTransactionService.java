package com.tibco.dovetail.container.cordapp;

import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.container.corda.CordaUtil;
import com.tibco.dovetail.core.runtime.transaction.ITransactionService;
import com.tibco.dovetail.core.runtime.transaction.TxnInputAttribute;
import com.tibco.dovetail.core.runtime.util.JsonUtil;

import net.corda.core.identity.Party;

public class AppTransactionService implements ITransactionService {

	private LinkedHashMap<String, Object> flowInputs;
	private String ourIdentity;
	private String transactionName;
	
	public AppTransactionService( LinkedHashMap<String, Object> flowInputs,  String flowName, Party ourIdentity) {
		this.flowInputs = flowInputs;
		this.ourIdentity = AppContainer.partyToString(ourIdentity);
		this.transactionName = flowName;
	}
	
	@Override
	public Map<String, Object> resolveTransactionInput(List<TxnInputAttribute> txnInputs) {
		LinkedHashMap<String, Object> values = new LinkedHashMap<String, Object>();
		values.put("ourIdentity", this.ourIdentity);
		DocumentContext doc = JsonUtil.getJsonParser().parse("{}");
		
		txnInputs.forEach ( in -> {
			String attr = in.getName();
			Object value = flowInputs.get(attr);
			
			if(value != null) {
				doc.put("$", attr, CordaUtil.toJsonObject(value).json());
				/*
				if(value instanceof AbstractParty) {
					values.put(attr, CordaUtil.serialize(value));
        			}
        			else if(value instanceof List){
        				List<Object> objs = (List<Object>)value;
        				
	                if(objs.size() > 0) {
		                if(objs.get(0) instanceof AbstractParty) {
		                		List<String> parties = objs.stream().map(p -> CordaUtil.serialize(p)).collect(Collectors.toList());
		                		values.put(attr, parties);
		                } else if (objs.get(0) instanceof String || objs.get(0)  instanceof Long || objs.get(0)  instanceof Integer || objs.get(0)  instanceof Boolean || objs.get(0)  instanceof Double) {
		                		values.put(attr, objs);
		                } else {
		                		values.put(attr, CordaUtil.toJsonObject(objs));
		                }
		                
	                }
        				
	            } else if (value instanceof String || value instanceof Long || value instanceof Integer || value instanceof Boolean || value instanceof Double) {
	            		values.put(attr, value);
	            } else {
	            		values.put(attr, CordaUtil.toJsonObject(value));
	            }
	            */
			} else {
				//values.put(attr, in.getValue());
				doc.put("$", attr, CordaUtil.toJsonObject(in.getValue()).json());
			}
		});
		values.put("transactionInput", doc);
		return values;
	}

	@Override
	public boolean isTransactionSecuritySupported() {
		// TODO Auto-generated method stub
		return false;
	}

	@Override
	public String getTransactionName() {
		return this.transactionName;
	}

	@Override
	public String getTransactionInitiator() {
		// TODO Auto-generated method stub
		return null;
	}

	@Override
	public String getInitiatorCertAttribute(String attr) {
		// TODO Auto-generated method stub
		return null;
	}
	
}
