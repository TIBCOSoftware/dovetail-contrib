package cordapp.activity.initiatorsign;

import com.tibco.dovetail.core.runtime.activity.IActivity;
import com.tibco.dovetail.core.runtime.engine.Context;

import co.paralleluniverse.fibers.Suspendable;
import net.corda.core.node.ServiceHub;
import net.corda.core.transactions.TransactionBuilder;

public class initiatorsign implements IActivity{

	@Override
	@Suspendable
	public void eval(Context context) throws IllegalArgumentException {
		TransactionBuilder builder = (TransactionBuilder) context.getInput("_txnbuilder_");
		ServiceHub services = (ServiceHub) context.getInput("_servicehub_");
		services.signInitialTransaction(builder);
	}

}
