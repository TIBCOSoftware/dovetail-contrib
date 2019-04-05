package com.tibco.dovetail.corda.json;

import java.io.IOException;
import java.util.UUID;

import com.fasterxml.jackson.core.JsonParser;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.DeserializationContext;
import com.fasterxml.jackson.databind.deser.std.StdDeserializer;

import net.corda.core.contracts.UniqueIdentifier;

public class LinearIdDeserializer extends StdDeserializer<UniqueIdentifier> {

	/**
	 * 
	 */
	private static final long serialVersionUID = 1L;

	public LinearIdDeserializer() {
		this(UniqueIdentifier.class);
	}
	protected LinearIdDeserializer(Class<?> vc) {
		super(vc);
		// TODO Auto-generated constructor stub
	}

	@Override
	public UniqueIdentifier deserialize(JsonParser arg0, DeserializationContext arg1)
			throws IOException, JsonProcessingException {
		
		String id = arg0.getValueAsString();
		return LinearIdDeserializer.fromString(id);
	}
	
	public static UniqueIdentifier fromString(String id) {
		String[] ids = id.split("#");
		if(ids.length < 2)
			return new UniqueIdentifier(id, UUID.randomUUID());
		else
			return new UniqueIdentifier(ids[0], UUID.fromString(ids[1]));
	}

}
