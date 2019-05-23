package cordapp.trigger.flowinitiator;

import java.util.LinkedHashMap;
import java.util.Map;

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

import smartcontract.trigger.transaction.MetadataParser;
import smartcontract.trigger.transaction.model.composer.HLCResource;

public class flowinitiator implements ITrigger{
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
        			
            		String schema = handlerConfigs[j].getOutputs().getTransactionInput().getMetadata();
            		HLCResource txnResource = MetadataParser.parseSingleSchema(schema);
	            txnResource.getAttributes().forEach(a -> {
	            		TxnInputAttribute txnAttr = new TxnInputAttribute();
	            		txnAttr.setName(a.getName());
	            		txnAttr.setType(a.getType());
	            		txnAttr.setArray(a.isArray());

	            		if(a.getType().equals("net.corda.core.identity.Party")) {
	            				txnAttr.setParticipant(true);
	            		} else {
	            			txnAttr.setParticipant(false);
	            		}
	            		
	            		flow.addFlowInput(txnAttr);
	            });
	            
	            
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
