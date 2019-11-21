package dovetail_ledger.activity.payment;


import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.core.runtime.activity.IActivity;
import com.tibco.dovetail.core.runtime.engine.Context;

public class payment implements IActivity {

	@Override
	public void eval(Context ctx) throws IllegalArgumentException {
		DocumentContext input = (DocumentContext) ctx.getInput("input");
		ctx.getContainerService().getDataService().processPayment(input);
	}

}
