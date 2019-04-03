package cordapp.activity.txnbuilder;

import java.io.Serializable;
import java.util.LinkedHashMap;
import java.util.List;

import com.fasterxml.jackson.core.type.TypeReference;
import com.google.common.collect.Lists;
import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.container.corda.CordaCommandDataWithData;
import com.tibco.dovetail.container.corda.CordaFlowContract;
import com.tibco.dovetail.container.corda.CordaUtil;
import com.tibco.dovetail.container.cordapp.AppContainer;
import com.tibco.dovetail.container.cordapp.AppDataService;
import com.tibco.dovetail.container.cordapp.AppFlow;
import com.tibco.dovetail.core.runtime.activity.IActivity;
import com.tibco.dovetail.core.runtime.engine.Context;

import kotlin.Pair;

import net.corda.core.contracts.ContractState;
import net.corda.core.contracts.StateAndRef;

public class txnbuilder implements IActivity{
	@Override
	public void eval(Context context) throws IllegalArgumentException {
		
		String contractClass = context.getInput("contractClass").toString();
		String command = context.getInput("command").toString();
		String schema = context.getInput("inputSchema").toString();
		
		DocumentContext input = (DocumentContext)context.getInput("input");
		
		try {
			
			AppFlow txservice = ((AppContainer) context.getContainerService()).getFlowService();
			AppDataService dataservice = (AppDataService) context.getContainerService().getDataService();
			
			LinkedHashMap<String, Object> inputs = (LinkedHashMap<String, Object>)input.json();
			List<BuilderSchemaAttribute> attrs = (List<BuilderSchemaAttribute>) CordaUtil.deserialize(schema, new TypeReference<List<BuilderSchemaAttribute>>() {});
			
			CordaCommandDataWithData cordacmd = new CordaCommandDataWithData();
			
			for( BuilderSchemaAttribute attr : attrs) {
				Object value = inputs.get(attr.getName());
				if(value != null) {
					if(attr.isAsset() && attr.isRef()) {
						StateAndRef state  = dataservice.getStateRef(value.toString());
						txservice.addInputState(state);
						value = state.getState().getData();
					}
				} 
				
				cordacmd.putData(attr.getName(), value);
			};
			
			cordacmd.putData("command", command);
			cordacmd.serialize();
			txservice.addCommand(cordacmd);
			
			CordaFlowContract contract = (CordaFlowContract)Class.forName(contractClass).newInstance();
			List<Pair<String, DocumentContext>> outputs = contract.runCommand(cordacmd, Lists.newArrayList());
			outputs.forEach(o -> {
         	 	try {
						ContractState s = (ContractState) CordaUtil.deserialize(o.getSecond().jsonString(), Class.forName(o.getFirst()));
						txservice.addOutputState(s);
         	 		} catch (ClassNotFoundException e) {
						throw new RuntimeException(e);
					}
			 });
			
			
		} catch (Exception e) {
			throw new IllegalArgumentException(e);
		} 
	}

}
