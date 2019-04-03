/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package smartcontract.trigger.transaction;

import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.core.model.flow.HandlerConfig;
import com.tibco.dovetail.core.model.flow.Resources;
import com.tibco.dovetail.core.model.flow.TriggerConfig;
import com.tibco.dovetail.core.runtime.compilers.FlowCompiler;
import com.tibco.dovetail.core.runtime.flow.ReplyData;
import com.tibco.dovetail.core.runtime.flow.TransactionFlow;
import com.tibco.dovetail.core.runtime.services.IContainerService;
import com.tibco.dovetail.core.runtime.transaction.TxnInputAttribute;
import com.tibco.dovetail.core.runtime.transaction.ITransactionService;
import com.tibco.dovetail.core.runtime.transaction.TxnACL;
import com.tibco.dovetail.core.runtime.trigger.ITrigger;

import co.paralleluniverse.fibers.Suspendable;
import smartcontract.trigger.transaction.model.composer.HLCDecorator;
import smartcontract.trigger.transaction.model.composer.HLCMetadata;
import smartcontract.trigger.transaction.model.composer.HLCMetadata.ResourceType;
import smartcontract.trigger.transaction.model.composer.HLCResource;

public class transaction implements ITrigger{
	private Map<String, TransactionFlow> handlers = new LinkedHashMap<String, TransactionFlow>();
	
	@Override
	public Map<String, ITrigger> Initialize(TriggerConfig triggerConfig)  {
		try {
			 String schema = triggerConfig.getSetting("schemas");
	         Map<String, HLCResource> metadatas = MetadataParser.parse(schema);
	            
			HandlerConfig[] handlerConfigs = triggerConfig.getHandlers();
			if(handlerConfigs == null || handlerConfigs.length == 0)
				throw new RuntimeException("No handlers defined for trigger " + triggerConfig.getName());
			
			Map<String, ITrigger> lookup = new LinkedHashMap<String, ITrigger>();
			
			for(int j=0; j<handlerConfigs.length; j++) {
				String txnName = handlerConfigs[j].getSetting("transaction").toString();
				Resources r = handlerConfigs[j].getFlow();
	
	             TransactionFlow flow = FlowCompiler.compile(r);
	
	            //flow inputs/outputs
	            HLCResource txnResource =  metadatas.get(txnName);
	            HLCMetadata metadata = txnResource.getMetadata();
	            txnResource.getAttributes().forEach(a -> {
	            		TxnInputAttribute txnAttr = new TxnInputAttribute();
	            		txnAttr.setName(a.getName());
	            		txnAttr.setType(a.getType());
	            		txnAttr.setArray(a.isArray());

	            		HLCResource attrMetadata = metadatas.get(a.getType());
	            		if(attrMetadata != null) {
		            		ResourceType rtype = attrMetadata.getMetadata().getType();
		            		if(rtype == ResourceType.Asset) {
		            			txnAttr.setAssetName(a.getType());
		            			txnAttr.setAssetRef(a.isRef());
		            		} 
		            		
		            		txnAttr.setParticipant(rtype == ResourceType.Participant);
	            		} else {
	            			txnAttr.setAssetRef(false);
	            			txnAttr.setParticipant(false);
	            		}
	            		
	            		flow.addFlowInput(txnAttr);
	            });
	            
	            if(r.getData().getMetadata() != null) {
	                flow.setFlowOutputs(r.getData().getMetadata().getOutput());
	            }
	
	            //set ACL InitiatedBy decorator
	            //first argument is comma delimited list of authorized parties in the format of $tx.path.to attribute, use * for any
	            //second argument is comma delimited list of conditions that initiator cert must meet, in the format of attribute=value
	            //TODO: testing is not done
	            HLCDecorator acldec = metadata.getDecorator("InitiatedBy");
	            if(acldec != null) {
	            		List<String> parties = new ArrayList<String>();
	            		String[] args = acldec.getArgs();
			        	if (args == null || args.length == 0) {
			        		parties.add("*");
			        	} else {
			        		for(String p : args[0].split(",")) {
			        			parties.add(p.trim().substring(4));
			        		}
		
			        		if (args.length > 1) {
			        			//attributes
			        			Map<String, String>conditions = new LinkedHashMap<String, String>();
			        			for(String c : args[1].split(",")) {
			        				String[] values = c.split("=");
			        				if(values.length != 1)
			        					throw new RuntimeException("Decorator InitiatedBy sencond argument(condtion) does not follow key=value format");
			        				
			        				conditions.put(values[0].trim(), values[1].trim());
			        				
			        			}
			        			flow.setAcl(new TxnACL(parties, conditions));
			        		}
			        	}
	            }
	            //store handler by txnName, txn without namespace, and flowId
	            String txnNoNS = txnName.substring(txnName.lastIndexOf('.'));
	            handlers.put(txnName, flow);
	            handlers.put(txnNoNS, flow);
	            handlers.put(handlerConfigs[j].getFlowName(), flow);
	            
	            lookup.put(txnName, this);
	            lookup.put(handlerConfigs[j].getFlowName(), this);
			}
			
			 return lookup;
		}catch(Exception e) {
			throw new RuntimeException(e);
		}
	}

	@Override
	public ReplyData invoke(IContainerService stub, ITransactionService txn) throws RuntimeException{
		TransactionFlow handler = handlers.get(txn.getTransactionName());
		if(handler == null)
			throw new RuntimeException("Transaction flow " + txn.getTransactionName() + " is not found");
		
		Map<String, Object> triggerData = txn.resolveTransactionInput(handler.getFlowInputs());
		triggerData.put("containerServiceStub",stub);
		
		if(txn.isTransactionSecuritySupported() && handler.getAcl() != null) {
			TxnACL acl = handler.getAcl();
			
			boolean authorized = false;
			//TO: test security
			if (acl.getAuthorizedParties().size() > 0) {
				for (String participant : acl.getAuthorizedParties()) {
					if (participant.equals("*")){
						authorized = true;
						break;
					}

					//$tx.path.to.party
					String id = findValueInMap(triggerData, participant);
					
					if (id.equals(txn.getTransactionInitiator())) {
						authorized = true;
						break;
					} 
				}
			}
			
			if(authorized)
				authorized = isAuthorized(txn, acl);
			
			if(!authorized)
				throw new RuntimeException("Security violation, " + txn.getTransactionInitiator() + " is not authorized for transaction " + txn.getTransactionName());
		}
		
		return handler.handle(stub, triggerData);
	}
	
	private boolean isAuthorized(ITransactionService txn, TxnACL acl) {
		for(String attr : acl.getConditions().keySet()) {
			String value = txn.getInitiatorCertAttribute(attr);
			
			if (value == null || !value.equals(acl.getConditions().get(attr) )) {
				return false;
			}
		}

		return true;
	}
	
	private String findValueInMap(Map<String, Object> values, String k) {
		
		Object v = values.get(k);
		if (v == null)
			return null;
		
	    if(v instanceof DocumentContext) {
	    		DocumentContext doc = (DocumentContext) v;
	        String path = "$." + k;
	        List<Object> list= doc.read(path);
	        if(list != null && list.size() > 0)
	            return (String)list.get(0);
	        else
	            return null;
	    } else {
	    		return (String)v;
	    }
	}
	
	@Override
	public TransactionFlow getHandler(String name) {
		return handlers.get(name);
	}
}

