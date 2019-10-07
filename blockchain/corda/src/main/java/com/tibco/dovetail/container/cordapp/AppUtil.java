package com.tibco.dovetail.container.cordapp;

import java.security.PublicKey;

import com.fasterxml.jackson.core.JsonGenerator;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import com.fasterxml.jackson.databind.module.SimpleModule;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.fasterxml.jackson.module.kotlin.KotlinModule;
import com.tibco.dovetail.corda.json.deserializer.AbstractPartyDeserializer;
import com.tibco.dovetail.corda.json.deserializer.CashDeserializer;
import com.tibco.dovetail.corda.json.deserializer.LinearIdDeserializer;
import com.tibco.dovetail.corda.json.deserializer.MoneyAmtDeserializer;
import com.tibco.dovetail.corda.json.deserializer.PartyDeserializer;

import net.corda.core.contracts.Amount;
import net.corda.core.contracts.UniqueIdentifier;
import net.corda.core.crypto.Base58;
import net.corda.core.identity.AbstractParty;
import net.corda.core.identity.AnonymousParty;
import net.corda.core.identity.CordaX500Name;
import net.corda.core.identity.Party;
import net.corda.core.node.ServiceHub;
import net.corda.finance.contracts.asset.Cash;

public class AppUtil {
	static ServiceHub serviceHub;
	
	private static ObjectMapper der_mapper = null;

	static {
		try {
				der_mapper = new ObjectMapper();
				der_mapper.enable(DeserializationFeature.USE_BIG_DECIMAL_FOR_FLOATS);
				der_mapper.enable(DeserializationFeature.USE_LONG_FOR_INTS);
				der_mapper.enable(JsonGenerator.Feature.WRITE_BIGDECIMAL_AS_PLAIN);
				der_mapper.disable(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS);
				der_mapper.registerModule(new KotlinModule());
				
				SimpleModule module = new SimpleModule();
			
				module.addDeserializer(AbstractParty.class, new AbstractPartyDeserializer());	
				module.addDeserializer(Party.class, new PartyDeserializer());
				module.addDeserializer(Amount.class, new MoneyAmtDeserializer());
				module.addDeserializer(UniqueIdentifier.class, new LinearIdDeserializer());
				module.addDeserializer(Cash.State.class, new CashDeserializer());

				der_mapper.registerModule(module);
				der_mapper.registerModule(new JavaTimeModule());
			
		}catch(Exception e) {
			System.out.println(e);
			throw new RuntimeException (e);
		}
	}
	
	 public static AbstractParty partyFromString(String s) {
		 if(serviceHub == null)
			 throw new RuntimeException("serviceHub is not initalized");
		 
		PublicKey key = net.corda.core.crypto.Crypto.decodePublicKey(Base58.decode(s));
	    try {
			if(serviceHub.getNetworkMapCache().getNodesByLegalIdentityKey(key).isEmpty())
				return new AnonymousParty(key);
			else
				return serviceHub.getIdentityService().partyFromKey(key);
	    }catch(java.lang.UnsupportedOperationException e) {
	    		return serviceHub.getIdentityService().partyFromKey(key);
	    }
	
	 }
	
	public static AbstractParty partyFromCommonName(String s) {
		 if(serviceHub == null)
			 throw new RuntimeException("serviceHub is not initalized");
		 
		 return serviceHub.getIdentityService().wellKnownPartyFromX500Name(CordaX500Name.parse(s));
	 }
	 
	 public static void setServiceHub(ServiceHub hub) {
		 serviceHub = hub;
	 }
	 
	 @SuppressWarnings({ "rawtypes", "unchecked" })
		public static Object deserialize(String json, Class clazz)  {
	    		try {	
					return der_mapper.readValue(json, clazz);
				} catch (Exception e) {
					throw new RuntimeException(e);
				}
	    }
}
