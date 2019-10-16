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
export class TxnFilterActivityContributionHandler extends WiServiceHandlerContribution {
    constructor(private injector: Injector, private http: Http,) {
        super(injector, http);
    }
   
    value = (fieldName: string, context: IActivityContribution): any | Observable<any> => {
        let filterby = context.getField("filterby").value;
        let asset = context.getField("assets").value;
        
        switch(fieldName) {
            case "assets":
                if (filterby === "Input State" || filterby === "Output State" || filterby === "Reference State"){
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
                            connectionRefs.push({
                                "unique_id": "com.tibco.dovetail.system.Cash",
                                "name": "Cash"
                            });
                            observer.next(connectionRefs);
                        });
                    });
                } else
                    return null;
            case "assetName":
                if(asset === "com.tibco.dovetail.system.Cash")
                    return asset;

                return Observable.create(observer => {
                    WiContributionUtils.getConnection(this.http, asset)
                                    .map(data => data)
                                    .subscribe(data => {
                                        var ns = this.getSettingValue(data, "module")
                                        var nm = this.getSettingValue(data, "name")
                                        observer.next(ns + "." + nm);
                                                
                                    })
                });
            case "output":
                switch(filterby){
                    case "Input State" :
                    case "Output State" :
                    case "Reference State":
                        if(Boolean(asset) == false)
                            return null;

                        if(asset === "com.tibco.dovetail.system.Cash")
                            return "{\"$schema\":\"http://json-schema.org/draft-04/schema#\",\"type\":\"array\",\"items\":{\"type\":\"object\",\"properties\":{\"ref\":{\"type\":\"string\"},\"data\":{\"type\":\"object\", \"properties\": {\"amt\":{\"type\":\"object\",\"properties\":{\"currency\":{\"type\":\"string\"},\"quantity\":{\"type\":\"integer\"}}},\"issuer\":{\"type\":\"string\"},\"issuerRef\":{\"type\":\"string\"},\"owner\":{\"type\":\"string\"}}}}}}";

                        return Observable.create(observer => {
                            this.getSchema(asset).subscribe( assetschema => {
                                let schema = JSON.parse(assetschema);
                                
                                let newSchema = {};
                                newSchema["$schema"] = schema["$schema"];
                                newSchema["title"] = schema["title"];
                                newSchema["type"] = "array"
                                newSchema["items"] = {type: "object", properties:{}};
                                newSchema["items"].properties = schema.properties;
                                newSchema["description"] = schema.description;
                                observer.next(JSON.stringify(newSchema));
                            });
                        });
                case  "Commands":
                    let cmds = {
                        "$schema":"http://json-schema.org/draft-07/schema#",
                        "type": "array",
                        "items": {
                            "type": "object",
                            "properties": {
                                "command": {
                                    "type": "string"
                                }
                            }
                        }
                    };
                    return JSON.stringify(cmds);
                case "Attachment":
                    let att = {
                        "$schema":"http://json-schema.org/draft-07/schema#",
                        "type": "object",
                        "properties": {
                            "secureHash": {
                                "type": "array",
                                "items": {
                                    "type": "string"
                                }
                            }
                        }
                    };
                    return JSON.stringify(att);
                case "Time Window":
                    let tw = {
                        "$schema":"http://json-schema.org/draft-07/schema#",
                        "type": "object",
                        "properties": {
                            "from": {
                                "type": "string",
                                "format": "date-time"
                            },
                            "until": {
                                "type": "string",
                                "format": "date-time"
                            },
                            "duration": {
                                "type": "string"
                            }
                        }
                    };
                    return JSON.stringify(tw);
                case "Notary":
                    let notary = {
                        "$schema":"http://json-schema.org/draft-07/schema#",
                        "type": "object",
                        "properties": {
                            "notary": {
                                "type": "string"
                            }
                        }
                    };
                    return JSON.stringify(notary);
            }
        }
        return null;
    }
 
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        let filterby = context.getField("filterby").value;
        let visible = false;
        switch(fieldName){
            case "assets":
                if (filterby === "Input State" || filterby === "Output State" || filterby === "Reference State"){
                    visible = true
                } 

                return Observable.create(observer => {
                    let vresult: IValidationResult = ValidationResult.newValidationResult();
                    vresult.setVisible(visible);
                    observer.next(vresult);
                });
        }
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