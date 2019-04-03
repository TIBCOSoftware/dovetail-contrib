package com.tibco.dovetail.corda.json;

import java.io.IOException;
import java.util.Currency;
import java.util.UUID;

import org.junit.Assert;
import org.junit.Test;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.module.SimpleModule;
import com.tibco.dovetail.container.corda.CordaUtil;
import com.tibco.dovetail.container.cordapp.AppContainer;

import net.corda.core.contracts.Amount;
import net.corda.core.contracts.UniqueIdentifier;
import net.corda.core.identity.AbstractParty;
import net.corda.core.identity.CordaX500Name;
import net.corda.core.identity.Party;
import net.corda.finance.contracts.asset.Cash;
import net.corda.testing.core.TestIdentity;
import net.corda.testing.node.MockServices;

public class TestSerializer {
	static TestIdentity bob;
	static MockServices mock;
	static AppContainer ctnr;
	static {
		
		SimpleModule module = new SimpleModule();
		
		bob = new TestIdentity(new CordaX500Name("BigCorp", "New York", "GB"));
		mock = new MockServices(bob);
		//CordaUtil.initWithCordaRuntime(mock);
		ctnr = new AppContainer(mock);
		
}
	@Test
	public void testAmount() {
		System.out.println("testAmount....");
		Amount<Currency> amt = new Amount<Currency>(100L, Currency.getInstance("USD"));
	
		try {
			System.out.println(amt.toString());
			
			String json = CordaUtil.serialize(amt);
			System.out.println(json); 
			
			Amount amt2 = (Amount) CordaUtil.deserialize( json, Amount.class);
			
			System.out.println(amt2.toString());
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
		
	}
	
	@Test
	public void testCash() {
		System.out.println("testCash....");
		
		Amount<Currency> amt = new Amount<Currency>(100L, Currency.getInstance("USD"));
		
		Cash.State cash = new Cash.State(bob.ref("0".getBytes()), amt, bob.getParty());
		
		try {
			
			String json = CordaUtil.serialize(cash);
			System.out.println(json); 
			
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
		
	}
	
	@Test
	public void testParty() {
		System.out.println("testParty....");
		Party pbob = bob.getParty();
		
		try {
			System.out.println("bob=" + pbob.toString());
			
			String json = CordaUtil.serialize( pbob);
			System.out.println("json=" + json); 
			
			Party amt2 = (Party) CordaUtil.deserialize(json, AbstractParty.class);
			
			System.out.println("bob2="+ amt2.toString());
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
		
	}
	
	@Test
	public void testPartyAmt() {
		System.out.println("testPartyAmt....");
			
		try {
			
			String json = CordaUtil.serialize(new TestSerializer.PartyAmount());
			System.out.println("json=" + json); 
			
			TestSerializer.PartyAmount amt2 = (TestSerializer.PartyAmount) CordaUtil.deserialize(json, TestSerializer.PartyAmount.class);
			
			System.out.println("bob2="+ CordaUtil.serialize( amt2));
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
		
	}
	
	@Test
	public void testUniqueIdentifier() {
		System.out.println("testUniqueIdentifier....");
			
		try {
			UniqueIdentifier id1 = new UniqueIdentifier("test", UUID.fromString("bb09aeb4-c053-4295-8718-964f348a4ebf"));
			String json = CordaUtil.serialize(id1);
			System.out.println("json=" + json); 
			
			UniqueIdentifier id2 = (UniqueIdentifier) CordaUtil.deserialize(json, UniqueIdentifier.class);
			String json2 = CordaUtil.serialize( id2);
			System.out.println("id2="+ json2);
			
			Assert.assertEquals(json, json2);
			
			
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
		
	}
	static class PartyAmount {
		AbstractParty p = bob.getParty();
		public AbstractParty getP() {
			return p;
		}
		public void setP(AbstractParty p) {
			this.p = p;
		}
		public Amount<Currency> getAmt() {
			return amt;
		}
		public void setAmt(Amount<Currency> amt) {
			this.amt = amt;
		}
		public String getName() {
			return name;
		}
		public void setName(String name) {
			this.name = name;
		}
		Amount<Currency> amt = new Amount<Currency>(100L, Currency.getInstance("USD"));
		String name = "test";
	}

}
