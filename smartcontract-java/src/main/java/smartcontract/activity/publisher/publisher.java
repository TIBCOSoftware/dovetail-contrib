package smartcontract.activity.publisher;

import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.core.runtime.engine.Context;
import com.tibco.dovetail.core.runtime.activity.IActivity;
import com.tibco.dovetail.core.runtime.services.IContainerService;

public class publisher implements IActivity {

    public void eval(Context context) throws IllegalArgumentException{
    		String evtName = context.getInput("event").toString();
    		String evtMeta = "";
    		Object objevtMeta = context.getInput("eventMetadata");
    		if(objevtMeta != null)
    			evtMeta = objevtMeta.toString();
    		
    		Object ctnr = context.getInput("containerServiceStub");
    		if (ctnr == null)
    			ctnr = context.getContainerService();
    		
    		if(ctnr == null)
    			throw new IllegalArgumentException("containerServicesStub is not mapped");
    		
    		Object payload = context.getInput("inout");
    		if(payload != null) {
    			String json = ((DocumentContext)payload).jsonString();
    			 ((IContainerService)ctnr).getEventService().publish(evtName, evtMeta, json);
    		} else {
    			((IContainerService)ctnr).getEventService().publish(evtName, evtMeta, null);
    		}
    }
}
