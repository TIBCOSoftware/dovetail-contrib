package dovetail_ledger_corda.trigger.scheduler;

import java.util.Map;

import com.tibco.dovetail.core.model.flow.HandlerConfig;
import com.tibco.dovetail.core.runtime.flow.TransactionFlow;
import com.tibco.dovetail.core.runtime.transaction.TxnInputAttribute;
import com.tibco.dovetail.core.runtime.trigger.DefaultTriggerImpl;


public class scheduler extends DefaultTriggerImpl {

	@Override
	protected void processTxnInput(TransactionFlow flow, HandlerConfig cfg) throws Exception {
		 //flow properties
		Map<String, Object> properties = cfg.getSettings();
	    flow.setProperties(properties);
		
    		TxnInputAttribute txnAttr = new TxnInputAttribute();
    		txnAttr.setName("transactionInput");
    		txnAttr.setType(properties.get("assetName").toString());
    		txnAttr.setArray(false);
    		txnAttr.setAssetRef(true);
    		
    		
    		flow.addTxnInput(txnAttr);
		
	}
	
}
