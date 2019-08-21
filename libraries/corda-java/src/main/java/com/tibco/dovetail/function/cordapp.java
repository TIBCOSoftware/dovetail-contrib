package com.tibco.dovetail.function;

import java.util.UUID;

public class cordapp {
	public static String createLinearId() {
		return "null#" + UUID.randomUUID().toString();
	}
	
	public static String createLinearIdFromExternalId(String extId) {
		return extId + "#" + UUID.randomUUID().toString();
	}
	
	public static String toLinearId(String extId, String uuid) {
		return extId + "#" + uuid;
	}
}
