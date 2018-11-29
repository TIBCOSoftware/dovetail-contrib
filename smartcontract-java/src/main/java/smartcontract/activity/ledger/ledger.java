/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package smartcontract.activity.ledger;

import java.util.List;

import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.core.runtime.engine.Context;
import com.tibco.dovetail.core.runtime.activity.IActivity;
import com.tibco.dovetail.core.runtime.services.IContainerService;
import com.tibco.dovetail.core.runtime.services.IDataService;
import com.tibco.dovetail.core.runtime.util.JsonUtil;

import net.minidev.json.JSONArray;

public class ledger implements IActivity {
    public void eval(Context context) throws IllegalArgumentException {
    		String op = context.getInput("operation").toString();
    		String assetName = context.getInput("assetName").toString();
    		String identifier = context.getInput("identifier").toString();
    	
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
    			output = delete(service, isArray, data, assetName, identifier);
    			break;
    		case "PUT":
    			output = put(service, isArray, data, assetName, identifier);
    			break;
    		case "GET":
    			output = get(service, isArray, data, assetName, identifier);
    			break;
    		case "LOOKUP":
    			Object compositeKey = context.getInput("compositeKey");
    			if(compositeKey == null)
    				throw new IllegalArgumentException("Composite key is not set for ledger LOOKUP operation");
        		
    			output = lookup(service, isArray, data, assetName, compositeKey.toString());
    			break;
    		default:
    			throw new IllegalArgumentException("Unsupported operation:" + op);
    		}
    		
    		context.setOutput("output", output);
    }
    
    private DocumentContext delete(IDataService service, boolean isArray, DocumentContext data, String assetName, String identifier) {
    		
    		if(isArray) {
    			DocumentContext doc = JsonUtil.getJsonParser().parse("[]");
			((JSONArray)data.json()).forEach( v -> {
				DocumentContext value = JsonUtil.getJsonParser().parse(v);
				DocumentContext obj = service.deleteState(assetName, identifier , value);
				if (obj != null) {
					doc.add("$", obj.json());
				}
			});
			return doc;
		} else {
			return service.deleteState(assetName, identifier , data);
		}
    }
    
    private DocumentContext get(IDataService service, boolean isArray, DocumentContext data, String assetName, String identifier) {
		
    		if(isArray) {
    			DocumentContext doc = JsonUtil.getJsonParser().parse("[]");
			((JSONArray)data.json()).forEach( v -> {
				DocumentContext value = JsonUtil.getJsonParser().parse(v);
				DocumentContext obj = service.getState(assetName, identifier , value);
				if (obj != null) {
					doc.add("$", obj.json());
				}
			});
			return doc;
		} else {
			return service.getState(assetName, identifier , data);
		}
    }
    
    private DocumentContext put(IDataService service, boolean isArray, DocumentContext data, String assetName, String identifier) {
		if(isArray) {
			((JSONArray)data.json()).forEach( v -> {
				DocumentContext value = JsonUtil.getJsonParser().parse(v);
				service.putState(assetName, identifier , value);
			});
		} else {
			service.putState(assetName, identifier , data);
		}
		
		return data;
    }
    
    private DocumentContext lookup(IDataService service, boolean isArray, DocumentContext data, String assetName, String identifier) {
	    	if(isArray) {
				DocumentContext doc = JsonUtil.getJsonParser().parse("[]");
			((JSONArray)data.json()).forEach( v -> {
				DocumentContext value = JsonUtil.getJsonParser().parse(v);
				List<DocumentContext> objs = service.lookupState(assetName, identifier , value);
				for (DocumentContext obj : objs) {
					doc.add("$", obj.json());
				}
			});
			return doc;
		} else {
			DocumentContext doc = JsonUtil.getJsonParser().parse("{}");
			Object obj = service.lookupState(assetName, identifier , data);
			if (obj != null) {
				doc.add("$", obj);
			}
			return doc;
		}
    }
}
