package cordapp.activity.simplevaultquery;

import java.util.Arrays;
import java.util.LinkedHashMap;
import java.util.List;

import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.container.cordapp.VaultQuery;
import com.tibco.dovetail.corda.json.deserializer.LinearIdDeserializer;
import com.tibco.dovetail.core.runtime.activity.IActivity;
import com.tibco.dovetail.core.runtime.engine.Context;
import com.tibco.dovetail.core.runtime.services.IDataService;

import net.corda.core.contracts.UniqueIdentifier;
import net.corda.core.node.services.Vault;
import net.corda.core.node.services.vault.QueryCriteria;

public class simplevaultquery implements IActivity {

	@Override
	public void eval(Context context) throws IllegalArgumentException {
		String asset = context.getInput("assetName").toString();
		String status = context.getInput("status").toString();
		String assetType = context.getInput("assetType").toString();
		IDataService dataservice = context.getContainerService().getDataService();
		
		switch (assetType) {
		case "LinearState":
			Object input = context.getInput("input");
			if (input == null)
				throw new IllegalArgumentException("simplevaultquery: input is not mapped");
			
			DocumentContext indoc = (DocumentContext)input;
			String id = ((LinkedHashMap)indoc.json()).get("linearId").toString();
			UniqueIdentifier linearId = LinearIdDeserializer.fromString(id);
			QueryCriteria.LinearStateQueryCriteria criteria = new QueryCriteria.LinearStateQueryCriteria(null, Arrays.asList(linearId.getId()), Arrays.asList(linearId.getExternalId()), Vault.StateStatus.valueOf(status));

			List<DocumentContext> output = dataservice.queryState(new VaultQuery(asset, criteria));
			context.setOutput("output", output);
			context.setOutput("size", output.size());
			
			break;
		case "FungibleAsset":
			break;
		}
	}

}
