package com.tibco.dovetail.corda.json;

import java.io.IOException;

import com.fasterxml.jackson.core.JsonGenerator;
import com.fasterxml.jackson.databind.SerializerProvider;
import com.fasterxml.jackson.databind.ser.std.StdSerializer;

import net.corda.core.contracts.StateAndRef;

@SuppressWarnings("rawtypes")
public class StateAndRefSerializer extends StdSerializer<StateAndRef> {

	/**
	 * 
	 */
	private static final long serialVersionUID = 1L;
	
	public StateAndRefSerializer() {
		this(StateAndRef.class);
	}
	protected StateAndRefSerializer(Class<StateAndRef> t) {
		super(t);
		// TODO Auto-generated constructor stub
	}

	@Override
	public void serialize(StateAndRef arg0, JsonGenerator gen, SerializerProvider arg2) throws IOException {
		
		String ref = getRef(arg0);
		
		gen.writeStartObject();
		gen.writeObjectField("data", arg0.getState().getData());
		gen.writeStringField("ref",ref);
		gen.writeEndObject();
	}
	
	public static String getRef(StateAndRef ref) {
		return ref.getRef().getTxhash().toString() + "###" + ref.getRef().getIndex();
	}

}
