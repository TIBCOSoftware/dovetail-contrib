package com.tibco.dovetail.corda.json;

import java.io.IOException;

import com.fasterxml.jackson.core.JsonGenerator;
import com.fasterxml.jackson.databind.SerializerProvider;
import com.fasterxml.jackson.databind.ser.std.StdSerializer;
import com.tibco.dovetail.container.cordapp.AppContainer;

import net.corda.core.identity.Party;

public class PartySerializer extends StdSerializer<Party>{

	/**
	 * 
	 */
	private static final long serialVersionUID = 1L;

	public PartySerializer() {
		this(null);
	}
	public PartySerializer(Class<Party> t) {
		super(t);
	}

	@Override
	public void serialize(Party value, JsonGenerator gen, SerializerProvider provider) throws IOException {
		gen.writeString(AppContainer.partyToString(value));
		
	}

}
