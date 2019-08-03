package cordapp.activity.subflow;

import java.util.LinkedHashMap;

import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.container.corda.CordaUtil;
import com.tibco.dovetail.container.cordapp.AppContainer;
import com.tibco.dovetail.corda.json.LinearIdDeserializer;
import com.tibco.dovetail.corda.json.MoneyAmtDeserializer;
import com.tibco.dovetail.core.model.composer.HLCAttribute;
import com.tibco.dovetail.core.model.composer.HLCResource;
import com.tibco.dovetail.core.runtime.activity.IActivity;
import com.tibco.dovetail.core.runtime.engine.Context;

public class subflow implements IActivity {

	@Override
	public void eval(Context context) throws IllegalArgumentException {
		//parse out flow input arguments
		Object input = context.getInput("input");
		if(input != null) {
			LinkedHashMap inputvalues = ((DocumentContext)input).json();
			
			HLCResource flowInputMetadata = (HLCResource) context.getInput("input_metadata");
			
			LinkedHashMap<String, Object> task = new LinkedHashMap<String, Object>();
			task.put("FlowName", context.getInput("flowName"));
			
			LinkedHashMap<String, Object> flowparams = new LinkedHashMap<String, Object>();
			for(int i=0; i<flowInputMetadata.getAttributes().size(); i++) {
				HLCAttribute attr = flowInputMetadata.getAttributes().get(i);
				switch(attr.getType()) {
				case "net.corda.core.identity.Party":
					//flowparams.put(attr.getName(), CordaUtil.partyFromString(inputvalues.get(attr.getName()).toString()));
					flowparams.put(attr.getName(), CordaUtil.partyFromCommonName(inputvalues.get(attr.getName()).toString()));
					break;
				case "net.corda.core.contracts.Amount<Currency>":
					flowparams.put(attr.getName(), MoneyAmtDeserializer.parseAmount(CordaUtil.toJsonNode(inputvalues.get(attr.getName()))));
					break;
				case "net.corda.core.contracts.UniqueIdentifier":
					flowparams.put(attr.getName(), LinearIdDeserializer.fromString(inputvalues.get(attr.getName()).toString()));
					break;
				case "java.time.LocalDate":
					flowparams.put(attr.getName(), java.time.LocalDate.parse(inputvalues.get(attr.getName()).toString()));
					break;
				case "java.time.Instant":
					flowparams.put(attr.getName(), java.time.Instant.parse(inputvalues.get(attr.getName()).toString()));
					break;
				default:
					flowparams.put(attr.getName(),  inputvalues.get(attr.getName()));
				}
			}
			task.put("Arguments", flowparams);
		
			context.getContainerService().addContainerAsyncTask(AppContainer.TASK_SUBFLOW, task);

		} else {
			throw new IllegalArgumentException("input is required for subflow activity");
		}
		

	}

}
