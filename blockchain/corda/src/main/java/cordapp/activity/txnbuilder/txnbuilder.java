package cordapp.activity.txnbuilder;

import java.io.Serializable;
import java.util.ArrayList;
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
import com.tibco.dovetail.corda.json.LinearIdDeserializer;
import com.tibco.dovetail.corda.json.LinearIdSerializer;
import com.tibco.dovetail.core.runtime.activity.IActivity;
import com.tibco.dovetail.core.runtime.engine.Context;

import kotlin.Pair;
import kotlin.Triple;
import net.corda.core.contracts.CommandData;
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
					if(attr.isAsset()) {
						if(attr.isArray()) {
							List<Object> results = new ArrayList<Object>();
							((List)value).forEach(v -> results.add(processAsset(dataservice, txservice, attr.isRef(), v)));
							value = results;
							
						} else {
							value = processAsset(dataservice, txservice, attr.isRef(), value);
						}
					}
				} 
				
				cordacmd.putData(attr.getName(), value);
			}
			
			cordacmd.putData("command", command);
			cordacmd.serialize();
		
			txservice.addCommand(cordacmd);
			
			CordaFlowContract contract = (CordaFlowContract)Class.forName(contractClass).newInstance();
			List<Triple<String, DocumentContext, CommandData>> outputs = contract.runCommand(cordacmd, Lists.newArrayList());
			outputs.forEach(o -> {
         	 	try {
						ContractState s = (ContractState) CordaUtil.deserialize(o.getSecond().jsonString(), Class.forName(o.getFirst()));
						txservice.addOutputState(s);
						
						if(o.getThird() != null) {
							if(o.getThird() instanceof CordaCommandDataWithData)
								((CordaCommandDataWithData)o.getThird()).serialize();
							
							txservice.addCommand(o.getThird());
						}
         	 		} catch (ClassNotFoundException e) {
						throw new RuntimeException(e);
					}
			 });
			
			
		} catch (Exception e) {
			throw new IllegalArgumentException(e);
		} 
	}
	
	private Object processAsset(AppDataService dataservice, AppFlow txservice, boolean isref, Object value) {
		if (isref) {
			StateAndRef state  = dataservice.getStateRef(value.toString());
			txservice.addInputState(state);
			return state.getState().getData();
		} else {
			//fix linear id as workaround until functions are available
			LinkedHashMap assetvalue = (LinkedHashMap)value;
			if(assetvalue.get("linearId") != null){
				assetvalue.put("linearId", LinearIdSerializer.toString(LinearIdDeserializer.fromString(assetvalue.get("linearId").toString())));
			}	
			return assetvalue;
		}
	}

}
