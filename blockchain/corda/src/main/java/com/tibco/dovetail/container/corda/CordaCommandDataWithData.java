/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package com.tibco.dovetail.container.corda;

import net.corda.core.contracts.CommandData;
import net.corda.core.serialization.CordaSerializable;

import java.util.LinkedHashMap;
import java.util.Map;

import com.fasterxml.jackson.core.type.TypeReference;
import com.tibco.dovetail.container.corda.CordaUtil;

@CordaSerializable
public class CordaCommandDataWithData implements CommandData {
	
    transient Map<String, Object> data = new LinkedHashMap<String, Object>();
    
    //to work around R3 issue
    String serializedData;
    
    public String getSerializedData() {
		return serializedData;
	}

	public void setSerializedData(String serializedData) {
		this.serializedData = serializedData;
	}

	public void putData(String param, Object value){
    		data.put(param, value);
    }

    public Object getData(String param){
    		return data.get(param);
    }
    
    public void serialize() {
    		this.serializedData = CordaUtil.getInstance().serialize(data);
    }
    
    public void deserialize() {
    		data = (Map<String, Object>) CordaUtil.getInstance().deserialize(serializedData, new TypeReference<Map<String, Object>>(){});
    }
    
    public String getCommand() {
    		return (String) data.get("command");
    }
}
