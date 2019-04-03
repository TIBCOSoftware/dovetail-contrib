/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package smartcontract.trigger.transaction;

import java.io.IOException;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;
import com.fasterxml.jackson.core.JsonParseException;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.JsonMappingException;
import com.fasterxml.jackson.databind.ObjectMapper;

import smartcontract.trigger.transaction.model.composer.HLCResource;

public class MetadataParser {
	public static Map<String, HLCResource> parse(String schema) throws JsonParseException, JsonMappingException, IOException {
		Map<String, HLCResource> hlcResources = new LinkedHashMap<String, HLCResource>();
		ObjectMapper mapper = new ObjectMapper();
		
		//String[0] = resource name
		//String[1] = json schema
		List<String[]> schemas = mapper.readValue(schema, new TypeReference<List<String[]>>() {});
		for (String[] s : schemas){
			String nm = s[0];
			String json = s[1];
			
			Schema metaSchema = mapper.readValue(json, Schema.class);
			HLCResource hlcResource = mapper.readValue(metaSchema.getDescription(), HLCResource.class );
			hlcResources.put(nm, hlcResource);
		};
		
		return hlcResources;
	}
	
	 @JsonIgnoreProperties(ignoreUnknown = true)
	 public static class Schema {
		//only care about meta data which is stored in description field
		private  String description;

		public String getDescription() {
			return description;
		}

		public void setDescription(String description) {
			this.description = description;
		}
	 }
	 
	 public static HLCResource parseSingleSchema(String schema) throws JsonParseException, JsonMappingException, IOException {
		ObjectMapper mapper = new ObjectMapper();
		Schema metaSchema = mapper.readValue(schema, Schema.class);
		HLCResource hlcResource = mapper.readValue(metaSchema.getDescription(), HLCResource.class );
			
		return hlcResource;
	}
	 
}
