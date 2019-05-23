/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package general.activity.txnreply;

import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.core.runtime.engine.Context;
import com.tibco.dovetail.core.runtime.activity.IActivity;

public class txnreply implements IActivity {

    public void eval(Context context) throws IllegalArgumentException{
    		String status = context.getInput("status").toString();
    		String message = null;
    		String payload = null;
    		
    		if(status.equalsIgnoreCase("Error with Message")){
				Object omessage = context.getInput("message");
				if(omessage instanceof DocumentContext)
					message = ((DocumentContext)omessage).jsonString();
				else
					message = omessage.toString();
    			throw new IllegalArgumentException(message);
    		}
    	
    		
    		if(status.equalsIgnoreCase("Success With Data")) {
	    		DocumentContext doc = null;
	    		if(context.getInput("input") != null) {
	    			 doc = (DocumentContext) context.getInput("input");
	    		} else if(context.getInput("userInput") != null) {
	   			 doc = (DocumentContext) context.getInput("userInput");
	   		}
    		
	    		if (doc != null)
	    			payload = doc.jsonString();
    		}
    		
    		context.getReplyHandler().setReply(status, message, payload);
    		return;
    }
}
