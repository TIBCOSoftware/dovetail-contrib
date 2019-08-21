/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
import com.fasterxml.jackson.databind.ObjectMapper;
import com.jayway.jsonpath.DocumentContext;
import com.jayway.jsonpath.JsonPath;
import com.tibco.dovetail.core.model.activity.ActivityModel;
import com.tibco.dovetail.core.runtime.util.CompareUtil;
import com.tibco.dovetail.core.runtime.util.JsonUtil;

import net.minidev.json.JSONArray;

import org.junit.Test;

import java.io.InputStream;
import java.util.ArrayList;
import java.util.LinkedHashMap;
import java.util.List;

public class TestResource {
    @Test
    public void testActivityJson(){
    	System.out.println("testActivityJson.......");
        try {
           ObjectMapper mapper = new ObjectMapper();
            ActivityModel model = ActivityModel.loadModel(mapper, "smartcontract/activity/ledger");
            System.out.println(model.getRef());
        }catch (Exception e){e.printStackTrace();}
    }
    
    @Test
    public void testDocumentContextArray() {
    	System.out.println("testDocumentContextArray.......");
    		String json = "[{\"a\": \"av\", \"b\":\"bv\"},{\"a\": \"av2\", \"b\":\"bv2\"}]";
    		
    		DocumentContext doc = JsonPath.parse(json);
    		List<Object> data = doc.json();
    		
    		for(int i=0; i<data.size(); i++) {
    			LinkedHashMap<String, Object> d = (LinkedHashMap<String, Object>) data.get(i);
    			d.put("_index_", i);
    		}
    		System.out.println("count=" + data.size());
    		
    		System.out.println(doc.jsonString());
    		
    		
    }
    
    @Test
    public void testPrimitiveToObjectArray() {
     	System.out.println("testPrimitiveToObjectArray.......");
    		String json = "[\"A\", \"B\"]";
		
		DocumentContext doc = JsonPath.parse(json);
		
    		DocumentContext doc2 = JsonUtil.getJsonParser().parse("[]");
		List<Object> data = doc.json();
		data.forEach(d -> {
			LinkedHashMap<String, Object> value = new LinkedHashMap<String, Object>();
			value.put("field", d);
			doc2.add("$", value);
		});
		
		System.out.println(doc2.jsonString());
    }
    
    @Test
    public void testObjectToPrimitiveArray() {
    	System.out.println("testObjectToPrimitiveArray.......");
    	String json = "[{\"field\": \"av\"},{\"field\": \"av2\"}]";
		
		DocumentContext doc = JsonPath.parse(json);
		
    		DocumentContext doc2 = JsonUtil.getJsonParser().parse("[]");
		List<Object> data = doc.json();
		data.forEach(d -> {
			LinkedHashMap<String, Object> value = (LinkedHashMap<String, Object>)d;
			doc2.add("$", value.get("field"));
		});
		
		System.out.println(doc2.jsonString());
    }
    
    @Test
    public void testMerge() {
    		System.out.println("testMerge.......");
    		String json = "{\"input1\":[{\"field\": \"av\"},{\"field\": \"av2\"}], \"input2\":[{\"field\": \"abv\"},{\"field\": \"abv2\"}]}";
		
		DocumentContext input = JsonPath.parse(json);
		
		DocumentContext doc = JsonUtil.getJsonParser().parse("[]");
		List<Object> input1 = input.read("input1");
		List<Object> input2 = input.read("input2");
		
		if (input1 != null && input1.size() > 0) {
			input1.forEach(in -> doc.add("$", in));
		}
		
		if (input2 != null && input2.size() > 0) {
			input2.forEach(in -> doc.add("$", in));
		}
		
		System.out.println(doc.jsonString());
    }
    
    @Test
    public void testFilter() {
    		System.out.println("testFilter.......");
    		String json = "{\"filterValue\":\"av\", \"dataset\":[{\"field\": {\"field2\":\"av\"}},{\"field\": {\"field2\":\"av2\"}}]}";
		
		DocumentContext input = JsonPath.parse(json);
		
		DocumentContext doc = JsonUtil.getJsonParser().parse("{}");
		List<Object> trueset = new ArrayList<Object>();
		List<Object> falseset = new ArrayList<Object>();
		
		String field = "$dataset.field.field2";
		String op = "==";
		
		LinkedHashMap<String, Object> data = input.json();
		Object filterValue = data.get("filterValue");
		Object values = data.get("dataset");
		
		if (values != null) {
			//$dataset.path.to.field
			String[] tokens = field.split("\\.");
			if(tokens.length < 2 || !field.startsWith("$dataset")) {
				throw new IllegalArgumentException("collection filter field " + field + " should be in the format of $dataset.path.to.field");
			}
			
			List<LinkedHashMap<String, Object>> dataset = (List<LinkedHashMap<String, Object>>)values;
			dataset.forEach(in -> {
				Object v = in.get(tokens[1]);
				for(int i= 2; i<tokens.length; i++) {
					if(v != null) {
						v = ((LinkedHashMap<String, Object>)v).get(tokens[i]);
					} else {
						break;
					}
				}
				
				if (v == null || !CompareUtil.compare(v, filterValue, op)) {
					falseset.add(in);
				} else {
					trueset.add(in);
				}
			});
		}
		
		doc.put(JsonPath.compile("$"), "trueSet", trueset);
		doc.put(JsonPath.compile("$"), "falseSet", falseset);

		
		System.out.println(doc.jsonString());
    }
}
