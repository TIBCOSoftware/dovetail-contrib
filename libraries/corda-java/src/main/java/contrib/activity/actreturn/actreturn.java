/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package contrib.activity.actreturn;

import com.tibco.dovetail.core.runtime.engine.Context;

import java.util.Map;

import com.tibco.dovetail.core.runtime.activity.IActivity;

public class actreturn implements IActivity {

    public void eval(Context context) throws IllegalArgumentException{
    		Map<String, Object> data = context.getInputs();
    		if(data != null)
    			context.getReplyHandler().setReply(data);
    		else
    			throw new IllegalArgumentException("reply data is not set");
    		
    		return;
    }
}
