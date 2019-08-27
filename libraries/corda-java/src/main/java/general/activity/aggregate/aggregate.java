/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package general.activity.aggregate;

import java.math.BigDecimal;
import java.math.RoundingMode;
import java.util.List;
import java.util.stream.Collectors;

import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.core.runtime.engine.Context;
import com.tibco.dovetail.function.math;
import com.tibco.dovetail.core.runtime.activity.IActivity;
import com.tibco.dovetail.core.runtime.util.JsonUtil;

public class aggregate implements IActivity {

	@Override
	public void eval(Context context) throws IllegalArgumentException {
		String op = context.getInput("operation").toString();
		String datatype = context.getInput("dataType").toString();
		int scale = 0;
		int precision = 0;
		RoundingMode rounding = RoundingMode.HALF_EVEN;
		
		if(datatype.equals("Double") || op.equals("AVG")) {
			scale = Integer.valueOf(context.getInput("scale").toString());
			precision = Integer.valueOf(context.getInput("precision").toString());
			rounding = RoundingMode.valueOf(context.getInput("scale").toString());
		}
		
		DocumentContext input = (DocumentContext)context.getInput("input");
		List<Object> data = input.read("$..data");
		if(data.size() == 0)
			return;
	
		data = data.stream().map(d -> {
								switch(datatype) {
								case "Integer":
									return Integer.valueOf(d.toString());
								case "Long":
									return Long.valueOf(d.toString());
								default:
									return new BigDecimal(d.toString());
								}
							})
							.collect(Collectors.toList());
		
		Object result = null;
		
		switch(op) {
		case "MIN":
			result = math.min(rounding, precision, scale, data.toArray());
			break;
		case "MAX":
			result = math.max(rounding, precision, scale, data.toArray());
			break;
		case "AVG":
			result = math.avg(rounding, precision, scale, data.toArray());
			break;
		case "SUM":
			result = math.sum(rounding, precision, scale, data.toArray());
			break;
		}
		
		DocumentContext doc = JsonUtil.getJsonParser().parse("{}");
		doc.put("$", "result", result);
		context.setOutput("output", doc);
	}

}
