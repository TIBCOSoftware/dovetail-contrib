/*
* Copyright © 2018. TIBCO Software Inc.
* This file is subject to the license terms contained
* in the license file that is distributed with this file.
 */
package smartcontract.trigger.transaction.model.composer;

import com.tibco.dovetail.core.model.common.SimpleAttribute;

public class HLCAttribute extends SimpleAttribute {

    boolean required = true;
    boolean isRef = false;
    boolean isArray = false;
    String minValue = null;
    String maxValue = null;
    String pattern = null;

    public void setRequired(boolean required) {
        this.required = required;
    }

    public boolean isRef() {
        return isRef;
    }

    public void setIsRef(boolean ref) {
        isRef = ref;
    }

    public boolean isArray() {
        return isArray;
    }

    public void setIsArray(boolean array) {
        isArray = array;
    }

    public String getMinValue() {
        return minValue;
    }

    public void setMinValue(String minValue) {
        this.minValue = minValue;
    }

    public String getMaxValue() {
        return maxValue;
    }

    public void setMaxValue(String maxValue) {
        this.maxValue = maxValue;
    }

    public String getPattern() {
        return pattern;
    }

    public void setPattern(String pattern) {
        this.pattern = pattern;
    }

    public boolean isRequired() {
        return required;
    }


}
