package cordapp.trigger.schedulable;

import java.util.LinkedHashMap;
import java.util.Map;

import com.tibco.dovetail.core.model.composer.HLCResource;
import com.tibco.dovetail.core.model.composer.MetadataParser;
import com.tibco.dovetail.core.model.flow.HandlerConfig;
import com.tibco.dovetail.core.model.flow.Resources;
import com.tibco.dovetail.core.model.flow.TriggerConfig;
import com.tibco.dovetail.core.runtime.compilers.FlowCompiler;
import com.tibco.dovetail.core.runtime.flow.ReplyData;
import com.tibco.dovetail.core.runtime.flow.TransactionFlow;
import com.tibco.dovetail.core.runtime.services.IContainerService;
import com.tibco.dovetail.core.runtime.transaction.ITransactionService;
import com.tibco.dovetail.core.runtime.transaction.TxnInputAttribute;
import com.tibco.dovetail.core.runtime.trigger.ITrigger;

public class schedulable implements ITrigger {
	private Map<String, TransactionFlow> handlers = new LinkedHashMap<String, TransactionFlow>();
	
	@Override
	public Map<String, ITrigger> Initialize(TriggerConfig triggerConfig) {
		try {
			HandlerConfig[] handlerConfigs = triggerConfig.getHandlers();
			if(handlerConfigs == null || handlerConfigs.length == 0)
				throw new RuntimeException("No handlers defined for trigger " + triggerConfig.getName());
			
			Map<String, ITrigger> lookup = new LinkedHashMap<String, ITrigger>();
			
			for(int j=0; j<handlerConfigs.length; j++) {
				Resources r = handlerConfigs[j].getFlow();
				TransactionFlow flow = FlowCompiler.compile(r);
	
	            //flow properties
				Map<String, Object> properties = handlerConfigs[j].getSettings();
			    flow.setProperties(properties);
        		
            		TxnInputAttribute txnAttr = new TxnInputAttribute();
            		txnAttr.setName("transactionInput");
            		txnAttr.setType(properties.get("asset").toString());
            		txnAttr.setArray(false);
            		txnAttr.setAssetRef(true);
            		
            		flow.addFlowInput(txnAttr);
	  
	            handlers.put(handlerConfigs[j].getFlowName(), flow);
	            lookup.put(handlerConfigs[j].getFlowName(), this);
			}
			
			 return lookup;
		}catch(Exception e) {
			throw new RuntimeException(e);
		}
	}

	@Override
	public ReplyData invoke(IContainerService stub, ITransactionService txn) {
		TransactionFlow handler = handlers.get(txn.getTransactionName());
		if(handler == null)
			throw new RuntimeException("Transaction flow " + txn.getTransactionName() + " is not found");
		
		Map<String, Object> triggerData = txn.resolveTransactionInput(handler.getFlowInputs());
		
		return handler.handle(stub, triggerData);
	}

	@Override
	public TransactionFlow getHandler(String name) {
		return handlers.get(name);
	}

}
