package cordapp.activity.signandcommit;

import com.tibco.dovetail.container.cordapp.AppContainer;
import com.tibco.dovetail.container.cordapp.AppFlowService;
import com.tibco.dovetail.core.runtime.activity.IActivity;
import com.tibco.dovetail.core.runtime.engine.Context;

import co.paralleluniverse.fibers.Suspendable;
import net.corda.core.flows.FlowException;

public class signandcommit implements IActivity{
	@Override
	//@Suspendable
	public void eval(Context context) throws IllegalArgumentException {
		/*AppFlowService txnservice = ((AppContainer) context.getContainerService()).getFlowService();
		try {
			txnservice.signAndCommit();
		} catch (FlowException e) {
			throw new IllegalArgumentException(e);
		}*/
	}
}
