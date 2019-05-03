package com.tibco.dovetail.corda.json;

import java.io.IOException;

import com.fasterxml.jackson.core.JsonGenerator;
import com.fasterxml.jackson.databind.SerializerProvider;
import com.fasterxml.jackson.databind.ser.std.StdSerializer;
import com.tibco.dovetail.container.corda.CordaUtil;

import net.corda.core.identity.AbstractParty;

public class AbstractPartySerializer extends StdSerializer<AbstractParty>{

	/**
	 * 
	 */
	private static final long serialVersionUID = 1L;

	public AbstractPartySerializer() {
		this(null);
	}
	public AbstractPartySerializer(Class<AbstractParty> t) {
		super(t);
	}

	@Override
	public void serialize(AbstractParty value, JsonGenerator gen, SerializerProvider provider) throws IOException {
		gen.writeString(CordaUtil.partyToString(value));
		
	}

}
