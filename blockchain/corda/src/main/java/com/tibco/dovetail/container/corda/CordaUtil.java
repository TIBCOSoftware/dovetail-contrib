/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package com.tibco.dovetail.container.corda;

import com.fasterxml.jackson.core.JsonGenerator;
import com.fasterxml.jackson.core.JsonProcessingException;
import com.fasterxml.jackson.core.type.TypeReference;
import com.fasterxml.jackson.databind.DeserializationFeature;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.SerializationFeature;
import com.fasterxml.jackson.databind.module.SimpleModule;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import com.fasterxml.jackson.module.kotlin.KotlinModule;
import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.corda.json.serializer.AbstractPartySerializer;
import com.tibco.dovetail.corda.json.serializer.CashSerializer;
import com.tibco.dovetail.corda.json.serializer.LinearIdSerializer;
import com.tibco.dovetail.corda.json.serializer.MoneyAmtSerializer;
import com.tibco.dovetail.corda.json.serializer.PartySerializer;
import com.tibco.dovetail.corda.json.serializer.PublicKeySerializer;
import com.tibco.dovetail.corda.json.serializer.StateAndRefSerializer;
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
import net.corda.core.identity.Party;

public class CordaUtil {
	private ObjectMapper ser_mapper = null;
	private static CordaUtil cordaUtil = null;
	
	private CordaUtil() {
		ser_mapper = new ObjectMapper();
		ser_mapper.enable(JsonGenerator.Feature.WRITE_BIGDECIMAL_AS_PLAIN);
		//ser_mapper.configure(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS, false);
		ser_mapper.disable(SerializationFeature.WRITE_DATES_AS_TIMESTAMPS);
		ser_mapper.enable(DeserializationFeature.USE_BIG_DECIMAL_FOR_FLOATS);
		ser_mapper.enable(DeserializationFeature.USE_LONG_FOR_INTS);
	
		ser_mapper.registerModule(new KotlinModule());
		
		SimpleModule module = new SimpleModule();
	
		module.addSerializer(AbstractParty.class, new AbstractPartySerializer());
		module.addSerializer(Party.class, new PartySerializer());
		module.addSerializer(Amount.class, new MoneyAmtSerializer("java.util.Currency"));
		module.addSerializer(UniqueIdentifier.class, new LinearIdSerializer());			
		module.addSerializer(Cash.State.class, new CashSerializer());
		module.addSerializer(StateAndRef.class, new StateAndRefSerializer());
		module.addSerializer(PublicKey.class, new PublicKeySerializer());

		ser_mapper.registerModule(module);
		ser_mapper.registerModule(new JavaTimeModule());
	}
	public static synchronized CordaUtil getInstance() {
		if(cordaUtil == null) {
			cordaUtil = new CordaUtil();
		}
		return cordaUtil;
	}
    public  DocumentContext toJsonObject(Object state){
    		if(state instanceof DocumentContext)
    			return (DocumentContext) state;
    		else {
    			String json = serialize (state);

    			return JsonUtil.getJsonParser().parse(json);
    		}
    }

    public  String serialize(Object o )  {
    		try {
				return ser_mapper.writeValueAsString(o);
			} catch (JsonProcessingException e) {
				throw new RuntimeException(e);
			}
    }
    
    
	public void compare(List<DocumentContext> actual, List<DocumentContext> results) {
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
	        }
        }catch(Exception e) {
        		throw new IllegalArgumentException(e);
        }
        
        ContractsDSL.requireThat(check -> {
            check.using("expected inputs/outputs have same values as what is in LedgerTransaction: txIn=" + CordaUtil.getInstance().serialize(av) + ", flow=" + CordaUtil.getInstance().serialize(rv),av.containsAll(rv));

            return null;
        });
    }

    public String sha256Hash(String value){
        try {
            MessageDigest digest = MessageDigest.getInstance("SHA-256");
            byte[] hash = digest.digest(value.getBytes(StandardCharsets.UTF_8));
            return new String(Hex.encode(hash));
        }catch (Exception e){
            throw new IllegalArgumentException(e);
        }
    }
    
    public String partyToString(AbstractParty p) {
 		return Base58.encode(p.getOwningKey().getEncoded());
    }
	 
	 public static PublicKey decodeKey(String key) {
		 return net.corda.core.crypto.Crypto.decodePublicKey(Base58.decode(key));
	 }
	 
	 public JsonNode toJsonNode(Object v) {
		 return ser_mapper.valueToTree(v);
	 }
	 
	
	    
    @SuppressWarnings({ "rawtypes" })
   	public  Object deserialize(String json, TypeReference clazz)  {
       		try {	
   				return ser_mapper.readValue(json, clazz);
   			} catch (Exception e) {
   				throw new RuntimeException(e);
   			}
       }
	
}
