'use strict';

function createFlowInputSchema(inputParams) {
    if(Boolean(inputParams) == false)
        return "{}";

   let inputs = JSON.parse(inputParams);
   let schema = {schema:"http://json-schema.org/draft-04/schema#", type: "object", properties:{}, description:""}
   let metadata = {metadata: {type: "Transaction"}, attributes: []};

   if(inputs) {
       for(let i=0; i<inputs.length; i++){
            let name = inputs[i].parameterName;
            let tp = inputs[i].type;
            let repeating = inputs[i].repeating;
            let partyType = inputs[i].partyRole;

            let datatype = {type: tp.toLowerCase()};
            let javatype = tp;
            let isRef = false;
            let isArray = false;
            let attr = {};

            switch (tp) {
                case "Party":
                    datatype.type = "string";
                    javatype = "net.corda.core.identity.Party";
                    isRef = true;
                    break;
                case "LinearId":
                    datatype.type = "string";
                   // datatype.type = "object";
                  //  datatype["properties"] = {uuid: {type: "string"}, externalId: {type: "string"}};
                    javatype = "net.corda.core.contracts.UniqueIdentifier";
                    break;
                case "Amount<Currency>":
                    datatype.type = "object";
                    datatype["properties"] = {currency: {type: "string"}, quantity: {type: "number"}};
                    javatype = "net.corda.core.contracts.Amount<Currency>";
                    break;
                case "Integer":
                case "Long":
                    datatype.type = "number";
                    javatype = "java.lang.Long";
                    break;
                case "Decimal":
                    datatype.type = "string";
                    javatype = "java.math.BigDecimal";
                    break;
                case "LocalDate":
                    datatype.type = "string";
                    datatype["format"] = "date-time";
                    javatype = "java.time.LocalDate";
                    break;
                case "DateTime":
                    datatype.type = "string";
                    datatype["format"] = "date-time";
                    javatype = "java.time.Instant"
                    break;
            }
            if(repeating === "true"){
                schema.properties[name] = {type: "array", items: {datatype}}
                isArray = true;
            } else {
                schema.properties[name] = datatype
            }

            attr["name"] = name;
            attr["type"] = javatype;
            attr["isRef"] = isRef
            attr["isArray"] = isArray;
            attr["partyType"] = partyType;
        
            metadata.attributes.push(attr);
       }
       schema.description = JSON.stringify(metadata)
       return JSON.stringify(schema);
   } else {
       return "{}";
   }
}

module.exports = {"createFlowInputSchema": createFlowInputSchema};