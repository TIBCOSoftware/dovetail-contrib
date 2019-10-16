import {Observable} from "rxjs/Observable";
import {Inject, Injectable, Injector} from "@angular/core";
import {Http} from "@angular/http";
import {
    WiContrib,
    WiServiceHandlerContribution,
    IValidationResult,
    ValidationResult,
    IActivityContribution,
    WiContributionUtils,
    IConnectorContribution,
    WiContribModelService,
    IFlow,
    IFlowElement,
} from "wi-studio/app/contrib/wi-contrib";

import * as lodash from "lodash";

const keyschema = {
    "type": "object",
    "$schema": "http://json-schema.org/draft-07/schema#",
    "properties": {}
  };

const arraykeyschema = {
       "type": "array", 
       "$schema": "http://json-schema.org/draft-07/schema#", 
       "items": {
         "type": "object", 
         "properties": {}
      }
    };

@WiContrib({})
@Injectable()
export class LedgerActivityContributionHandler extends WiServiceHandlerContribution {
    constructor(@Inject(Injector) injector, private http: Http, private contribModelService: WiContribModelService) {
        super(injector, http, contribModelService);
    }
   
    value = (fieldName: string, context: IActivityContribution): any | Observable<any> => {
        let action = context.getField("operation").value;
        let isArray = false;
        if(Boolean(context.getField("isArray"))){
            isArray = context.getField("isArray").value;
        }
     
        switch(fieldName){
            case "asset":
                return this.getTxnAsset(context)
            case "txn":
                return context.getCurrentFlowName()
            case "input":
                return Observable.create(observer => {
                    this.getAssetSchema(context).subscribe(data => {

                        let schema = JSON.parse(data);
                        let metadata = JSON.parse(schema.description);
                        let parent = metadata.metadata.parent;
                        if(action === "DELETE" || action == "GET") {
                            if(parent === "com.tibco.dovetail.system.LinearState"){
                                observer.next(this.createKeySchema("linearId", isArray));
                            } else {
                                observer.next(JSON.stringify(schema));
                            }
                        } 
                        else {
                            if(isArray && schema.type == "object"){
                                observer.next(this.createArraySchema(schema));
                            }
                            else
                                observer.next(JSON.stringify(schema));
                        }
                    });
                });
            
            case "output":
                
                return Observable.create(observer => {
                    this.getAssetSchema(context).subscribe(data => {

                        let schema = JSON.parse(data);
                        let metadata = JSON.parse(schema.description);
                        let parent = metadata.metadata.parent;
                        if(action === "DELETE" || action == "GET") {
                            if(parent === "com.tibco.dovetail.system.LinearState"){
                                observer.next(this.createKeySchema("linearId", isArray));
                            } else {
                                observer.next(JSON.stringify(schema));
                            }
                        } 
                        else {
                            if(isArray && schema.type == "object"){
                                observer.next(this.createArraySchema(schema));
                            }
                            else
                                observer.next(JSON.stringify(schema));
                        }
                    });
                });
            default:
                return null;
        }
    }
 
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        switch(fieldName){
            case "input":
                return Observable.create(observer => {
                    let vresult: IValidationResult = ValidationResult.newValidationResult();
                    vresult.setReadOnly(true);
                    observer.next(vresult);
                });
            default:
                return null; 
        }
    }

    getAssetSchema(context):  Observable<any> {
        var asset = this.getTxnAsset(context)
        var schema
        return Observable.create(observer => {
            WiContributionUtils.getConnections(this.http, "Dovetail-Ledger").subscribe((data: IConnectorContribution[]) => {
                for(var connection of data){
                    if ((<any>connection).isValid) {  
                        if(connection.name === "AssetSchemaConnector"){
                            var name = this.getSettingValue(connection, "name")
                            var ns = this.getSettingValue(connection, "module")
                            if(ns+"."+name=== asset){
                                schema = this.getSettingValue(connection, "schema") 
                                break
                            }
                        }      
                    }
                }
                observer.next(schema);
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

    createKeySchema(key, isArray): string {
        var p = {}
        p[key] = {type:"string"};
        
        let idschema 
        if(isArray){
          idschema = lodash.cloneDeep(arraykeyschema);
          idschema.items.properties = p;
        }
        else{
           idschema = lodash.cloneDeep(keyschema);
           idschema.properties = p;
        }
        
        idschema["required"] = key;

        return JSON.stringify(idschema);
    }

    getTxnAsset(context):string {
        let modelService = this.getModelService();
        let applicationModel = modelService.getApplication();
        let asset;
        if (applicationModel) {
            let triggerMappings = applicationModel.getTriggerFlowModelMaps();
            triggerMappings.map(triggerMapping => {
                if ((context.getCurrentFlowName() === triggerMapping.getFlowModel().getName()) && !triggerMapping.getFlowModel().isTriggerFlow()) {
                    var schema = JSON.parse(triggerMapping.getFlowModel().getFlowInputSchema().json);
                    var desc = schema.description
                    if(Boolean(desc) === false){
                        desc = schema.properties.transactionInput.description
                    }
                    var metadata = JSON.parse(desc)
                    asset = metadata.metadata.asset
                }
            })
        };

        return asset
    }
}