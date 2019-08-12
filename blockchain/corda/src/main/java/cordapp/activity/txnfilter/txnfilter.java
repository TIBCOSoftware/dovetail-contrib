package cordapp.activity.txnfilter;

import java.util.List;
import java.util.function.Predicate;
import java.util.stream.Collectors;

import com.google.common.collect.Sets;
import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.container.corda.CordaUtil;
import com.tibco.dovetail.core.runtime.activity.IActivity;
import com.tibco.dovetail.core.runtime.engine.Context;
import com.tibco.dovetail.core.runtime.services.IDataService;
import com.tibco.dovetail.core.runtime.util.JsonUtil;

import net.corda.core.contracts.ContractState;
import net.corda.core.contracts.TimeWindow;
import net.corda.core.transactions.LedgerTransaction;
import net.corda.core.transactions.SignedTransaction;

public class txnfilter implements IActivity {

	@Override
	public void eval(Context context) throws IllegalArgumentException {
		String filter = context.getInput("filterby").toString();
		Object obj = context.getInput("ledgerTxn");
		if (obj == null)
			throw new IllegalArgumentException("txnfilter::ledgerTxn is required");
		
		Object asset = context.getInput("assetName");
		SignedTransaction txn = (SignedTransaction)obj;
		IDataService dataService = context.getContainerService().getDataService();
		int size = 0;
		DocumentContext doc = null;
		
		switch(filter) {
		case "Input State":
			List<ContractState> inputs = dataService.getStates(asset.toString(), null, Sets.newHashSet(txn.getInputs()));
			doc = CordaUtil.getInstance().toJsonObject(inputs);
			size = inputs.size();
			break;
		case "Reference State":
			List<ContractState> refs = dataService.getStates(asset.toString(), null, Sets.newHashSet(txn.getReferences()));
			doc = CordaUtil.getInstance().toJsonObject(refs);
			size = refs.size();
			break;
		case "Output State":
			List<ContractState> outs = txn.getTx().getOutputStates().stream().filter(new StateFilter(asset.toString())).collect(Collectors.toList());
			doc = CordaUtil.getInstance().toJsonObject(outs);
			size = outs.size();
			break;
		case "Command":
			List<String> cmds = txn.getTx().getCommands().stream().map(c -> c.getClass().getName()).filter( c -> c.equals(asset.toString())).collect(Collectors.toList());
			doc = JsonUtil.getJsonParser().parse("{}");
			if(cmds.size() > 0)
				doc.put("$", "command", cmds.get(0));
			
			size = 1;
			break;
		case "Notary":
			String notary = txn.getNotary().getName().getCommonName();
			doc = JsonUtil.getJsonParser().parse("{}");
			doc.put("$", "notary", notary);
			size = 1;
			break;
		case"Time Window":
			TimeWindow tw = txn.getTx().getTimeWindow();
			if (tw != null) {
				doc = JsonUtil.getJsonParser().parse("{}");
				doc.put("$", "from", tw.getFromTime()== null? "":tw.getFromTime().toString());
				doc.put("$", "until", tw.getUntilTime()== null? "":tw.getUntilTime().toString());
				doc.put("$", "duration", tw.getLength()==null? "": tw.getLength().toString());
				size = 1;
			}
			break;
		}
		
		context.setOutput("output", doc);
		context.setOutput("size", size);
	}
	
	public static class StateFilter implements Predicate<ContractState> {
		String assetName;
		StateFilter(String asset){
			assetName = asset;
		}
		@Override
		public boolean test(ContractState t) {
			return t.getClass().getName().equals(this.assetName)?true:false;
		}
		
	}

}
