package com.tibco.dovetail.corda.json;

import java.io.IOException;

import com.fasterxml.jackson.core.JsonGenerator;
import com.fasterxml.jackson.databind.SerializerProvider;
import com.fasterxml.jackson.databind.ser.std.StdSerializer;
import net.corda.core.contracts.UniqueIdentifier;

public class LinearIdSerializer extends StdSerializer<UniqueIdentifier> {
	
    /**
	 * 
	 */
	private static final long serialVersionUID = 1L;

	public LinearIdSerializer() {
    		this(UniqueIdentifier.class);
    }
	protected LinearIdSerializer(Class<UniqueIdentifier> t) {
		super(t);
		// TODO Auto-generated constructor stub
	}

	@Override
	public void serialize(UniqueIdentifier arg0, JsonGenerator arg1, SerializerProvider arg2) throws IOException {
		arg1.writeString(arg0.getExternalId() + "#" + arg0.getId().toString());
		
	}

	public static String toString(UniqueIdentifier arg0) {
		return arg0.getExternalId() + "#" + arg0.getId().toString();
	}
}
