/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package smartcontract.trigger.transaction.model.composer;

import java.util.Collection;
import java.util.LinkedHashMap;
import java.util.List;
import java.util.Map;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

@JsonIgnoreProperties(ignoreUnknown = true)
public class HLCMetadata {
    private String name;
    private ResourceType type;
    private String identifiedBy = null;
    private String parent = null;
    private String cordaClass = null;
    private Map<String, HLCDecorator> decorators = new LinkedHashMap<String, HLCDecorator>();

    public String getName() {
        return name;
    }
    public void setName(String name) {
        this.name = name;
    }

    public ResourceType getType() {
        return type;
    }
    public void setType(String type) {
        this.type = ResourceType.valueOf(type);
    }

    public String getIdentifiedBy() {
        return identifiedBy;
    }
    public void setIdentifiedBy(String identifiedBy) {
        this.identifiedBy = identifiedBy;
    }

    public String getParent() {
        return parent;
    }
    public void setParent(String parent) {
        this.parent = parent;
    }

    public String getCordaClass() {
        return cordaClass;
    }
    public void setCordaClass(String cordaClass) {
        this.cordaClass = cordaClass;
    }

    public void setDecorators(List<HLCDecorator> inputs) {
    	inputs.forEach(in -> this.decorators.put(in.getName(), in));
    }
    
    public Collection<HLCDecorator> getDecorators() {
       return decorators.values();
    }

    public HLCDecorator getDecorator(String name) {
    		return this.decorators.get(name);
    }
    
	public static enum ResourceType {
        Asset, Participant, Transaction, Event, Concept
    }
}
