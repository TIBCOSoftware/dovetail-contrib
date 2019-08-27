package com.tibco.dovetail.corda.json.deserializer;

import java.io.IOException;

import com.fasterxml.jackson.core.JsonParser;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.DeserializationContext;
import com.fasterxml.jackson.databind.deser.std.StdDeserializer;
import com.tibco.dovetail.container.cordapp.AppUtil;

import net.corda.core.identity.AbstractParty;

public class AbstractPartyDeserializer extends StdDeserializer<AbstractParty>{

	/**
	 * 
	 */
	private static final long serialVersionUID = 1L;

	
	public AbstractPartyDeserializer()  {
		this(AbstractParty.class);
	}
	
	protected AbstractPartyDeserializer(Class<?> vc) {
		super(vc);
	}

	@Override
	public AbstractParty deserialize(JsonParser p, DeserializationContext ctxt)
			throws IOException, JsonProcessingException {
		
			return AppUtil.partyFromString(p.getText());
	}

}
