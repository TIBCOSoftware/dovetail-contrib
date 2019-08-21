package com.tibco.dovetail.corda.json.serializer;

import java.io.IOException;
import java.util.Currency;

import com.fasterxml.jackson.core.JsonGenerator;
import com.fasterxml.jackson.databind.SerializerProvider;
import com.fasterxml.jackson.databind.ser.std.StdSerializer;

import net.corda.core.contracts.Amount;

@SuppressWarnings("rawtypes")
public class MoneyAmtSerializer extends StdSerializer<Amount>{

	/**
	 * 
	 */
	private static final long serialVersionUID = 1L;
	private String amtType;

	public MoneyAmtSerializer(String javatype) {
		this(Amount.class);
		this.amtType = javatype;
	}
	
	protected MoneyAmtSerializer(Class<Amount> t) {
		super(t);
	}

	@Override
	public void serialize(Amount value, JsonGenerator gen, SerializerProvider provider) throws IOException {
		gen.writeStartObject();
		if (amtType.equals("java.util.Currency")) {
			Currency c = (Currency) value.getToken();
			gen.writeStringField("currency", c.getCurrencyCode());
		}
		gen.writeNumberField("quantity", value.getQuantity());
		gen.writeEndObject();
		
	}
/*
	@Override
	public void serialize(Amount<?> value, JsonGenerator gen, SerializerProvider provider) throws IOException {
		if (value.getToken() instanceof Currency) {
			gen.writeStartObject();
			Currency c = (Currency) value.getToken();
			gen.writeStringField("currency", c.getCurrencyCode());
			gen.writeNumberField("quantity", value.getQuantity());
			gen.writeEndObject();
		}
		
	}
*/
}
