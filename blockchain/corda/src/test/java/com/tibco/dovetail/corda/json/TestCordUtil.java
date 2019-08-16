package com.tibco.dovetail.corda.json;

import java.util.List;

import org.junit.Test;

import com.fasterxml.jackson.core.type.TypeReference;
import com.tibco.dovetail.container.corda.CordaUtil;

import cordapp.activity.txnbuilder.BuilderSchemaAttribute;

public class TestCordUtil {
	@Test
	public void testDeserializeAttrs() {
		String json = "[{\"name\":\"iou\",\"isOptional\":false,\"type\":\"com.example.iou.IOU\",\"isAsset\":true},{\"name\":\"transactionId\",\"isOptional\":false,\"type\":\"String\"},{\"name\":\"timestamp\",\"isOptional\":false,\"type\":\"DateTime\"}]";
		List<BuilderSchemaAttribute> attrs = (List<BuilderSchemaAttribute>) CordaUtil.getInstance().deserialize(json, new TypeReference<List<BuilderSchemaAttribute>>() {});
		
		System.out.println(CordaUtil.getInstance().serialize(attrs));
	}
}
