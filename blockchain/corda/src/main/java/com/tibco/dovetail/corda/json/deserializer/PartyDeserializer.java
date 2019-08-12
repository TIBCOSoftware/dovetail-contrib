package com.tibco.dovetail.corda.json.deserializer;

import java.io.IOException;

import com.fasterxml.jackson.core.JsonParser;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.DeserializationContext;
import com.fasterxml.jackson.databind.deser.std.StdDeserializer;
import com.tibco.dovetail.container.cordapp.AppUtil;

import net.corda.core.identity.AbstractParty;
import net.corda.core.identity.Party;

public class PartyDeserializer extends StdDeserializer<Party>{

	/**
	 * 
	 */
	private static final long serialVersionUID = 1L;

	
	public PartyDeserializer()  {
		this(AbstractParty.class);
	}
	
	protected PartyDeserializer(Class<?> vc) {
		super(vc);
	}

	@Override
	public Party deserialize(JsonParser p, DeserializationContext ctxt)
			throws IOException, JsonProcessingException {
		
			return (Party) AppUtil.partyFromString(p.getText());
	}

}
