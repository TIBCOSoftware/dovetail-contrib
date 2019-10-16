import { Inject, Injectable, Injector } from "@angular/core";
import { WiContrib, WiServiceHandlerContribution, WiContributionUtils,AUTHENTICATION_TYPE } from "wi-studio/app/contrib/wi-contrib";
import { IConnectorContribution } from "wi-studio/common/models/contrib";
import { IValidationResult, ValidationResult, ValidationError } from "wi-studio/common/models/validation";
import { IActionResult, ActionResult } from "wi-studio/common/models/contrib";
import { Observable } from "rxjs/Observable";

@Injectable()
@WiContrib({})
export class StructConnectorService extends WiServiceHandlerContribution {
    constructor( @Inject(Injector) injector) {
        super(injector);
    }

    value = (fieldName: string, context: IConnectorContribution): Observable<any> | any => {
        return null;    
    }

    validate = (fieldName: string, context: IConnectorContribution): Observable<IValidationResult> | IValidationResult => {
       return null;
    }

    action = (name: string, context: IConnectorContribution): Observable<IActionResult> | IActionResult => {
        if(name === "Finish"){
            return Observable.create(observer => {
                try {
                        var   schema = this.createSchema(context)
                        observer.next(this.processModelOutput(context, schema));
                } catch(err) {
                    console.log(name + ":" + err.message);
                    return Observable.create(observer => {
                        observer.next(ActionResult.newActionResult().setSuccess(false).setResult(new ValidationError("struct schema-1002", "Action failed :" + err.message)));
                        observer.complete();
                    });
                }
            });
        } else {
            return Observable.create(observer => {
                observer.next(ActionResult.newActionResult().setSuccess(false).setResult(new ValidationError("struct schema-1001", "Invalid action :" + name)));
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

    createSchema(context) : any{
        var inputParams = this.getFieldValue(context, "inputParams").value
        
        if(Boolean(inputParams) == false)
            return null

       let inputs = JSON.parse(inputParams);
      
       let schema = {schema:"http://json-schema.org/draft-04/schema#", type: "object", properties:{}}
       let metadata = {metadata: {type:"Struct", name:"", module: ""}, attributes:[]}
      
       metadata.metadata.name = this.getFieldValue(context, "name")
       metadata.metadata.module = this.getFieldValue(context, "module")
       if(inputs) {
           for(let i=0; i<inputs.length; i++){
                let name = inputs[i].name;
                let tp = inputs[i].type;
                let repeating = inputs[i].repeating;
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
                    case "Decimal":
                        datatype.type = "string"
                        systype = "Decimal"
                        break
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
                    schema.properties[name] = {type: "array", items: datatype}
                    isArray = true;
                } else {
                    schema.properties[name] = datatype
                }
    
                attr["name"] = name;
                attr["type"] = systype;
                attr["isRef"] = isRef
                attr["isArray"] = isArray;
                metadata.attributes.push(attr);
           }
           schema["description"] = JSON.stringify(metadata);
           return JSON.stringify(schema);
       } else {
           return "{}";
       }
    }
}
