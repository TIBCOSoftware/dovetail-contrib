package com.tibco.dovetail.container.cordapp;

import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.container.corda.CordaUtil;
import com.tibco.dovetail.core.runtime.transaction.ITransactionService;
import com.tibco.dovetail.core.runtime.transaction.TxnInputAttribute;
import com.tibco.dovetail.core.runtime.util.JsonUtil;

import net.corda.core.identity.AbstractParty;

public class AppTransactionService implements ITransactionService {

	private LinkedHashMap<String, Object> flowInputs;
	private String ourIdentity;
	private String transactionName;
	
	public AppTransactionService( LinkedHashMap<String, Object> flowInputs,  String flowName, AbstractParty selfIdentity) {
		this.flowInputs = flowInputs;
		this.ourIdentity = CordaUtil.partyToString(selfIdentity);
		this.transactionName = flowName;
	}
	
	@Override
	public Map<String, Object> resolveTransactionInput(List<TxnInputAttribute> txnInputs) {
		LinkedHashMap<String, Object> values = new LinkedHashMap<String, Object>();
		values.put("ourIdentity", this.ourIdentity);
	
		if (txnInputs.size() == 1 && txnInputs.get(0).getName() == "transactionInput") {
			values.put("transactionInput", CordaUtil.toJsonObject(flowInputs.get("transactionInput")));
		} else {
			DocumentContext doc = JsonUtil.getJsonParser().parse("{}");
			txnInputs.forEach ( in -> {
				String attr = in.getName();
				if(in.getType().equals("any")) {
					values.put(in.getName(), flowInputs.get(in.getName())); //ledgerTxn
				} else {
					Object value = flowInputs.get(attr);
					
					if(value != null) {
						doc.put("$", attr, CordaUtil.toJsonObject(value).json());
					} else {
						doc.put("$", attr, CordaUtil.toJsonObject(in.getValue()).json());
					}
					values.put("transactionInput", doc);
				}
				
			});
		}
		
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
