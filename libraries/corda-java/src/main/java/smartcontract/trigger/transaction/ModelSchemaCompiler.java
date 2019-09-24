/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package smartcontract.trigger.transaction;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.tibco.dovetail.core.model.composer.HLCResource;

import java.io.IOException;
import java.io.InputStream;
import java.util.LinkedHashMap;
import java.util.Map;

public class ModelSchemaCompiler {
    public static Map<String, HLCResource> compile(InputStream model) throws IOException{
        Map<String, HLCResource> hlcResources = new LinkedHashMap<String, HLCResource>();

        ObjectMapper mapper = new ObjectMapper();
        hlcResources = mapper.readValue(model, new TypeReference<Map<String, HLCResource>>(){} ); 

        return hlcResources;
    }
}
