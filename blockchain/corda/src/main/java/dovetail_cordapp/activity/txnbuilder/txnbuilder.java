package dovetail_cordapp.activity.txnbuilder;

import java.security.PublicKey;
import java.util.ArrayList;
import java.util.HashSet;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Set;

import com.google.common.collect.Lists;
import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.container.corda.ContractCommandOutput;
import com.tibco.dovetail.container.corda.CordaCommandDataWithData;
import com.tibco.dovetail.container.corda.CordaFlowContract;
import com.tibco.dovetail.container.corda.CordaUtil;
import com.tibco.dovetail.container.cordapp.AppDataService;
import com.tibco.dovetail.container.cordapp.AppFlow;
import com.tibco.dovetail.container.cordapp.AppUtil;
import com.tibco.dovetail.corda.json.deserializer.LinearIdDeserializer;
import com.tibco.dovetail.corda.json.serializer.LinearIdSerializer;
import com.tibco.dovetail.core.model.metadata.Attribute;
import com.tibco.dovetail.core.model.metadata.ResourceDescriptor;
import com.tibco.dovetail.core.runtime.activity.IActivity;
import com.tibco.dovetail.core.runtime.engine.Context;

import net.corda.core.contracts.ContractState;
import net.corda.core.contracts.StateAndRef;

public class txnbuilder implements IActivity{
	@Override
	public void eval(Context context) throws IllegalArgumentException {
		
		String command = context.getInput("command").toString();
		ResourceDescriptor schema = (ResourceDescriptor) context.getSetting("input_metadata");
		String contractClass = schema.getMetadata().getAsset() + "Contract";
		Set<PublicKey> signingkeys = new HashSet<PublicKey>();
		
		DocumentContext input = (DocumentContext)context.getInput("input");
		
		try {
			
			AppFlow txservice = (AppFlow) context.getContainerService().getContainerProperty("FlowService");
			AppDataService dataservice = (AppDataService) context.getContainerService().getDataService();
			
			LinkedHashMap<String, Object> inputs = (LinkedHashMap<String, Object>)input.json();
			LinkedHashMap<String, Set<PublicKey>> signingkeysbystate = new LinkedHashMap<String, Set<PublicKey>>();
					
			CordaCommandDataWithData cordacmd = new CordaCommandDataWithData();
			
			for(Attribute attr : schema.getAttributes()) {
				Object value = inputs.get(attr.getName());

				if(value != null) {
					if(attr.isAsset()) {
						if(attr.isArray()) {
							List<Object> results = new ArrayList<Object>();
							((List)value).forEach(v -> results.add(processAsset(dataservice, txservice, attr, v, signingkeys)));
							value = results;
							
						} else {
							value = processAsset(dataservice, txservice, attr, value, signingkeys);
						}
						
						if(attr.isRef()) {
							signingkeysbystate.put(attr.getType(), new HashSet<PublicKey>(signingkeys));
						}
					}
				} 
				
				cordacmd.putData(attr.getName(), value);
			}
			
			cordacmd.putData("command", command);
			cordacmd.serialize();
			
			CordaFlowContract contract = (CordaFlowContract)Class.forName(contractClass).newInstance();
			//output: <assetname, assetvalue, command>
			ContractCommandOutput outputs = contract.runCommand(cordacmd, Lists.newArrayList());
			outputs.getOutputStates().forEach(out -> {
				try {
					ContractState s = (ContractState) AppUtil.deserialize(out.getSecond().jsonString(), Class.forName(out.getFirst()));
					txservice.addOutputState(s);
					s.getParticipants().forEach(p -> signingkeys.add(p.getOwningKey()));
					
					//work around to get signatures for embedded command other than cash
					/*if(out.getThird() instanceof CordaCommandDataWithData) {
						CordaCommandDataWithData cmd = (CordaCommandDataWithData)out.getThird();
						if(!cmd.getCommand().equals(command)) {
							Set<PublicKey> keys = new HashSet<PublicKey>();
							Set<PublicKey> inkeys = signingkeysbystate.get(out.getFirst());
							if(inkeys != null)
								keys.addAll(inkeys);
							s.getParticipants().forEach(p -> keys.add(p.getOwningKey()));
							txservice.addCommand(cmd, keys);
							System.out.println("txbuilder::build commands, cmd=" + cmd.getCommand() + ", keys=" + CordaUtil.getInstance().serialize(keys) );
						}
					}*/
					
     	 		} catch (ClassNotFoundException e) {
					throw new RuntimeException(e);
				}
			});
			
			//a smart contract transaction (command) could invoke another command, each command has its own set of signing keys
			outputs.getEmbeddedCommands().forEach((c, keys) -> {
				if(c instanceof CordaCommandDataWithData)
					((CordaCommandDataWithData)c).serialize();
				
				System.out.println("txbuilder::build commands, cmd=" + c.getClass().getSimpleName() + ", keys=" + CordaUtil.getInstance().serialize(keys) );
				txservice.addCommand(c, keys);
			});
			
			System.out.println("txbuilder::build commands, cmd=" + cordacmd.getCommand() + ", keys=" + CordaUtil.getInstance().serialize(signingkeys) );
			txservice.addCommand(cordacmd, signingkeys);
		} catch (Exception e) {
			throw new IllegalArgumentException(e);
		} 
	}
	
	private Object processAsset(AppDataService dataservice, AppFlow txservice, Attribute attr, Object value, Set<PublicKey> keys) {
		LinkedHashMap assetvalue = (LinkedHashMap)value;
		if (attr.isRef()) {
			StateAndRef state  = dataservice.getStateRef(assetvalue.get("ref").toString());
			if(attr.isReferenceData()) {
				txservice.addRefState(state);
			} else {
				txservice.addInputState(state);
				state.getState().getData().getParticipants().forEach(p -> keys.add(p.getOwningKey()));
			}
			return state.getState().getData();
		} else {
			//fix linear id as workaround until functions are available
			if(assetvalue.get("linearId") != null){
				assetvalue.put("linearId", LinearIdSerializer.toString(LinearIdDeserializer.fromString(assetvalue.get("linearId").toString())));
			}	
			return assetvalue;
		}
	}

}
