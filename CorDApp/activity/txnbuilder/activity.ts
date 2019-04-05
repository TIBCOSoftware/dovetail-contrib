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

let inputschema = "{\"$schema\": \"http://json-schema.org/draft-07/schema#\",\"type\": \"object\",\"properties\": {\"txnHash\": {\"type\": \"string\"}, \"index\":{\"type\":\"integer\"}}}"

@WiContrib({})
@Injectable()
export class TxnBuilderActivityContributionHandler extends WiServiceHandlerContribution {
    constructor(private injector: Injector, private http: Http,) {
        super(injector, http);
    }
   
    value = (fieldName: string, context: IActivityContribution): any | Observable<any> => {
        let conId = context.getField("contract").value;
        let cmd = context.getField("command").value;
        switch(fieldName) {
            case "contract":
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
            case "command":
                if(Boolean(conId) == false)
                    return null;

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
            case "inputSchema":
                    if(Boolean(conId) == false || Boolean(cmd) == false)
                        return null;

                    return Observable.create(observer => {
                        this.getSchemas(conId).subscribe( schemas => {
                            
                            let schema = JSON.parse(schemas[cmd])
                            let metadata = JSON.parse(schema.description);
                            for (let attr of metadata.attributes) {
                                if(Boolean(schemas[attr.type])) {
                                    let typeschema = JSON.parse(schemas[attr.type])
                                    let typemetadata = JSON.parse(typeschema.description)
                                    if(typemetadata.metadata.type == "Asset"){
                                        attr["isAsset"] = true
                                    } else if(typemetadata.metadata.type == "Participant"){
                                        attr["isParty"] = true
                                    }
                                }
                            }
                            observer.next(JSON.stringify(metadata.attributes));
                        });
                    });  
            case "input":
                if(Boolean(conId) == false || Boolean(cmd) == false)
                    return null;

                return Observable.create(observer => {
                    this.getSchemas(conId).subscribe( schemas => {
                        
                        let schema = JSON.parse(schemas[cmd])
                        let metadata = JSON.parse(schema.description);

                        for (let attr of metadata.attributes) {
                            if(attr.isRef) {
                                let typeschema = JSON.parse(schemas[attr.type])
                                let typemetadata = JSON.parse(typeschema.description)
                                if(typemetadata.metadata.type == "Asset"){
                                    if(attr.isArray){
                                        schema.properties[attr.name] = {type: "array", items: {type: "string"}}
                                    } else {
                                        schema.properties[attr.name] = {type: "string"}
                                    }
                                }
                            }
                        }


                        observer.next(JSON.stringify(schema));
                    });
                });  
        }
        return null;
    }
 
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
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