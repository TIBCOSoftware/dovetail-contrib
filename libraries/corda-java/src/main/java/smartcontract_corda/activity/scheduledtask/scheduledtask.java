package smartcontract_corda.activity.scheduledtask;

import java.time.Instant;
import java.util.LinkedHashMap;

import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.container.corda.CordaContainer;
import com.tibco.dovetail.core.runtime.activity.IActivity;
import com.tibco.dovetail.core.runtime.engine.Context;
import net.corda.core.contracts.ScheduledActivity;
import net.corda.core.contracts.StateRef;
import net.corda.core.flows.FlowLogicRef;
import net.corda.core.flows.FlowLogicRefFactory;

public class scheduledtask implements IActivity {

	@Override
	public void eval(Context context) throws IllegalArgumentException {
		CordaContainer ctnr = (CordaContainer)context.getContainerService();
	
		Object input = context.getInput("input");
		if(input != null) {
			LinkedHashMap inputvalues = ((DocumentContext)input).json();
			if(inputvalues.get("scheduledAt") == null) {
				throw new IllegalArgumentException("scheduledAt is required for scheduled activity");
			}
			
			String scheduledtime = inputvalues.get("scheduledAt").toString();
			String flowName = context.getInput("schedulableFlowClassName").toString();
			
			String flowClass = ctnr.getContainerProperty("Namespace") + "." + flowName + "Impl";
			
			FlowLogicRefFactory factory = (FlowLogicRefFactory)ctnr.getContainerProperty("FlowLogicRefFactory");
			StateRef ref = (StateRef)ctnr.getContainerProperty("StateRef");
			FlowLogicRef flowRef = factory.create(flowClass, ref);
			ctnr.getLogService().info(flowClass + "is scheduled at " + scheduledtime + ", flowRef=" + factory.toFlowLogic(flowRef).getClass().getName());
		
			ctnr.addContainerAsyncTask(CordaContainer.TASK_SCHEDULEDACTIVITY, new ScheduledActivity(flowRef, Instant.parse(scheduledtime)));

		} else {
			throw new IllegalArgumentException("scheduledAt is required for scheduled activity");
		}
		
		
	}

}
