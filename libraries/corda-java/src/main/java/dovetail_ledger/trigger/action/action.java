package dovetail_ledger.trigger.action;

import java.util.LinkedHashMap;
import java.util.Map;

import com.tibco.dovetail.core.model.flow.HandlerConfig;
import com.tibco.dovetail.core.model.metadata.ResourceDescriptor;
import com.tibco.dovetail.core.runtime.flow.TransactionFlow;
import com.tibco.dovetail.core.runtime.transaction.TxnACL;
import com.tibco.dovetail.core.runtime.transaction.TxnInputAttribute;
import com.tibco.dovetail.core.runtime.trigger.DefaultTriggerImpl;

public class action extends DefaultTriggerImpl {

	@Override
	protected void processTxnInput(TransactionFlow flow, HandlerConfig cfg) throws Exception {
		ResourceDescriptor txnResource = cfg.getTransactionInputMetadata();
        txnResource.getAttributes().forEach(a -> {
	        	TxnInputAttribute txnAttr = new TxnInputAttribute();
	        	txnAttr.setName(a.getName());
	        	txnAttr.setType(a.getType());
	        	txnAttr.setArray(a.isArray());
        		txnAttr.setAssetName(a.getType());
        		txnAttr.setAssetRef(a.isAsset() && a.isRef());
        		txnAttr.setParticipant(false);
        		txnAttr.setAsset(a.isAsset());
        		txnAttr.setReferenceData(a.isReferenceData());
	            	
	        	flow.addTxnInput(txnAttr);
        });
        
        Map<String, String> auth = txnResource.getMetadata().getAuthorizedUserAndCerts();
        TxnACL acl = new TxnACL();
        auth.forEach((p, c) -> {
        		Map<String, String>conditions = new LinkedHashMap<String, String>();
        		if(c != null) {
				for(String s : c.split(",")) {
					String[] values = s.split("=");
					if(values.length != 1)
						throw new RuntimeException("actor certification condition does not follow comma delimited key=value format");
					
					conditions.put(values[0].trim(), values[1].trim());
				}
        		}
        		acl.addAthorizedParty(p, conditions);
        });
			//attributes
			
		flow.setAcl(acl);
		flow.setTimewindowl(txnResource.getMetadata().getTimewindow());
		int idx = txnResource.getMetadata().getAsset().lastIndexOf(".");
		String ns = txnResource.getMetadata().getAsset().substring(0, idx);
		flow.setTransactionName(ns + "." + flow.getTransactionName());
	}

}