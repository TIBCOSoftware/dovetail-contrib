package dovetail_cordapp.trigger.flowinitiator;

import java.util.Map;

import com.tibco.dovetail.core.model.flow.HandlerConfig;
import com.tibco.dovetail.core.model.metadata.ResourceDescriptor;
import com.tibco.dovetail.core.runtime.flow.TransactionFlow;
import com.tibco.dovetail.core.runtime.transaction.TxnInputAttribute;
import com.tibco.dovetail.core.runtime.trigger.DefaultTriggerImpl;

public class flowinitiator extends DefaultTriggerImpl {

	@Override
	protected void processTxnInput(TransactionFlow flow, HandlerConfig cfg) throws Exception {
		//flow properties
		Map<String, Object> properties = cfg.getSettings();
	    flow.setProperties(properties);
			
    		ResourceDescriptor txnResource = cfg.getTransactionInputMetadata();
        txnResource.getAttributes().forEach(a -> {
        		TxnInputAttribute txnAttr = new TxnInputAttribute();
        		txnAttr.setName(a.getName());
        		txnAttr.setType(a.getType());
        		txnAttr.setArray(a.isArray());

        		flow.addTxnInput(txnAttr);
        });
        
		
	}
}
