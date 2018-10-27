package smartcontract.activity.query;

import java.util.LinkedHashMap;
import java.util.List;

import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.core.runtime.engine.Context;
import com.tibco.dovetail.core.runtime.activity.IActivity;

public class query implements IActivity {

    public void eval(Context context) throws IllegalArgumentException{
    		String queryString = context.getInput("queryString").toString();
    		if (queryString == null)
    			throw new IllegalArgumentException("query string is not set");
    		
    		Object objParams = context.getInput("param");
    		List<DocumentContext> result = null;
    		if (objParams == null) {
          result = context.getContainerService().getDataService().queryState(queryString);
    		}
    		else {
    			DocumentContext params = (DocumentContext)objParams;
    			LinkedHashMap<String, Object> pvs = params.json();
    			pvs.forEach((k,v) -> {
    				queryString.replaceAll("_$"+k, v.toString());
    			});
    		}
         context.setOutput("output", result);
    }
}
