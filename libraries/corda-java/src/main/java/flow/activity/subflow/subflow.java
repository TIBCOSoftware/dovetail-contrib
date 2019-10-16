package flow.activity.subflow;

import com.tibco.dovetail.core.runtime.activity.IActivity;
import com.tibco.dovetail.core.runtime.engine.Context;
import com.tibco.dovetail.core.runtime.engine.FlowEngine;
import com.tibco.dovetail.core.runtime.engine.Scope;
import com.tibco.dovetail.core.runtime.flow.BasicTransactionFlow;
import com.tibco.dovetail.core.runtime.flow.ReplyData;

public class subflow implements IActivity {

	@Override
	public void eval(Context context) throws IllegalArgumentException {
		String flowName = context.getSetting("flowURI").toString();
		Object flow = context.getSetting("subflow");
		if(flow != null) {
			FlowEngine engine = new FlowEngine((BasicTransactionFlow) flow);
			ReplyData reply = engine.execute(context, new Scope(null, true));
			if(reply != null && reply.getData() != null) {
				reply.getData().forEach((k,v) -> {
					context.setOutput(k, v);
				});
			}
		} else {
			throw new IllegalArgumentException("flow " + flowName + " does not exist");
		}
	}

}
