package com.tibco.dovetail.container.corda;

import net.corda.core.contracts.CommandData;
import net.corda.core.serialization.CordaSerializable;

import java.util.LinkedHashMap;
import java.util.Map;

@CordaSerializable
public class CordaCommandDataWithData implements CommandData {
    Map<String, Object> data = new LinkedHashMap<String, Object>();
    public void putData(String param, Object value){
        data.put(param, value);
    }

    public Object getData(String param){
        return data.get(param);
    }
}
