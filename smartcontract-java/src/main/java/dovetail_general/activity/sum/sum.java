/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package dovetail_general.activity.sum;

import java.math.RoundingMode;
import java.util.List;

import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.core.runtime.engine.Context;
import com.tibco.dovetail.function.math;
import com.tibco.dovetail.core.runtime.activity.IActivity;
import com.tibco.dovetail.core.runtime.util.JsonUtil;

public class sum implements IActivity {

	@Override
	public void eval(Context context) throws IllegalArgumentException {
		
		int scale = Integer.valueOf(context.getInput("scale").toString());
		String rounding = context.getInput("rounding").toString();
		
		//TODO: groupby
		//boolean groupby = Boolean.valueOf(context.getInput("groupby").toString());
		
		DocumentContext input = (DocumentContext)context.getInput("input");
		DocumentContext doc = JsonUtil.getJsonParser().parse("{}");
		
		List<Object> data = input.read("$.dataset[*].field");
		
		if(data.size() == 0) {
			doc.put("$", "result", 0);
		} else {
			Object result = math.sumBigDecimal(RoundingMode.valueOf(rounding), scale, data.toArray());
			doc.put("$", "result", result);
		}
	
		context.setOutput("output", doc);
	}

}
