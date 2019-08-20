package cordapp.trigger.schedulable;

import java.util.Map;

import com.tibco.dovetail.core.model.flow.HandlerConfig;
import com.tibco.dovetail.core.runtime.flow.TransactionFlow;
import com.tibco.dovetail.core.runtime.transaction.TxnInputAttribute;
import com.tibco.dovetail.core.runtime.trigger.DefaultTriggerImpl;

public class schedulable extends DefaultTriggerImpl {

	@Override
	protected void processTxnInput(TransactionFlow flow, HandlerConfig cfg) throws Exception {
		//flow properties
		Map<String, Object> properties = cfg.getSettings();
	    flow.setProperties(properties);
		
    		TxnInputAttribute txnAttr = new TxnInputAttribute();
    		txnAttr.setName("transactionInput");
    		txnAttr.setType(properties.get("asset").toString());
    		txnAttr.setArray(false);
    		txnAttr.setAssetRef(true);
    		
    		flow.addTxnInput(txnAttr);
		
	}
	/*
	private Map<String, TransactionFlow> handlers = new LinkedHashMap<String, TransactionFlow>();
	private List<AppProperty> properties;
	
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
        		
            		TxnInputAttribute txnAttr = new TxnInputAttribute();
            		txnAttr.setName("transactionInput");
            		txnAttr.setType(properties.get("asset").toString());
            		txnAttr.setArray(false);
            		txnAttr.setAssetRef(true);
            		
            		flow.addTxnInput(txnAttr);
	  
	            handlers.put(handlerConfigs[j].getFlowName(), flow);
	            lookup.put(handlerConfigs[j].getFlowName(), this);
			}
			
			 return lookup;
		}catch(Exception e) {
			throw new RuntimeException(e);
		}
	}*/



}
