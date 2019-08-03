package cordapp.trigger.flowinitiator;

import java.util.Map;

import com.tibco.dovetail.core.model.composer.HLCResource;
import com.tibco.dovetail.core.model.composer.MetadataParser;
import com.tibco.dovetail.core.model.flow.HandlerConfig;
import com.tibco.dovetail.core.runtime.flow.TransactionFlow;
import com.tibco.dovetail.core.runtime.transaction.TxnInputAttribute;
import com.tibco.dovetail.core.runtime.trigger.DefaultTriggerImpl;

public class flowinitiator extends DefaultTriggerImpl {
	/*private Map<String, TransactionFlow> handlers = new LinkedHashMap<String, TransactionFlow>();
	List<AppProperty> properties;
	
	@Override
	public Map<String, ITrigger> Initialize(TriggerConfig triggerConfig, List<AppProperty> pp) {
		try {
			this.properties = pp;
			
			HandlerConfig[] handlerConfigs = triggerConfig.getHandlers();
			if(handlerConfigs == null || handlerConfigs.length == 0)
				throw new RuntimeException("No handlers defined for trigger " + triggerConfig.getName());
			
			Map<String, ITrigger> lookup = new LinkedHashMap<String, ITrigger>();
			
			for(int j=0; j<handlerConfigs.length; j++) {
				TransactionFlow flow = FlowCompiler.compile(handlerConfigs[j]);
	
	            //flow properties
				Map<String, Object> properties = handlerConfigs[j].getSettings();
			    flow.setProperties(properties);
        			
            		String schema = handlerConfigs[j].getTransactionInputMetadata();
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
	            		
	            		flow.addTxnInput(txnAttr);
	            });
	            
	            
	            handlers.put(handlerConfigs[j].getFlowName(), flow);
	            lookup.put(handlerConfigs[j].getFlowName(), this);
			}
			
			 return lookup;
		}catch(Exception e) {
			throw new RuntimeException(e);
		}
	
	}
*/

	@Override
	protected void processTxnInput(TransactionFlow flow, HandlerConfig cfg) throws Exception {
		//flow properties
		Map<String, Object> properties = cfg.getSettings();
	    flow.setProperties(properties);
			
    		String schema = cfg.getTransactionInputMetadata();
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
        		
        		flow.addTxnInput(txnAttr);
        });
        
		
	}
}
