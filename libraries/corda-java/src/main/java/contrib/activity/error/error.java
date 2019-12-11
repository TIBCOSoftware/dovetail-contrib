/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package contrib.activity.error;

import com.tibco.dovetail.core.runtime.engine.Context;
import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.core.runtime.activity.IActivity;

public class error implements IActivity {

    public void eval(Context context) throws IllegalArgumentException{
    		String message = context.getInput("message").toString();
    		Object data = context.getInput("data");
    		if(data != null) {
    			String detail = null;
    			if(data instanceof DocumentContext)
    				detail = ((DocumentContext)data).jsonString();
    			else
    				detail = data.toString();
    			throw new IllegalArgumentException(message + ". error data:" + detail);
    		} else {
    			throw new IllegalArgumentException(message);
    		}
    }
}
