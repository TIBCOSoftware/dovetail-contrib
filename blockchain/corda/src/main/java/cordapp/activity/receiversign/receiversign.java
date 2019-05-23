package cordapp.activity.receiversign;

import com.tibco.dovetail.container.cordapp.AppContainer;
import com.tibco.dovetail.container.cordapp.flows.DefaultSignTransactionFlow;
import com.tibco.dovetail.core.runtime.activity.IActivity;
import com.tibco.dovetail.core.runtime.engine.Context;

import co.paralleluniverse.fibers.Suspendable;
import net.corda.core.flows.FlowException;
import net.corda.core.flows.SignTransactionFlow;
import net.corda.core.transactions.SignedTransaction;

public class receiversign implements IActivity {

	@Override
	//@Suspendable
	public void eval(Context context) throws IllegalArgumentException {
		

	}

}
