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

const GenericSchema = {
    "type": "object",
    "$schema": "http://json-schema.org/draft-07/schema#",
    "properties":{
        "dataset":{
            "type":"array",
            "items": {
            }
        }
    },
    "required":["dataset"]
  };

  const MergeSchema = {
    "type": "object",
    "$schema": "http://json-schema.org/draft-07/schema#",
    "properties":{
        "dataset1":{
            "type":"array",
            "items": {
            }
        },
        "dataset2":{
            "type":"array",
            "items": {
            }
        }
    },
    "required":["dataset1", "dataset2"]
  };

  const FilterInputSchema = {
    "type": "object",
    "$schema": "http://json-schema.org/draft-07/schema#",
    "properties":{
        "dataset":{
            "type":"array",
            "items": {
            }
        },
        "filterExpression": {
            "type":"string"
        }
    },
    "required":["dataset", "filterExpression"]
  };

  const FilterOutputSchema = {
    "type": "object",
    "$schema": "http://json-schema.org/draft-07/schema#",
    "properties":{
        "true_dataset":{
            "type":"array",
            "items": {
            }
        },
        "false_dataset":{
            "type":"array",
            "items": {
            }
        },
        "true_dataset_size": {"type":"integer"},
        "false_dataset_size": {"type":"integer"}
    }
  };

@WiContrib({})
@Injectable()
export class CollectionActivityContributionHandler extends WiServiceHandlerContribution {
    constructor(private injector: Injector, private http: Http) {
        super(injector, http);
        }

    value = (fieldName: string, context: IActivityContribution): any | Observable<any> => {
        let op = context.getField("operation").value;
        switch(fieldName){
            case "dataType":
                let connectionRefs = [];
                connectionRefs.push({"unique_id": "string", "name":"string"})
                connectionRefs.push({"unique_id": "number", "name":"number"})
                connectionRefs.push({"unique_id": "boolean", "name":"boolean"})
                connectionRefs.push({"unique_id": "datetime", "name":"datetime"})
                return Observable.create(observer => {
                    WiContributionUtils.getConnections(this.http, "Dovetail-Ledger").subscribe((data: IConnectorContribution[]) => {
                        data.forEach(connection => {
                            if ((<any>connection).isValid) {
                                if(connection.name === "AssetSchemaConnector"){
                                    connectionRefs.push({
                                        "unique_id": WiContributionUtils.getUniqueId(connection),
                                        "name": this.getSettingValue(connection, "displayname")
                                    });
                                }
                            }
                        });
                        observer.next(connectionRefs);
                    });
                });
           
            case "input":
                if(op){
                    return Observable.create(observer => {
                        this.createSchema(context).subscribe(s => {
                            if(op === "MERGE"){
                                var schemaJSON = lodash.cloneDeep(MergeSchema);
                                schemaJSON.properties.dataset1.items = s
                                schemaJSON.properties.dataset2.items = s
                                observer.next(JSON.stringify(schemaJSON))
                            } else if(op === "FILTER") {
                                var fschemaJSON = lodash.cloneDeep(FilterInputSchema);
                                fschemaJSON.properties.dataset.items = s
                                fschemaJSON.properties.filterExpression = s
                                observer.next(JSON.stringify(fschemaJSON))
                            } else {
                                var gschemaJSON = lodash.cloneDeep(GenericSchema);
                                gschemaJSON.properties.dataset.items = s
                                observer.next(JSON.stringify(gschemaJSON))
                            }
                        })
                    })
                } else
                    return null
            case "output":
                if(op){
                    return Observable.create(observer => {
                        this.createSchema(context).subscribe(s => {
                            if(op === "FILTER") {
                                var foschemaJSON = lodash.cloneDeep(FilterOutputSchema);
                                foschemaJSON.properties.true_dataset.items = s
                                foschemaJSON.properties.false_dataset.items = s
                                observer.next(JSON.stringify(foschemaJSON))
                            } else {
                                var goschemaJSON = lodash.cloneDeep(GenericSchema);
                                goschemaJSON.properties.dataset.items = s
                                goschemaJSON.properties["size"] = {type:"integer"}
                                observer.next(JSON.stringify(goschemaJSON))
                            }
                        })
                    })
                } else
                    return null
        }

        return null
    }
 
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        return null; 
    }
   
    createSchema(context){
        var datatype = context.getField("dataType").value
        var schema = {type:"object", properties:{}}
        return Observable.create(observer => {
            if(datatype){
                switch(datatype){
                    case "string":
                    case "number":
                    case "boolean":
                        schema.properties = {field: {type:datatype}}
                        observer.next(schema);
                        break
                    case "datetime":
                        schema.properties = {field: {type:"string", format:"date-time"}}
                        observer.next(schema);
                        break
                    default:
                        WiContributionUtils.getConnection(this.http, datatype)
                        .map(data => data)
                        .subscribe(data => {
                            var fields = []
                            var assetschema = JSON.parse(this.getSettingValue(data, "schema"))
                            schema["properties"] = assetschema.properties
                            observer.next(schema);
                        });
                }
            } else {
                observer.next(null)
            }
        })
       
    }
    getAllAssets():Observable<Map<string, string>>  {
        var schemas = new Map()
        return Observable.create(observer => {
            WiContributionUtils.getConnections(this.http, "Dovetail-Ledger").subscribe((data: IConnectorContribution[]) => {
                data.forEach(connection => {
                    if ((<any>connection).isValid) {
                        if(connection.name === "AssetSchemaConnector"){
                            var name = this.getSettingValue(connection, "name")
                            var module = this.getSettingValue(connection, "module")
                            var schema = this.getSettingValue(connection, "schema")
                          schemas.set(module+"."+name, schema)
                        }
                    }
                });
                observer.next(schemas);
            });
        });
    }
    getSettingValue(connection, setting):string {
        for(let i=0; i < connection.settings.length; i++) {
            if(connection.settings[i].name === setting){
                return connection.settings[i].value
            }
        }
    }
}