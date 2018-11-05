import {Observable} from "rxjs/Observable";
import {Injectable, Injector, Inject} from "@angular/core";
import { Http, Response } from "@angular/http";
import {
    WiContrib,
    WiServiceHandlerContribution,
    IValidationResult,
    ValidationResult,
    IActivityContribution,
    WiContributionUtils,
    IConnectorContribution
} from "wi-studio/app/contrib/wi-contrib";
import * as lodash from "lodash";

const stringInputSchema = {
    "type": "array",
    "$schema": "http://json-schema.org/draft-07/schema#",
    "items": {
      "type": "object",
      "properties": {
        "data": {
          "type": "string"
        }
      }
    }
  };

  const stringOutputSchema = {
    "type": "object",
    "$schema": "http://json-schema.org/draft-07/schema#",
      "properties": {
        "result": {
          "type": "string"
        }
      }
  };

const distinctOutputSchema = {
    "type": "object",
    "$schema": "http://json-schema.org/draft-07/schema#",
    "required": ["result"],
    "properties": {
      "result": {
        "type": "array",
        "items": {
          "type": "string"
        }
      },
      "count" :{
        "type": "integer"
      }
    }
  };

  const countOutputSchema = {
    "type": "object",
    "$schema": "http://json-schema.org/draft-07/schema#",
    "properties": {
      "count" :{
        "type": "integer"
      }
    }
  };

@WiContrib({})
@Injectable()
export class CollectionActivityContributionHandler extends WiServiceHandlerContribution {
    constructor(private injector: Injector, private http: Http) {
        super(injector, http);
        }

    value = (fieldName: string, context: IActivityContribution): any | Observable<any> => {
        let conId = context.getField("model").value;
        let op = context.getField("operation").value;
        let datatype = context.getField("dataType").value;
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
                if(op == "DISTINCT" || op == "REDUCE-JOIN"){
                    return ["String"];
                } else {
                    if(Boolean(conId)) {
                        return Observable.create(observer => {
                            WiContributionUtils.getConnection(this.http, conId)
                                                .map(data => data)
                                                .subscribe(data => {
                                                    let types = [];
                                                    if(Boolean(data)){
                                                        for (let setting of data.settings) {
                                                            if (setting.name === "assets" || setting.name === "concepts") {
                                                                types = types.concat(setting.value);
                                                            }
                                                        }

                                                        types = types.concat("User Defined...");
                                                        observer.next(types);
                                                    } else {
                                                        observer.next(["User Defined..."]);
                                                    }
                                                });
                        });
                    } else {
                        return ["User Defined..."];
                    }
                }
            case "input":
                if (op == "DISTINCT" || op == "REDUCE-JOIN"){
                    return Observable.create(observer => {
                        let schemaJSON = lodash.cloneDeep(stringInputSchema);
                        observer.next(JSON.stringify(schemaJSON));
                    });
                } else {
                    if(datatype == "" || datatype == "User Defined...")
                            return null;

                    return Observable.create(observer => {
                        this.getAssetSchemas(conId).subscribe( schemas => {
                            if(Boolean(schemas[datatype])) {
                                let schema = JSON.parse(schemas[datatype]);    
                                if (op == "MERGE") {
                                    observer.next(this.createMergeInputSchema(schema));
                                } else if (op == "FILTER"){
                                    observer.next(this.createFilterInputSchema(schema, context));
                                } else {
                                    if(schema.type == "object"){
                                        observer.next(this.createArraySchema(schema));
                                    }
                                    else
                                        observer.next(schemas[datatype]);
                                }
                            } else {
                                observer.next(null);
                            }
                            
                        });
                    });
                }
            case "output":
                if (op == "DISTINCT"){
                    return Observable.create(observer => {
                        let schemaJSON = lodash.cloneDeep(distinctOutputSchema);
                        observer.next(JSON.stringify(schemaJSON));
                    });
                } else if(op == "COUNT"){
                    return Observable.create(observer => {
                        let schemaJSON = lodash.cloneDeep(countOutputSchema);
                        observer.next(JSON.stringify(schemaJSON));
                    });
                } else if (op == "REDUCE-JOIN") {
                    return Observable.create(observer => {
                        let schemaJSON = lodash.cloneDeep(stringOutputSchema);
                        observer.next(JSON.stringify(schemaJSON));
                    });
                } else {
                    let value: string;
                    let inschema;
                    let schema;
                    switch(datatype){
                        case "User Defined...":
                            inschema = context.getField("userInput").value;     
                            break;
                        default:
                            inschema = context.getField("input").value;
                            break;
                    }
                    if(Boolean(inschema)){
                        value = inschema.value;
                        if(Boolean(value)){
                            schema = JSON.parse(value);
                            if (op == "MERGE") {
                                return this.createMergeOutputSchema(schema);
                            } else if (op == "FILTER") {
                                return this.createFilterOutputSchema(schema);
                            } else if (op == "INDEXING") {  
                                return this.createIndexingOutputSchema(schema);
                            } 
                        }
                    } 
                    
                    return null;
                }
        }

        return null
    }
 
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        let datatype = context.getField("dataType").value;
        let op = context.getField("operation").value;
        switch(fieldName){
            case "userInput":
                return Observable.create(observer => {
                    let vresult: IValidationResult = ValidationResult.newValidationResult();
                    if(datatype == "User Defined..."){
                        vresult.setVisible(true);
                        let input = context.getField("userInput").value.value
                        if(Boolean(input)){
                            let schema = JSON.parse(input)
                            if(Boolean(schema["$schema"])) {
                                if(op == "MERGE"){
                                    if(schema.type != "object" || !Boolean(schema.properties.input1) || !Boolean(schema.properties.input2 || schema.properties.input1.type != "array" || schema.properties.input2.type != "array")){
                                        vresult.setValid(false).setError("INVALID_SCHEMA", `Expected {"input1":[{"field1":"value1", ...}], "input2":[{"field1":"value1", ...}]}`);
                                    }
                                } else if (op == "FILTER") {
                                    if(schema.type != "object" || !Boolean(schema.properties.filterField) || !Boolean(schema.properties.filterValue) || !Boolean(schema.properties.dataset || schema.properties.dataset.type != "array")){
                                        vresult.setValid(false).setError("INVALID_SCHEMA", `Expected {"filterField":"$dataset.path.to.field", "filterValue": "value", dataset":[{"field1":"value1", ...}]}`);
                                    }
                                }else {
                                    if(!Boolean(schema.items))
                                        vresult.setValid(false).setError("INVALID_SCHEMA", "Array is expected");

                                } 
                            } else {
                                if(op == "MERGE"){
                                    if(!Boolean(schema.input1) || !Boolean(schema.input2 || schema.input1.length == 0 || schema.input2.length == 0)){
                                        vresult.setValid(false).setError("INVALID_SCHEMA", `Expected {"input1":[{"field1":"value1", ...}], "input2":[{"field1":"value1", ...}]}`);
                                    }
                                } else if (op == "FILTER") {
                                    if(!Boolean(schema.filterValue) || !Boolean(schema.dataset || schema.dataset.length == 0)){
                                        vresult.setValid(false).setError("INVALID_SCHEMA", `Expected {"filterField":"$dataset.path.to.field", "filterValue": "value", "dataset":[{"field1":"value1", ...}]}`);
                                    }
                                }else {
                                    if(!input.trim().startsWith("["))
                                        vresult.setValid(false).setError("INVALID_SCHEMA", "Array is expected");

                                }
                            }
                        } 
                    }
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
            case "delimiter":
                return Observable.create(observer => {
                    let vresult: IValidationResult = ValidationResult.newValidationResult();
                    if(op == "REDUCE-JOIN")
                        vresult.setVisible(true);
                    else 
                        vresult.setVisible(false);
                    observer.next(vresult);
                });
                
            case "filterFieldType":
            case "filterFieldOp":
                return Observable.create(observer => {
                    let vresult: IValidationResult = ValidationResult.newValidationResult();
                    if(context.getField("operation").value == "FILTER")
                        vresult.setVisible(true);
                    else 
                        vresult.setVisible(false);

                    observer.next(vresult);
                });
            case "model":
                return Observable.create(observer => {
                    let vresult: IValidationResult = ValidationResult.newValidationResult();
                    if(context.getField("operation").value == "REDUCE-JOIN")
                        vresult.setVisible(false);
                    else 
                        vresult.setVisible(true);

                    observer.next(vresult);
                });
        }
        return null; 
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

    createMergeInputSchema(schema): string {
        let newSchema = {};
        newSchema["$schema"] = schema["$schema"];
        newSchema["title"] = schema["title"];
        newSchema["type"] = "object"
        newSchema["properties"] = {input1:{type:"array"}, "input2":{type:"array"}};
        newSchema["properties"].input1["items"]={type: "object", properties:{}};
        newSchema["properties"].input2["items"]={type: "object", properties:{}};
        newSchema["properties"].input1["items"].properties = schema.properties;
        newSchema["properties"].input2["items"].properties = schema.properties;
        return JSON.stringify(newSchema);
    }

    createMergeOutputSchema(schema): string {
        let newSchema = {};
        if(Boolean(schema["$schema"])){ 
            if(Boolean(schema.properties) && Boolean(schema.properties.input1)){
                newSchema["$schema"] = schema["$schema"];
                newSchema["title"] = schema["title"];
                newSchema["type"] = "array"
                newSchema["items"] = schema.properties.input1.items
                return JSON.stringify(newSchema);
            }
        } else {
            //user defined
            return JSON.stringify(schema.input1);
        }
    }

    createFilterOutputSchema(schema): string {
        let newSchema = {};   
        if(Boolean(schema["$schema"])){   
            if(Boolean(schema.properties) && Boolean(schema.properties.dataset)){
                newSchema["$schema"] = schema["$schema"];
                newSchema["title"] = schema["title"];
                newSchema["type"] = "object"
                newSchema["properties"] = {trueSet:{type:"array"}, falseSet:{type:"array"}};
                newSchema["properties"].trueSet["items"] = schema.properties.dataset.items;
                newSchema["properties"].falseSet["items"] = schema.properties.dataset.items;
            }
        } else {
            //user defined
            if(Boolean(schema.dataset)){
                newSchema["trueSet"] = schema.dataset
                newSchema["falseSet"] = schema.dataset
            }
                
        }
        return JSON.stringify(newSchema);
    }

    createFilterInputSchema(schema, context): string {
        let dataType = context.getField("filterFieldType").value.toLowerCase();
        let op = context.getField("filterFieldOp").value;

        let newSchema = {};
        newSchema["$schema"] = schema["$schema"];
        newSchema["title"] = schema["title"];
        newSchema["type"] = "object"
        newSchema["properties"] = {filterField:{}, filterValue:{}, dataset:{type:"array"}};
        newSchema["properties"].filterField = {type: "string"};
        if(op == "IN"){
            newSchema["properties"].filterValue = {type: "array", items: {type: "object", properties:{value: {type: dataType}}}};
        } else {
            newSchema["properties"].filterValue = {type: dataType};
        }
        newSchema["properties"].dataset["items"]={type: "object", properties:{}};
        newSchema["properties"].dataset["items"].properties = schema.properties;
        newSchema["description"] = schema.description;
        return JSON.stringify(newSchema);
    }

    createIndexingOutputSchema(schema): string {
        if(Boolean(schema["$schema"])){
            if(Boolean(schema.items.properties))
            //object array
                schema.items.properties["_index_"] = {type: "integer"};
            else {
            //change primitive array ["a", "b"] to object array
                let item = {type: "object"};
            
                item["properties"] = {_index_ : {type: "integer"}, data: {type: schema.items.type}}
                schema.items = item
            }
        } else {
            let arrayschema = []
            let item = {}
            if(typeof(schema[0]) == "object"){
                item = schema[0]
                item["_index_"] = 0
            } else {
                item["data"] = schema[0]
                item["_index_"] = 0
            }
            
            arrayschema.push(item)
            schema = arrayschema
        }

        return JSON.stringify(schema);
    }
}