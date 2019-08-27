import {Observable} from "rxjs/Observable";
import {Injectable, Injector, Inject} from "@angular/core";
import {Http} from "@angular/http";
import {
    WiContrib,
    WiServiceHandlerContribution,
    IValidationResult,
    ValidationResult,
    IActivityContribution,
} from "wi-studio/app/contrib/wi-contrib";


@WiContrib({})
@Injectable()
export class ScheduledTaskActivityContributionHandler extends WiServiceHandlerContribution {
    constructor(private injector: Injector, private http: Http,) {
        super(injector, http);
    }
   
    value = (fieldName: string, context: IActivityContribution): any | Observable<any> => {
        switch(fieldName) {
            case "input":
                let schema = {schema:"http://json-schema.org/draft-04/schema#", type: "object", properties:{scheduledAt:{type: "string", format: "date-time"}}, "required":["scheduledAt"]};
                return JSON.stringify(schema);
        }
        return null;
    }
 
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        return null;
    }

    createFlowInputSchema(inputParams):String {
        let schema = {schema:"http://json-schema.org/draft-04/schema#", type: "object", properties:{scheduledAt:{type: "string", format: "date-time"}}, "required":["scheduledAt"]}
        schema.properties["scheduledAt"] = {type: "string", format: "date-time"}
        let metadata = {metadata: {type: "Transaction"}, attributes: []};

        let inputs = JSON.parse(inputParams);
      
        if(inputs) {
            schema.properties["flowInputs"] = {type: "object", properties: {}}
            for(let i=0; i<inputs.length; i++){
                    let name = inputs[i].parameterName;
                    let tp = inputs[i].type;
                    let repeating = inputs[i].repeating;
                    let javatype = tp;
                    let isRef = false;
                    let isArray = false;
                    let attr = {};
                    let datatype = {type: tp.toLowerCase()};

                    switch (tp) {
                        case "Party":
                            datatype.type = "string";
                            javatype = "net.corda.core.identity.Party";
                            isRef = true;
                            break;
                        case "LinearId":
                            datatype.type = "string";
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
                            break;
                        case "LocalDate":
                            datatype.type = "string";
                            datatype["format"] = "date-time";
                            javatype = "java.time.LocalDate"
                            break;
                        case "DateTime":
                            datatype.type = "string";
                            datatype["format"] = "date-time";
                            javatype = "java.time.Instant"
                            break;
                    }
                    if(repeating === "true"){
                        schema.properties["flowInputs"].properties[name] = {type: "array", items: {datatype}}
                    } else {
                        schema.properties["flowInputs"].properties[name] = datatype
                    }

                    attr["name"] = name;
                    attr["type"] = javatype;
                    attr["isRef"] = isRef
                    attr["isArray"] = isArray;
                    metadata.attributes.push(attr);
            }
        } 
        schema["description"] = JSON.stringify(metadata);
        return JSON.stringify(schema);
    }
}