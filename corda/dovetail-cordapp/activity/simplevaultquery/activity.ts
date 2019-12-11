import {Observable} from "rxjs/Observable";
import {Injectable, Injector, Inject} from "@angular/core";
import {Http} from "@angular/http";
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

//let linearschema = "{\"$schema\": \"http://json-schema.org/draft-07/schema#\",\"type\": \"object\",\"properties\": {\"uuid\": {\"type\": \"string\"}, \"externalId\": {\"type\":\"string\"}}}"
let linearschema = "{\"linearId\": \"\"}"

@WiContrib({})
@Injectable()
export class VaultQueryActivityContributionHandler extends WiServiceHandlerContribution {
    constructor(private injector: Injector, private http: Http,) {
        super(injector, http);
    }
   
    value = (fieldName: string, context: IActivityContribution): any | Observable<any> => {
        let assetType = context.getField("assetType").value;
        let conId = context.getField("assets").value;
        switch(fieldName) {
            case "assets":
                let connectionRefs = [];
                return Observable.create(observer => {
                    WiContributionUtils.getConnections(this.http, "Dovetail-Ledger").subscribe((data: IConnectorContribution[]) => {
                        data.forEach(connection => {
                            if ((<any>connection).isValid) {
                                for(let i=0; i < connection.settings.length; i++) {
                                    if(connection.settings[i].name === "displayname"){
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
            case "assetName":
                return Observable.create(observer => {
                    WiContributionUtils.getConnection(this.http, conId)
                                    .map(data => data)
                                    .subscribe(data => {
                                        var ns = this.getSettingValue(data, "module")
                                        var nm = this.getSettingValue(data, "name")
                                        observer.next(ns + "." + nm);
                                               
                                    })
                });
                    
            case "input":
                if(assetType === "LinearState")
                    return linearschema;
                else 
                    return null;
            case "output":
                if(Boolean(conId) == false)
                    return null;

                return Observable.create(observer => {
                    this.getSchema(conId).subscribe( assetschema => {
                        let schema = JSON.parse(assetschema);
                        
                        let newSchema = {};
                        newSchema["$schema"] = schema["$schema"];
                        newSchema["title"] = schema["title"];
                        newSchema["type"] = "array"
                        newSchema["items"] = {type: "object", properties:{}};
                        let properties = {};
                        properties["data"] = {type: "object", properties: {}};
                        properties["data"].properties = schema.properties;
                        properties["ref"] = {type: "string"}
                        newSchema["items"].properties = properties;
                        newSchema["description"] = schema.description;
                        observer.next(JSON.stringify(newSchema));
                    });
                });
        }
        return null;
    }
 
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        return null;
    }

    getSchema(conId):  Observable<any> {
        return Observable.create(observer => {
            WiContributionUtils.getConnection(this.http, conId)
                            .map(data => data)
                            .subscribe(data => {
                                for (let setting of data.settings) {
                                    if(setting.name === "schema") {
                                        observer.next(setting.value);
                                        break;
                                    }
                                }
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