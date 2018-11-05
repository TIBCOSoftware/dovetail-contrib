package com.tibco.dovetail.container.corda;

import com.fasterxml.jackson.core.JsonParseException;
import com.fasterxml.jackson.databind.JsonMappingException;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.jayway.jsonpath.DocumentContext;
import com.tibco.dovetail.core.runtime.util.JsonUtil;
import net.corda.core.contracts.ContractState;
import net.corda.core.contracts.ContractsDSL;
import net.corda.finance.contracts.asset.Cash;

import org.bouncycastle.util.encoders.Hex;

import java.io.IOException;
import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.stream.Collectors;
import net.corda.core.identity.AbstractParty;
import net.corda.core.identity.AnonymousParty;
import java.util.Currency;

public class CordaUtil {

    public static DocumentContext toJsonObject(Object state){
        String json = toJsonString(state);

        return JsonUtil.getJsonParser().parse(json);
    }

    public static DocumentContext toJsonObject(List<Object> states){
        String json = "[" + states.stream().map(s -> toJsonString(s)).collect(Collectors.joining(",")) + "]";

        return JsonUtil.getJsonParser().parse(json);
    }

    public static String toJsonString(Object state){
        String json = null;
        if(state instanceof Cash.State){
            Cash.State cash = (Cash.State)state;
            String party = cash.getOwner().toString();
           // String party = Base58.encode(cash.getOwner().getOwningKey().getEncoded());
            json = "{\"owner\":\"" + party +
                    "\", \"amt\":{\"quantity\":" + cash.getAmount().getQuantity() +
                    ",\"currency\":\"" + cash.getAmount().getToken().getProduct().getCurrencyCode() +"\"}}";

        } else if (state instanceof net.corda.core.contracts.Amount<?>) {
        		net.corda.core.contracts.Amount<Currency> amt = (net.corda.core.contracts.Amount<Currency>)state;
        		json = "{\"quantity\":" + amt.getQuantity() +
                        ",\"currency\":\"" + amt.getToken().getCurrencyCode() +"\"}";
        } else {
            json = state.toString();
        }

        return json;
    }

    public static String toString(Object obj){
        String string = null;
        if(obj instanceof AnonymousParty) {
        		return ((AnonymousParty) obj).toString();
        } else if(obj instanceof AbstractParty){
            AbstractParty party = (AbstractParty)obj;
            string = party.toString();
        } else {
            string = obj.toString();
        }

        return string;
    }
    @SuppressWarnings("unchecked")
	public static void compare(List<DocumentContext> actual, List<DocumentContext> results) throws JsonParseException, JsonMappingException, IOException{
        ContractsDSL.requireThat(check -> {
        		String astring="";
        		String rstring ="";
        		if(actual.size() != results.size()) {
        			astring = actual.stream().map(v -> v.jsonString()).collect(Collectors.joining(","));
        			rstring = results.stream().map(v -> v.jsonString()).collect(Collectors.joining(","));
        		}
            check.using("expected outputs have same number as what is in LedgerTransaction:  txIn=" + astring + ", flowOutput=" + rstring, actual.size() == results.size());

            return null;
        });

        List<Map<String, Object>> av = new ArrayList<>();
        List<Map<String, Object>> rv = new ArrayList<>();

        ObjectMapper mapper = new ObjectMapper();
        for(int i=0; i < actual.size(); i++){
        		av.add(mapper.readValue(actual.get(i).jsonString(), Map.class));
        		rv.add(mapper.readValue(results.get(i).jsonString(), Map.class));
        }
        
        ContractsDSL.requireThat(check -> {
            check.using("expected outputs have same values as what is in LedgerTransaction: txIn=" + av.toString() + ", flowOutput=" + rv.toString(),av.containsAll(rv));

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
}
