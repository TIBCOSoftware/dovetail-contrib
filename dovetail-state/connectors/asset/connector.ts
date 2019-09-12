import { Inject, Injectable, Injector } from "@angular/core";
import { WiContrib, WiServiceHandlerContribution, WiContributionUtils,AUTHENTICATION_TYPE } from "wi-studio/app/contrib/wi-contrib";
import { IConnectorContribution } from "wi-studio/common/models/contrib";
import { IValidationResult, ValidationResult, ValidationError } from "wi-studio/common/models/validation";
import { IActionResult, ActionResult } from "wi-studio/common/models/contrib";
import { Observable } from "rxjs/Observable";

@Injectable()
@WiContrib({})
export class AssetConnectorService extends WiServiceHandlerContribution {
    constructor( @Inject(Injector) injector) {
        super(injector);
    }

    value = (fieldName: string, context: IConnectorContribution): Observable<any> | any => {
        let assettype = context.getField("assettype").value;
        switch(fieldName){
            case "inputLinear":
                if(assettype === "Linear State"){
                    return "[{\"name\":\"linearId\",\"type\":\"LinearId\", \"partyType\":\"N/A\", \"repeating\":\"False\"}]"
                } 
                return null
            case "assetfields":
                if(assettype === "Linear State"){
                    return "[{\"name\":\"linearId\",\"type\":\"LinearId\", \"repeating\":\"False\"}]"
                } 
                if(assettype === "Fungible Asset"){
                    let value = []
                    let owner = {name:"owner", type:"Party", repeating: "False"};
                    value.push(owner)

                    let issuer = {name:"issuer", type:"Party", repeating: "False"};
                    value.push(issuer)

                    let issuerRef = {name:"issuerRef", type:"String", repeating: "False"};
                    value.push(issuerRef)

                    let quantity = {name:"quantity", type:"Long", repeating: "False"};
                    value.push(quantity)
                    let quantityUnit = {name:"quantityUnit", type:"String", repeating: "False"};
                    value.push(quantityUnit)
                    return JSON.stringify(value)
                }
                return null
           // case "schema":
           //     return this.createSchema(context, assettype == "Linear State")
                
        }
        return null;    
    }

    validate = (fieldName: string, context: IConnectorContribution): Observable<IValidationResult> | IValidationResult => {
       
        return null;
    }

    action = (name: string, context: IConnectorContribution): Observable<IActionResult> | IActionResult => {
       
       if(name === "Done"){
            return Observable.create(observer => {
                try {
                    let schema = ""
                    if(this.getFieldValue(context, "assettype") === "Fungible Asset") {
                        schema = this.createSchema(context, false)
                    } else {
                        schema = this.createSchema(context, true)
                    }
                    console.log(schema)
                    
                    observer.next(this.processModelOutput(context, schema));
                } catch(err) {
                    console.log(name + ":" + err.message);
                    return Observable.create(observer => {
                        observer.next(ActionResult.newActionResult().setSuccess(false).setResult(new ValidationError("asset schema-1002", "Action failed :" + err.message)));
                        observer.complete();
                    });
                }
            });
        } else {
            return Observable.create(observer => {
                observer.next(ActionResult.newActionResult().setSuccess(false).setResult(new ValidationError("asset schema-1001", "Invalid action :" + name)));
            });
        }
    }

    processModelOutput(context, schema) : IActionResult{
        for(let i=0; i<context.settings.length; i++){

            if(context.settings[i].name==="schema"){
                context.settings[i].value = schema;
                break
            } 
        }
        
        let actionResult = {
            context: context,
            authType: AUTHENTICATION_TYPE.BASIC,
            authData: {}
        }
        return ActionResult.newActionResult().setResult(actionResult);
    }

    getFieldValue(context, fieldName): any {
        for(let i=0; i<context.settings.length; i++){
            if(context.settings[i].name=== fieldName){
                return context.settings[i].value;
            }
        }
    }

   createSchema(context, isLinear) : any{
        var inputParams = this.getFieldValue(context, "assetfields").value
        var extraParams = this.getFieldValue(context, "extrafields").value
        var partyTypes = this.getFieldValue(context, "partyType").value
        
        if(Boolean(inputParams) == false)
            return null

       let inputs = JSON.parse(inputParams);
      
       if(extraParams){
            var extras = JSON.parse(extraParams)
            inputs = inputs.concat(extras)
       }
       let schema = {schema:"http://json-schema.org/draft-04/schema#", type: "object", properties:{}}
       let metadata = {metadata: {type:"Asset", name:"", module: "", parent:"", issueSigners:[], moveSigners:[], exitSigners:[]}, attributes:[]}
       
       if(isLinear) {
            metadata.metadata.parent = "com.tibco.dovetail.system.LinearState"
       } else {
            metadata.metadata.parent = "com.tibco.dovetail.system.FungibleAsset" 
       }
       metadata.metadata.name = this.getFieldValue(context, "name")
       metadata.metadata.module = this.getFieldValue(context, "module")
       if(partyTypes){
            let parties = JSON.parse(partyTypes);
            if(parties) {
                for(var p=0; p<parties.length; p++){
                        if(parties[p].moveSigner == "True")
                            metadata.metadata.moveSigners.push(parties[p].party)
                        if(parties[p].issueSigner == "True")
                            metadata.metadata.issueSigners.push(parties[p].party)
                        if(parties[p].exitSigner == "True")
                            metadata.metadata.exitSigners.push(parties[p].party)
                }
            }
        }

       if(inputs) {
           for(let i=0; i<inputs.length; i++){
                let name = inputs[i].name;
                let tp = inputs[i].type;
                let repeating = inputs[i].repeating;
                let partyType = "";
                let isArray = false;
                let isRef = false;
                let attr = {};
                let datatype = {type: tp.toLowerCase()};
                let systype = tp;
                switch (tp) {
                    case "Party":
                        datatype.type = "string";
                        systype = "com.tibco.dovetail.system.Party";
                        isRef = true;
                        if(metadata.metadata.moveSigners.includes(name))
                            partyType = "Participant"
                        else if(metadata.metadata.issueSigners.includes(name) || metadata.metadata.exitSigners.includes(name))
                            partyType = "Signatory"
                        else
                            partyType = "Observer"
                        break;
                    case "LinearId":
                        datatype.type = "string";
                        systype = "com.tibco.dovetail.system.UniqueIdentifier";
                        break;
                    case "Amount<Currency>":
                        datatype.type = "object";
                        datatype["properties"] = {currency: {type: "string"}, quantity: {type: "number"}};
                        systype = "com.tibco.dovetail.system.Amount<Currency>";
                        break;
                    case "Integer":
                    case "Long":
                        datatype.type = "number";
                        break;
                    case "LocalDate":
                            datatype.type = "string";
                            datatype["format"] = "date-time";
                            systype = "com.tibco.dovetail.system.LocalDate"
                            break;
                    case "DateTime":
                        datatype.type = "string";
                        datatype["format"] = "date-time";
                        systype = "com.tibco.dovetail.system.Instant"
                        break;
                }
                if(repeating === "True"){
                    schema.properties[name] = {type: "array", items: {datatype}}
                    isArray = true;
                } else {
                    schema.properties[name] = datatype
                }
    
                attr["name"] = name;
                attr["type"] = systype;
                attr["isRef"] = isRef
                attr["isArray"] = isArray;
                attr["partyType"] = partyType;
                metadata.attributes.push(attr);
           }
           schema["description"] = JSON.stringify(metadata);
           return JSON.stringify(schema);
       } else {
           return "{}";
       }
    }
}
