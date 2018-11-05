package smartcontract.trigger.transaction;

import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.ObjectMapper;

import smartcontract.trigger.transaction.model.composer.HLCResource;

import java.io.IOException;
import java.io.InputStream;
import java.util.HashMap;
import java.util.Map;

public class ModelSchemaCompiler {
    public static Map<String, HLCResource> compile(InputStream model) throws IOException{
        Map<String, HLCResource> hlcResources = new HashMap<String, HLCResource>();

        ObjectMapper mapper = new ObjectMapper();
        hlcResources = mapper.readValue(model, new TypeReference<Map<String, HLCResource>>(){} ); 

        return hlcResources;
    }
}
