package com.tibco.dovetail.corda.json;

import java.io.IOException;
import java.util.Currency;

import javax.xml.bind.DatatypeConverter;

import com.fasterxml.jackson.core.JsonParser;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.databind.DeserializationContext;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.deser.std.StdDeserializer;
import com.tibco.dovetail.container.corda.CordaUtil;

import net.corda.core.contracts.Amount;
import net.corda.core.contracts.PartyAndReference;
import net.corda.core.utilities.OpaqueBytes;
import net.corda.finance.contracts.asset.Cash;
import net.corda.finance.contracts.asset.Cash.State;

public class CashDeserializer extends StdDeserializer<State>{

	/**
	 * 
	 */
	private static final long serialVersionUID = 1L;

	public CashDeserializer() {
		super(Cash.State.class);
	}
	
	protected CashDeserializer(Class<State> t) {
		super(t);
	}

	@Override
	public State deserialize(JsonParser p, DeserializationContext arg1) throws IOException, JsonProcessingException {
		JsonNode node = p.getCodec().readTree(p);
		
		JsonNode qnode = node.get("issuer"); 
		String issuer = qnode.textValue();
		
		qnode = node.get("issuerRef"); 
		String issuerRef = qnode.textValue();

		qnode = node.get("owner"); 
		String owner = qnode.textValue();
		
		Amount<Currency> amt = MoneyAmtDeserializer.parseAmount(node.get("amt"));
		
		return new Cash.State(new PartyAndReference(CordaUtil.partyFromString(issuer), OpaqueBytes.of(DatatypeConverter.parseBase64Binary(issuerRef))), amt, CordaUtil.partyFromString(owner));
	}

}
