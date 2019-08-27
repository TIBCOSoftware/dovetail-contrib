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
        let conId = context.getField("model").value;
        let asset = context.getField("assetName").value;
        
        switch(fieldName) {
            case "model":
                if (filterby === "Input State" || filterby === "Output State" || filterby === "Reference State" || filterby === "Command"){
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
                } else
                    return null;
                
            case "assetName":
                if (filterby === "Input State" || filterby === "Output State" || filterby === "Reference State" || filterby === "Command"){

                    if(Boolean(conId) == false)
                        return null;

                    if (filterby === "Input State" || filterby === "Output State" || filterby === "Reference State"){
                        return Observable.create(observer => {
                            WiContributionUtils.getConnection(this.http, conId)
                                                .map(data => data)
                                                .subscribe(data => {
                                                    for (let setting of data.settings) {
                                                        if (setting.name === "assets") {
                                                            observer.next(setting.value);
                                                            break;
                                                        }
                                                    }
                                                });
                        });
                    } else {
                        return Observable.create(observer => {
                            WiContributionUtils.getConnection(this.http, conId)
                                                .map(data => data)
                                                .subscribe(data => {
                                                    for (let setting of data.settings) {
                                                        if (setting.name === "transactions") {
                                                            observer.next(setting.value);
                                                            break;
                                                        }
                                                    }
                                                });
                        });
                    }
                } else 
                    return null;

            case "output":
                switch(filterby){
                    case "Input State" :
                    case "Output State" :
                    case "Reference State":
                        if(Boolean(conId) == false || Boolean(asset) == false)
                            return null;

                        return Observable.create(observer => {
                            this.getSchemas(conId).subscribe( schemas => {
                                let schema = JSON.parse(schemas[asset]);
                                
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
                case  "Command":
                    let cmds = {
                        "$schema":"http://json-schema.org/draft-07/schema#",
                        "type": "object",
                        "properties": {
                            "command": {
                                "type": "string"
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
            case "model":
            case "assetName":
                if (filterby === "Input State" || filterby === "Output State" || filterby === "Reference State" || filterby === "Command"){
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

    getSchemas(conId):  Observable<any> {
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
}