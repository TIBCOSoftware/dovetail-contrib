package cordapp.activity.txnbuilder;

import java.security.PublicKey;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Set;

import com.fasterxml.jackson.core.type.TypeReference;
import com.google.common.collect.Lists;
import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.container.corda.ContractCommandOutput;
import com.tibco.dovetail.container.corda.CordaCommandDataWithData;
import com.tibco.dovetail.container.corda.CordaFlowContract;
import com.tibco.dovetail.container.corda.CordaUtil;
import com.tibco.dovetail.container.cordapp.AppDataService;
import com.tibco.dovetail.container.cordapp.AppFlow;
import com.tibco.dovetail.corda.json.LinearIdDeserializer;
import com.tibco.dovetail.corda.json.LinearIdSerializer;
import com.tibco.dovetail.core.runtime.activity.IActivity;
import com.tibco.dovetail.core.runtime.engine.Context;

import net.corda.core.contracts.ContractState;
import net.corda.core.contracts.StateAndRef;

public class txnbuilder implements IActivity{
	@Override
	public void eval(Context context) throws IllegalArgumentException {
		
		String contractClass = context.getInput("contractClass").toString();
		String command = context.getInput("command").toString();
		String schema = context.getInput("inputSchema").toString();
		Set<PublicKey> signingkeys = new HashSet<PublicKey>();
		
		DocumentContext input = (DocumentContext)context.getInput("input");
		
		try {
			
			AppFlow txservice = (AppFlow) context.getContainerService().getContainerProperty("FlowService");
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
							((List)value).forEach(v -> results.add(processAsset(dataservice, txservice, attr.isRef(), v, signingkeys)));
							value = results;
							
						} else {
							value = processAsset(dataservice, txservice, attr.isRef(), value, signingkeys);
						}
					}
				} 
				
				cordacmd.putData(attr.getName(), value);
			}
			
			cordacmd.putData("command", command);
			cordacmd.serialize();
			
			CordaFlowContract contract = (CordaFlowContract)Class.forName(contractClass).newInstance();
			ContractCommandOutput outputs = contract.runCommand(cordacmd, Lists.newArrayList());
			outputs.getOutputStates().forEach(o -> {
				try {
					ContractState s = (ContractState) CordaUtil.deserialize(o.getSecond().jsonString(), Class.forName(o.getFirst()));
					txservice.addOutputState(s);
					
					//if the output state is created as part of another command, add signing keys to the executing command
					if(o.getThird() == null) {
						s.getParticipants().forEach(p -> signingkeys.add(p.getOwningKey()));
					}
					
     	 		} catch (ClassNotFoundException e) {
					throw new RuntimeException(e);
				}
			});
			
			//a smart contract transaction (command) could invoke another command, each command has its own set of signing keys
			outputs.getEmbeddedCommands().forEach((c, keys) -> {
				if(c instanceof CordaCommandDataWithData)
					((CordaCommandDataWithData)c).serialize();
				
				txservice.getLogger().info("txbuilder::build commands, cmd=" + c.getClass().getSimpleName() + ", keys=" + CordaUtil.serialize(keys) );
				txservice.addCommand(c, keys);
			});
			
			txservice.getLogger().info("txbuilder::build commands, cmd=" + cordacmd.getCommand() + ", keys=" + CordaUtil.serialize(signingkeys) );
			txservice.addCommand(cordacmd, signingkeys);
		} catch (Exception e) {
			throw new IllegalArgumentException(e);
		} 
	}
	
	private Object processAsset(AppDataService dataservice, AppFlow txservice, boolean isref, Object value, Set<PublicKey> keys) {
		if (isref) {
			StateAndRef state  = dataservice.getStateRef(value.toString());
			txservice.addInputState(state);
			state.getState().getData().getParticipants().forEach(p -> keys.add(p.getOwningKey()));
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
