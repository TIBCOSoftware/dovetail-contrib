/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package smartcontract.activity.history;

import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.core.runtime.engine.Context;
import com.tibco.dovetail.core.runtime.activity.IActivity;
import com.tibco.dovetail.core.runtime.services.IContainerService;

public class history implements IActivity {

    public void eval(Context context) throws IllegalArgumentException{
    		String assetName = context.getInput("assetName").toString();
    		String assetKey = context.getInput("identifier").toString();
    		Object keyValue = context.getInput("input");
    		if(keyValue == null)
    			throw new IllegalArgumentException("Input is not set");
    		
    		Object ctnr = context.getInput("containerServiceStub");
    		if (ctnr == null)
    			ctnr = context.getContainerService();
    		
    		if(ctnr == null)
    			throw new IllegalArgumentException("containerServicesStub is not mapped");
    		
         ((IContainerService)ctnr).getDataService().getHistory(assetName, assetKey, (DocumentContext)keyValue);
    }
}
