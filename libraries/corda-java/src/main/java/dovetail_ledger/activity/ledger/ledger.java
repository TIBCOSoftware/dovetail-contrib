/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package dovetail_ledger.activity.ledger;

import java.util.List;

import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.core.runtime.engine.Context;
import com.tibco.dovetail.core.runtime.activity.IActivity;
import com.tibco.dovetail.core.runtime.services.IContainerService;
import com.tibco.dovetail.core.runtime.services.IDataService;
import com.tibco.dovetail.core.runtime.util.JsonUtil;

public class ledger implements IActivity {
    public void eval(Context context) throws IllegalArgumentException {
    		String op = context.getInput("operation").toString();
    		String assetName = context.getInput("asset").toString();
    		String txn = context.getInput("txn").toString();
    	
    		Object input = context.getInput("input");
    		if(input == null)
    			throw new IllegalArgumentException("Input is not set");
    		
    		boolean isArray = false;
    		if (context.getInput("isArray") != null)
    			isArray = Boolean.valueOf(context.getInput("isArray").toString());
    		
    		Object ctnr = context.getInput("containerServiceStub");
    		if (ctnr == null)
    			ctnr = context.getContainerService();
    		
    		if(ctnr == null)
    			throw new IllegalArgumentException("containerServicesStub is not mapped");
    		
    		IDataService service = ((IContainerService)ctnr).getDataService();
    		DocumentContext data = (DocumentContext)input;
    		DocumentContext output;
    		
    		switch(op) {
    		case "DELETE":
    			output = put(service, isArray, data, assetName, "linerId", txn);
    			break;
    		case "PUT":
    			output = put(service, isArray, data, assetName, "linerId", txn);
    			break;
    		case "GET":
    			output = get(service, isArray, data, assetName, "linearId");
    			break;
    		
    		default:
    			throw new IllegalArgumentException("Unsupported operation:" + op);
    		}
    		
    		context.setOutput("output", output);
    }
    
    private DocumentContext delete(IDataService service, boolean isArray, DocumentContext data, String assetName, String identifier, String txn) {
    		
    		if(isArray) {
    			DocumentContext doc = JsonUtil.getJsonParser().parse("[]");
			((List<Object>)data.json()).forEach( v -> {
				DocumentContext value = JsonUtil.getJsonParser().parse(v);
				DocumentContext obj = (DocumentContext) service.deleteState(txn , assetName, identifier , value);
				if (obj != null) {
					doc.add("$", obj.json());
				}
			});
			return doc;
		} else {
			return (DocumentContext) service.deleteState(txn, assetName, identifier , data);
		}
    }
    
    private DocumentContext get(IDataService service, boolean isArray, DocumentContext data, String assetName, String identifier) {
		
    		if(isArray) {
    			DocumentContext doc = JsonUtil.getJsonParser().parse("[]");
			((List<Object>)data.json()).forEach( v -> {
				DocumentContext value = JsonUtil.getJsonParser().parse(v);
				DocumentContext obj = (DocumentContext) service.getState(assetName, identifier , value);
				if (obj != null) {
					doc.add("$", obj.json());
				}
			});
			return doc;
		} else {
			return (DocumentContext) service.getState(assetName, identifier , data);
		}
    }
    
    private DocumentContext put(IDataService service, boolean isArray, DocumentContext data, String assetName, String identifier, String txn) {
		if(isArray) {
			((List<Object>)data.json()).forEach( v -> {
				DocumentContext value = JsonUtil.getJsonParser().parse(v);
				service.putState(txn, assetName, identifier , value);
			});
		} else {
			service.putState(txn, assetName, identifier , data);
		}
		
		return data;
    }
 
}
