package com.tibco.dovetail.container.cordapp.flows;

import org.antlr.v4.runtime.CharStream;

import com.tibco.dovetail.core.runtime.expression.MapExprGrammarLexer;

import net.corda.core.serialization.SerializationCustomSerializer;

public class MapExprGrammarLexerCustomSerializer implements SerializationCustomSerializer<MapExprGrammarLexer, MapExprGrammarLexerCustomSerializer.Proxy>{

	public static class Proxy {
		CharStream input;
		public Proxy(CharStream in) {
			input = in;
		}
		
		public CharStream getCharStream() {
			return this.input;
		}
	}

	@Override
	public MapExprGrammarLexer fromProxy(Proxy proxy) {
		return new MapExprGrammarLexer(proxy.input);
	}

	@Override
	public Proxy toProxy(MapExprGrammarLexer obj) {
		return new Proxy(obj._input);
	}
	
}
