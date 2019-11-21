/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package dovetail_general.activity.mapper;

import java.math.BigDecimal;
import java.math.MathContext;
import java.math.RoundingMode;
import java.util.Arrays;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.stream.Collectors;

import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.core.runtime.engine.Context;
import com.tibco.dovetail.core.runtime.activity.IActivity;
import com.tibco.dovetail.core.runtime.util.JsonUtil;

public class mapper implements IActivity {

	@Override
	public void eval(Context context) throws IllegalArgumentException {
		DocumentContext input = getInput(context);
		String datatype = context.getInput("dataType").toString();
		boolean isArray = false;
		String inputArrayType = "Object Array";
		String outputArrayType = "Object Array";
		
		if(context.getInput("isArray") != null) {
			isArray = Boolean.valueOf(context.getInput("isArray").toString());
			inputArrayType = context.getInput("inputArrayType").toString();
			outputArrayType = context.getInput("outputArrayType").toString();
			
		}
		
		if(datatype.equalsIgnoreCase("User Defined...")) {
			context.setOutput("output", input);
		} else {
			if(datatype.equalsIgnoreCase("Double")) {
				input = handleDouble(context, input);
			}
			
			if(isPrimitive(datatype) && isArray && !inputArrayType.equalsIgnoreCase(outputArrayType)) {
				DocumentContext doc = JsonUtil.getJsonParser().parse("[]");
				List<Object> data = input.json();
				if(outputArrayType.equalsIgnoreCase("Object Array")) {
					//convert primitive to object array
					data.forEach(d -> {
						LinkedHashMap<String, Object> value = new LinkedHashMap<String, Object>();
						value.put("field", d);
						doc.add("$", value);
					});
				} else {
					//convert object to primitive array
					data.forEach(d -> {
						LinkedHashMap<String, Object> value = (LinkedHashMap<String, Object>)d;
						doc.add("$", value.get("field"));
					});
				}
				context.setOutput("output", doc);
			} else {
				context.setOutput("output", input);
			}
		}
	}
	
	private DocumentContext getInput(Context context) {
		String dataType = context.getInput("dataType").toString();
		if (dataType.equalsIgnoreCase("User Defined...")) {
			Object obj = context.getInput("userInput");
			if(obj != null) {
				DocumentContext input = (DocumentContext)obj;
				return input;
			}
		} else {
			Object obj = context.getInput("input");
			if(obj != null) {
				if(obj instanceof DocumentContext) {
					DocumentContext input = (DocumentContext)obj;
					return input;
				} else {
					List inputs;
					if(obj instanceof List)
						inputs = (List)obj;
					else if(obj.getClass().isArray()) {
						inputs = Arrays.stream((Object[])obj).collect(Collectors.toList());
					}
					else {
						return JsonUtil.getJsonParser().parse(obj);
					}
			
					DocumentContext doc = JsonUtil.getJsonParser().parse("[]");
					
					for(Object input : inputs) {
						doc.add("$", input);
					}
					return doc;
				}
			}
		}

		throw new IllegalArgumentException("Input data is not set");
	}
	
	private boolean isPrimitive(String dataType) {
		switch(dataType) {
		case "Integer":
		case "Double":
		case "Long":
		case "String":
		case "Boolean":
		case "Datetime":
			return true;
		default:
			return false;
		}
	}
	
	private DocumentContext handleDouble(Context context, DocumentContext doc) {
		int scale = Integer.valueOf(context.getInput("scale").toString());
		int precision = Integer.valueOf(context.getInput("precision").toString());
		RoundingMode rounding = RoundingMode.valueOf(context.getInput("scale").toString());
		
		List<Object> data = doc.read("$..field");
		data.forEach(d -> {
			if(d instanceof BigDecimal) {
				BigDecimal bd = new BigDecimal(d.toString(), new MathContext(precision, rounding));
				bd.setScale(scale);
				doc.put("$", "result", ((BigDecimal)bd).toPlainString());
			}
			else
				doc.put("$", "result", d.toString());
		});
		
		return doc;	
	}
}
