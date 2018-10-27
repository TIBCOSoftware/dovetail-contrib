package smartcontract.trigger.transaction.model.composer;

import com.fasterxml.jackson.annotation.JsonIgnoreProperties;

import java.util.HashMap;
import java.util.List;
import java.util.Map;

@JsonIgnoreProperties(ignoreUnknown = true)
public class HLCResource {
    private HLCMetadata metadata;
    private Map<String, HLCAttribute> attributes = new HashMap<String, HLCAttribute>();
    private List<HLCAttribute> attributeList;

    public HLCMetadata getMetadata() {
        return metadata;
    }

    public void setMetadata(HLCMetadata metadata) {
        this.metadata = metadata;
    }

    public HLCAttribute getAttribute(String attr) {
        return attributes.get(attr);
    }

    public void setAttributes(List<HLCAttribute> attributes) {
    		attributeList = attributes;
        attributes.forEach(attr -> this.attributes.put(attr.getName(), attr));
    }
    public List<HLCAttribute> getAttributes() {
       return attributeList;
    }
}
