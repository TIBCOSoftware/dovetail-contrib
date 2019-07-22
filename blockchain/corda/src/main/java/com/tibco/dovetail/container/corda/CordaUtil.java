/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package com.tibco.dovetail.container.corda;

import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.module.SimpleModule;
import com.fasterxml.jackson.module.kotlin.KotlinModule;
import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.corda.json.CashSerializer;
import com.tibco.dovetail.corda.json.LinearIdDeserializer;
import com.tibco.dovetail.corda.json.LinearIdSerializer;
import com.tibco.dovetail.corda.json.MoneyAmtDeserializer;
import com.tibco.dovetail.corda.json.MoneyAmtSerializer;
import com.tibco.dovetail.corda.json.PartyDeserializer;
import com.tibco.dovetail.corda.json.PartySerializer;
import com.tibco.dovetail.corda.json.PublicKeySerializer;
import com.tibco.dovetail.corda.json.StateAndRefSerializer;
import com.tibco.dovetail.corda.json.AbstractPartyDeserializer;
import com.tibco.dovetail.corda.json.AbstractPartySerializer;
import com.tibco.dovetail.corda.json.CashDeserializer;
import com.tibco.dovetail.core.runtime.util.JsonUtil;

import net.corda.core.contracts.Amount;
import net.corda.core.contracts.ContractsDSL;
import net.corda.core.contracts.StateAndRef;
import net.corda.core.contracts.UniqueIdentifier;
import net.corda.core.crypto.Base58;
import net.corda.finance.contracts.asset.Cash;

import org.bouncycastle.util.encoders.Hex;

import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.security.PublicKey;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;
import net.corda.core.identity.AbstractParty;
import net.corda.core.identity.AnonymousParty;
import net.corda.core.identity.CordaX500Name;
import net.corda.core.identity.Party;
import net.corda.core.node.ServiceHub;

public class CordaUtil {
	static ObjectMapper mapper;
	static ServiceHub serviceHub;
	
	static {
		mapper = new ObjectMapper();
		mapper.registerModule(new KotlinModule());
		
		SimpleModule module = new SimpleModule();
	
		module.addDeserializer(AbstractParty.class, new AbstractPartyDeserializer());
		module.addSerializer(AbstractParty.class, new AbstractPartySerializer());
		
		module.addDeserializer(Party.class, new PartyDeserializer());
		module.addSerializer(Party.class, new PartySerializer());

		module.addSerializer(Amount.class, new MoneyAmtSerializer("java.util.Currency"));
		module.addDeserializer(Amount.class, new MoneyAmtDeserializer());

		module.addSerializer(UniqueIdentifier.class, new LinearIdSerializer());
		module.addDeserializer(UniqueIdentifier.class, new LinearIdDeserializer());
		
		module.addSerializer(Cash.State.class, new CashSerializer());
		module.addDeserializer(Cash.State.class, new CashDeserializer());
		module.addSerializer(StateAndRef.class, new StateAndRefSerializer());
		module.addSerializer(PublicKey.class, new PublicKeySerializer());
		mapper.registerModule(module);
	}
	
    public static DocumentContext toJsonObject(Object state){
    		if(state instanceof DocumentContext)
    			return (DocumentContext) state;
    		else {
    			String json = serialize (state);

    			return JsonUtil.getJsonParser().parse(json);
    		}
    }

    public static String serialize(Object o )  {
    		try {
				return mapper.writeValueAsString(o);
			} catch (JsonProcessingException e) {
				throw new RuntimeException(e);
			}
    }
    
    @SuppressWarnings({ "rawtypes", "unchecked" })
	public static Object deserialize(String json, Class clazz)  {
    		try {	
				return mapper.readValue(json, clazz);
			} catch (Exception e) {
				throw new RuntimeException(e);
			}
    }
    
    @SuppressWarnings({ "rawtypes" })
   	public static Object deserialize(String json, TypeReference clazz)  {
       		try {	
   				return mapper.readValue(json, clazz);
   			} catch (Exception e) {
   				throw new RuntimeException(e);
   			}
       }
    
    
	public static void compare(List<DocumentContext> actual, List<DocumentContext> results) {
        ContractsDSL.requireThat(check -> {
        		String astring="";
        		String rstring ="";
        		if(actual.size() != results.size()) {
        			astring = actual.stream().map(v -> v.jsonString()).collect(Collectors.joining(","));
        			rstring = results.stream().map(v -> v.jsonString()).collect(Collectors.joining(","));
        		}
            check.using("expected inputs/outputs have same number as what is in LedgerTransaction:  txIn=" + astring + ", flow=" + rstring, actual.size() == results.size());

            return null;
        });

        List<Map<String, Object>> av = new ArrayList<>();
        List<Map<String, Object>> rv = new ArrayList<>();
        try {
	        ObjectMapper mapper = new ObjectMapper();
	        for(int i=0; i < actual.size(); i++){
	        		av.add(mapper.readValue(actual.get(i).jsonString(), Map.class));
	        		rv.add(mapper.readValue(results.get(i).jsonString(), Map.class));
	        		//av.add(actual.get(i).json());
	        		//rv.add(results.get(i).json());
	        }
        }catch(Exception e) {
        		throw new IllegalArgumentException(e);
        }
        
        ContractsDSL.requireThat(check -> {
            check.using("expected inputs/outputs have same values as what is in LedgerTransaction: txIn=" + CordaUtil.serialize(av) + ", flow=" + CordaUtil.serialize(rv),av.containsAll(rv));

            return null;
        });
    }

    public static String sha256Hash(String value){
        try {
            MessageDigest digest = MessageDigest.getInstance("SHA-256");
            byte[] hash = digest.digest(value.getBytes(StandardCharsets.UTF_8));
            return new String(Hex.encode(hash));
        }catch (Exception e){
            throw new IllegalArgumentException(e);
        }
    }
    
    public static String partyToString(AbstractParty p) {
 		return Base58.encode(p.getOwningKey().getEncoded());
    }
 
	public static AbstractParty partyFromString(String s) {
		 if(serviceHub == null)
			 throw new RuntimeException("serviceHub is not initalized");
		 
		PublicKey key = net.corda.core.crypto.Crypto.decodePublicKey(Base58.decode(s));

 		if(serviceHub.getNetworkMapCache().getNodesByLegalIdentityKey(key).isEmpty())
 			return new AnonymousParty(key);
 		else
 			return serviceHub.getIdentityService().partyFromKey(key);

	 }
	
	public static AbstractParty partyFromCommonName(String s) {
		 if(serviceHub == null)
			 throw new RuntimeException("serviceHub is not initalized");
		 
		 return serviceHub.getIdentityService().wellKnownPartyFromX500Name(CordaX500Name.parse(s));
	 }
	 
	 public static void setServiceHub(ServiceHub hub) {
		 serviceHub = hub;
	 }
	 
	 public static PublicKey decodeKey(String key) {
		 return net.corda.core.crypto.Crypto.decodePublicKey(Base58.decode(key));
	 }
	 
	 public static JsonNode toJsonNode(Object v) {
		 return mapper.valueToTree(v);
	 }
    
}
