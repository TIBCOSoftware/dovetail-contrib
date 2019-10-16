package dovetail_cordapp.trigger.flowreceiver;

import com.tibco.dovetail.core.model.flow.HandlerConfig;
import com.tibco.dovetail.core.runtime.flow.TransactionFlow;
import com.tibco.dovetail.core.runtime.transaction.TxnInputAttribute;
import com.tibco.dovetail.core.runtime.trigger.DefaultTriggerImpl;


public class flowreceiver extends DefaultTriggerImpl {
	
	@Override
	protected void processTxnInput(TransactionFlow flow, HandlerConfig cfg) throws Exception {
		
    		TxnInputAttribute txnAttr = new TxnInputAttribute();
    		txnAttr.setName("ledgerTxn");
    		txnAttr.setType("any");

    		flow.addTxnInput(txnAttr);
        
		return;
		
	}

}
