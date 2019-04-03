package com.tibco.dovetail.corda.json;

import java.io.IOException;
import java.util.Currency;

import com.fasterxml.jackson.core.JsonParser;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.DeserializationContext;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.deser.std.StdDeserializer;
import com.fasterxml.jackson.databind.node.IntNode;
import com.fasterxml.jackson.databind.node.LongNode;

import net.corda.core.contracts.Amount;

public class MoneyAmtDeserializer extends StdDeserializer<Amount<?>>{

	/**
	 * 
	 */
	private static final long serialVersionUID = 1L;

	public MoneyAmtDeserializer() {
		this(null);
	}
	protected MoneyAmtDeserializer(Class<?> vc) {
		super(vc);
	}

	@Override
	public Amount<?> deserialize(JsonParser p, DeserializationContext ctxt)
			throws IOException, JsonProcessingException {
		
		JsonNode node = p.getCodec().readTree(p);
		JsonNode qnode = node.get("quantity"); 
		long value = 0;
		if (qnode instanceof LongNode)
			value = (long)((LongNode) qnode).numberValue();
		else
			value = (int)((IntNode) qnode).numberValue();
		
        String currency = node.get("currency").asText();
       
		return new Amount<Currency>(value, Currency.getInstance(currency));
	}

}
