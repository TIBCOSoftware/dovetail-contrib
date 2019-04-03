package cordapp.activity.receiversign;

import com.tibco.dovetail.container.cordapp.AppContainer;
import com.tibco.dovetail.container.cordapp.AppFlowService;
import com.tibco.dovetail.container.cordapp.DefaultSignTransactionFlow;
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
		/*boolean verify = false;
		if(context.getInput("verify") != null)
			Boolean.valueOf(context.getInput("verify").toString());
		
		AppFlowService flowservice = ((AppContainer) context.getContainerService()).getFlowService();
		SignTransactionFlow verifyFlow = null;
		if(!verify)
			verifyFlow = new DefaultSignTransactionFlow(flowservice.getFlowOtherParty());
		
		try {
			flowservice.receiverVerifyAndSign(verifyFlow);
			//SignedTransaction tx = 
			//context.setOutput("SignedTransaction", tx);
		} catch (FlowException e) {
			throw new IllegalArgumentException(e);
		}*/

	}

}
