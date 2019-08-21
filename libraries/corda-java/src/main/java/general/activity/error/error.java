/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package general.activity.error;

import com.tibco.dovetail.core.runtime.engine.Context;
import com.tibco.dovetail.core.runtime.activity.IActivity;

public class error implements IActivity {

    public void eval(Context context) throws IllegalArgumentException{
    		String message = context.getInput("message").toString();
    		
    		throw new IllegalArgumentException(message);
    		
    }
}
