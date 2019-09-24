/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package general.activity.txnreply;

import com.tibco.dovetail.core.runtime.engine.Context;
import com.tibco.dovetail.core.runtime.activity.IActivity;

public class txnreply implements IActivity {

    public void eval(Context context) throws IllegalArgumentException{
    		Object data = context.getInput("data");
    		if(data != null)
    			context.getReplyHandler().setReply(data.toString());
    		else
    			throw new IllegalArgumentException("reply data is not set");
    		
    		return;
    }
}
