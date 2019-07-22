package com.tibco.dovetail.corda.json;

import java.io.IOException;
import java.security.PublicKey;

import com.fasterxml.jackson.core.JsonGenerator;
import com.fasterxml.jackson.databind.SerializerProvider;
import com.fasterxml.jackson.databind.ser.std.StdSerializer;

import net.corda.core.crypto.Base58;

public class PublicKeySerializer extends StdSerializer<PublicKey>{

	/**
	 * 
	 */
	private static final long serialVersionUID = 1L;

	public PublicKeySerializer() {
		this(null);
	}
	public PublicKeySerializer(Class<PublicKey> t) {
		super(t);
	}

	@Override
	public void serialize(PublicKey value, JsonGenerator gen, SerializerProvider provider) throws IOException {
		gen.writeString(Base58.encode(value.getEncoded()));
		
	}

}
