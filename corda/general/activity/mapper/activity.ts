import {Observable} from "rxjs/Observable";
import {Injectable, Injector, Inject} from "@angular/core";
import {Http} from "@angular/http";
import * as lodash from "lodash";
import {
    WiContrib,
    WiServiceHandlerContribution,
    IValidationResult,
    ValidationResult,
    IActivityContribution,
    WiContributionUtils,
    IConnectorContribution
} from "wi-studio/app/contrib/wi-contrib";


const numericSchema = {
    "type": "object",
    "$schema": "http://json-schema.org/draft-04/schema#",
    "required": ["field"],
    "properties": {
      "field": {
        "type": "number"
      }
    }
  };

  const stringSchema = {
    "type": "object",
    "$schema": "http://json-schema.org/draft-04/schema#",
    "required": ["field"],
    "properties": {
      "field": {
        "type": "string"
      }
    }
  };

  const booleanSchema = {
    "type": "object",
    "$schema": "http://json-schema.org/draft-04/schema#",
    "required": ["field"],
    "properties": {
      "field": {
        "type": "boolean"
      }
    }
  };
  const datetimeSchema = {
    "type": "object",
    "$schema": "http://json-schema.org/draft-04/schema#",
    "required": ["field"],
    "properties": {
      "field": {
        "type": "string",
        "format": "date-time"
      }
    }
  };
  
  const strPrimitiveArraySchema = {
    "type": "array",
    "$schema": "http://json-schema.org/draft-04/schema#",
    "items" :{
        "type": "string"
    }
  };

  const strObjectArraySchema = {
    "type": "array",
    "$schema": "http://json-schema.org/draft-04/schema#",
    "items" :{
        "type": "object",
        "properties": {
            "field": {
                "type": "string"
            }
        }
    }
  };

  const numberPrimitiveArraySchema = {
    "type": "array",
    "$schema": "http://json-schema.org/draft-04/schema#",
    "items" :{
        "type": "number"
    }
  };

  const numberObjectArraySchema = {
    "type": "array",
    "$schema": "http://json-schema.org/draft-04/schema#",
    "items" :{
        "type": "object",
        "properties": {
            "field": {
                "type": "number"
            }
        }
    }
  };

  const boolPrimitiveArraySchema = {
    "type": "array",
    "$schema": "http://json-schema.org/draft-04/schema#",
    "items" :{
        "type": "boolean"
    }
  };

  const boolObjectArraySchema = {
    "type": "array",
    "$schema": "http://json-schema.org/draft-04/schema#",
    "items" :{
        "type": "object",
        "properties": {
            "field": {
                "type": "boolean"
            }
        }
    }
  };
  const dtPrimitiveArraySchema = {
    "type": "array",
    "$schema": "http://json-schema.org/draft-04/schema#",
    "items" :{
        "type": "string",
        "format": "date-time"
    }
  };

  const dtObjectArraySchema = {
    "type": "array",
    "$schema": "http://json-schema.org/draft-04/schema#",
    "items" :{
        "type": "object",
        "properties": {
            "field": {
                "type": "string",
                "format": "date-time"
            }
        }
    }
  };
const primitiveTypes = ["Boolean","Datetime", "Double", "Integer", "Long", "String"];
const primivtePlusCustom = ["Boolean","Datetime", "Double", "Integer", "Long", "String", "User Defined..."];

@WiContrib({})
@Injectable()
export class MapperActivityContributionHandler extends WiServiceHandlerContribution {
  
  constructor(private injector: Injector, private http: Http) {
    super(injector, http);
    }
    
    value = (fieldName: string, context: IActivityContribution): any | Observable<any> => {  
        let datatype = context.getField("dataType").value;
        let isArray = context.getField("isArray").value;
        let conId = context.getField("model").value;
        
        switch(fieldName){
            case "model":
                let connectionRefs = [];
                return Observable.create(observer => {
                    WiContributionUtils.getConnections(this.http, "SmartContract").subscribe((data: IConnectorContribution[]) => {
                        data.forEach(connection => {
                            if ((<any>connection).isValid) {
                                for(let i=0; i < connection.settings.length; i++) {
                                    if(connection.settings[i].name === "name"){
                                        connectionRefs.push({
                                            "unique_id": WiContributionUtils.getUniqueId(connection),
                                            "name": connection.settings[i].value
                                        });
                                        break;
                                    }
                                }
                            }
                        });
                        observer.next(connectionRefs);
                    });
                });
            case "dataType":
                if(Boolean(conId)) {
                    return Observable.create(observer => {
                        WiContributionUtils.getConnection(this.http, conId)
                                            .map(data => data)
                                            .subscribe(data => {
                                                let types = [].concat(primitiveTypes);
                                                if(Boolean(data)){
                                                    for (let setting of data.settings) {
                                                        if (setting.name === "assets" || setting.name === "concepts") {
                                                            types = types.concat(setting.value);
                                                        }
                                                    }
                                                    types = types.concat("User Defined...");
                                                    observer.next(types);
                                                } else {
                                                    observer.next(primivtePlusCustom);
                                                }
                                            });
                    });
                } else {
                    return primivtePlusCustom;
                }
            case "input":
                if(this.isPrimtive(datatype) ) {
                    return JSON.stringify(this.getSchema(datatype, isArray, context.getField("inputArrayType").value));
                }
                
                if(Boolean(conId) == false || Boolean(datatype) == false || datatype == "User Defined...")
                        return null;

                return Observable.create(observer => {
                    this.getAssetSchemas(conId).subscribe( schemas => {
                        let schema = JSON.parse(schemas[datatype]);      
                        if(isArray && schema.type == "object"){
                            observer.next(this.createArraySchema(schema));
                        }
                        else
                            observer.next(schemas[datatype]);
                        
                    });
                });
            
            case "output":
                if (datatype ==  "User Defined..."){
                    let inschema = context.getField("userInput").value;
                    return inschema.value;
                } else {
                    if(this.isPrimtive(datatype) ) {
                        return JSON.stringify(this.getSchema(datatype, isArray, context.getField("outputArrayType").value));
                    } else {
                        let inschema = context.getField("input").value;
                        return inschema.value;
                    }
                }  
        }
       
    }
 
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        let datatype = context.getField("dataType").value;
        let isArray = context.getField("isArray").value
        switch(fieldName){
            case "precision":
            case "scale":
            case "rounding":
                return Observable.create(observer => {
                    let vresult: IValidationResult = ValidationResult.newValidationResult();
                    //TODO: will support in next release
                   // if(datatype === "Double")
                    //    vresult.setVisible(true);
                   // else
                    vresult.setVisible(false);
                    observer.next(vresult);
                });
            case "format":
                return Observable.create(observer => {
                    let vresult: IValidationResult = ValidationResult.newValidationResult();
                    //TODO: will support in next release
                   // if(datatype === "Datetime")
                    //    vresult.setVisible(true);
                   // else
                    vresult.setVisible(false);
                    observer.next(vresult);
                });
            case "isArray":
                return Observable.create(observer => {
                    let vresult: IValidationResult = ValidationResult.newValidationResult();
                    if(datatype === "User Defined...")
                        vresult.setVisible(false);
                    else 
                        vresult.setVisible(true);
                    observer.next(vresult);
                });
            case "inputArrayType":
            case "outputArrayType":
                return Observable.create(observer => {
                    let vresult: IValidationResult = ValidationResult.newValidationResult();
                    if(this.isPrimtive(datatype) && isArray)
                        vresult.setVisible(true);
                    else 
                        vresult.setVisible(false);
                    observer.next(vresult);
                });
            case "userInput":
                return Observable.create(observer => {
                    let vresult: IValidationResult = ValidationResult.newValidationResult();
                    if(datatype == "User Defined...")
                        vresult.setVisible(true);
                    else 
                        vresult.setVisible(false);
                    observer.next(vresult);
                });
            case "input":
                return Observable.create(observer => {
                    let vresult: IValidationResult = ValidationResult.newValidationResult();
                    if(datatype == "User Defined...")
                        vresult.setVisible(false);
                    else 
                        vresult.setVisible(true);
                    observer.next(vresult);
                });
            default:
                return null; 
        }
    }

    isPrimtive(dataType): boolean {
        for(let i=0; i<primitiveTypes.length; i++){
            if(primitiveTypes[i] == dataType)
                return true
        }
        return false;
    }
    getSchema(datatype, isArray, arrayType): any {
       let value;
       let inschema;
       if(isArray){

       }
        switch(datatype){
            case "Integer":
            case "Long":
            case "Double":
                if(isArray) {
                    if(arrayType == "Object Array")
                        value = lodash.cloneDeep(numberObjectArraySchema);
                    else
                        value = lodash.cloneDeep(numberPrimitiveArraySchema)
                } else {
                    value = lodash.cloneDeep(numericSchema);
                }
                break;
            case "Boolean":
                if(isArray) {
                    if(arrayType == "Object Array")
                    value = lodash.cloneDeep(boolObjectArraySchema);
                else
                    value = lodash.cloneDeep(boolPrimitiveArraySchema)
                }
                else
                    value = lodash.cloneDeep(booleanSchema);
                break;
            case "Datetime":
                if(isArray){
                    if(arrayType == "Object Array")
                        value = lodash.cloneDeep(dtObjectArraySchema);
                    else
                        value = lodash.cloneDeep(dtPrimitiveArraySchema)
                } else
                    value = lodash.cloneDeep(datetimeSchema);
                break;
            default:
                if(isArray){
                    if(arrayType == "Object Array")
                        value = lodash.cloneDeep(strObjectArraySchema);
                    else
                        value = lodash.cloneDeep(strPrimitiveArraySchema)
                } else
                    value = lodash.cloneDeep(stringSchema);
                break;           
        }

        return value
    }

    getAssetSchemas(conId):  Observable<any> {
        let schemas = new Map();
        return Observable.create(observer => {
            WiContributionUtils.getConnection(this.http, conId)
                            .map(data => data)
                            .subscribe(data => {
                                let schemas = new Map();
                                for (let setting of data.settings) {
                                    if(setting.name === "schemas") {
                                        setting.value.map(item => schemas[item[0]] = item[1]);
                                        observer.next(schemas);
                                        break;
                                    }
                                }
                            });
                        });
    }

    createArraySchema(schema): string {
        let newSchema = {};
        newSchema["$schema"] = schema["$schema"];
        newSchema["title"] = schema["title"];
        newSchema["type"] = "array"
        newSchema["items"] = {type: "object", properties:{}};
        newSchema["items"].properties = schema.properties;
        newSchema["description"] = schema.description;
        return JSON.stringify(newSchema);
    }
}