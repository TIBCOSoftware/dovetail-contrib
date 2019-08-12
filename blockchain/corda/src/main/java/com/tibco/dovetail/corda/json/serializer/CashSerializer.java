package com.tibco.dovetail.corda.json.serializer;

import java.io.IOException;

import javax.xml.bind.DatatypeConverter;

import com.fasterxml.jackson.core.JsonGenerator;
import com.fasterxml.jackson.databind.SerializerProvider;
import com.fasterxml.jackson.databind.ser.std.StdSerializer;
import com.tibco.dovetail.container.corda.CordaUtil;

import net.corda.finance.contracts.asset.Cash.State;

public class CashSerializer extends StdSerializer<State> {

	/**
	 * 
	 */
	private static final long serialVersionUID = 1L;

	public CashSerializer() {
		this(null);
	}
	protected CashSerializer(Class<State> t) {
		super(t);
	}

	@Override
	public void serialize(State value, JsonGenerator gen, SerializerProvider provider) throws IOException {
		
		gen.writeStartObject();
		gen.writeStringField("owner",CordaUtil.getInstance().partyToString(value.getOwner()));
		gen.writeStringField("issuer", CordaUtil.getInstance().partyToString(value.getAmount().getToken().getIssuer().getParty()));
		gen.writeStringField("issuerRef", DatatypeConverter.printBase64Binary(value.getAmount().getToken().getIssuer().getReference().getBytes()));
		gen.writeObjectFieldStart("amt");
		gen.writeNumberField("quantity", value.getAmount().getQuantity());
		gen.writeStringField("currency", value.getAmount().getToken().getProduct().getCurrencyCode());
		gen.writeEndObject();
		gen.writeEndObject();
	}

}
