import {Injectable, Injector, Inject} from "@angular/core";
import {Http} from "@angular/http";
import {Observable} from "rxjs/Observable";
import {
    WiContrib,
    WiServiceHandlerContribution,
    IValidationResult,
    ValidationResult,
    IContributionTypes,
    IFieldDefinition,
    ActionResult,
    IActionResult,
    WiContributionUtils,
    IConnectorContribution,
    IActivityContribution
} from "wi-studio/app/contrib/wi-contrib";

@WiContrib({})
@Injectable()
export class fabrequestHandler extends WiServiceHandlerContribution {
    metadataMap: Map<string, object>

    constructor(private injector: Injector, private http: Http) {
        super(injector, http);
        this.metadataMap = new Map<string, object>()
    }

    value = (fieldName: string, context: IActivityContribution): Observable<any> | any => {
        if (fieldName === "connectionName") {
            // return list of connector refs
            return Observable.create(observer => {
                WiContributionUtils.getConnections(this.http, "fabclient").subscribe((data: IConnectorContribution[]) => {
                    let connectionRefs = [];
                    data.forEach(connection => {
                        if ((<any>connection).isValid) {
                            for (let i = 0; i < connection.settings.length; i++) {
                                if (connection.settings[i].name === "name") {
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
        } else if (fieldName === "transactionName" ) {
            let connectorId = context.getField("connectionName").value;
            if (connectorId) {
                return Observable.create(observer => {
                    WiContributionUtils.getConnection(this.http, connectorId).map(data => data)
                    .subscribe(data => {
                        this.setMetadata(connectorId, data)
                        let md = this.getMetadata(connectorId)
                        if (md) {
                            let txn = [""];
                            let con = md["contract"]
                            Object.keys(con["transactions"]).forEach( (t) => {
                                txn.push(t);
                            });
                            observer.next(txn);
                        }
                    });
                });
            }
        } else if (fieldName === "requestType" ) {
            let connectorId = context.getField("connectionName").value;
            let md = this.getMetadata(connectorId)
            let txnName = context.getField("transactionName").value;
            if (txnName && md) {
                let con = md["contract"];
                let txn = con["transactions"][txnName];
                return txn["operation"];
            }
        } else if (fieldName === "parameters" || fieldName === "transient" || fieldName === "result") {
            let txnName = context.getField("transactionName").value;
            let connectorId = context.getField("connectionName").value;
            if (txnName && connectorId) {
                return Observable.create(observer => {
                    WiContributionUtils.getConnection(this.http, connectorId).map(data => data)
                    .subscribe(data => {
                        this.setMetadata(connectorId, data)
                        if (this.getMetadata(connectorId)) {
                            let attr = fieldName;
                            if (fieldName === "result") {
                                attr = "returns";
                            }
                            let schema = this.getDataSchema(connectorId, txnName, attr)
//                            console.log("schema of " + fieldName + ": " + schema)
                            observer.next(schema);
                        }
                    });
                });
            }
        }
        return null;
    }

    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        if (fieldName === "parameters" || fieldName === "transient" || fieldName === "result") {
            let vresult: IValidationResult = ValidationResult.newValidationResult();
            let valueField: IFieldDefinition = context.getField(fieldName);
            if (valueField.value && valueField.value.value) {
                try {
                    let valRes;
                    valRes = JSON.parse(valueField.value.value);
                    valRes = JSON.stringify(valRes);
                } catch (e) {
                    vresult.setError("FABTIC-REQUEST-1000", "Invalid JSON: " + e.toString());
                }
            }
            vresult.setReadOnly(true);
            return vresult;
        }
        return null;
    }

    getDataSchema = (connectorId: string, txnName: string, attr: string): string => {
        let md = this.getMetadata(connectorId)
        let con = md["contract"];
        let txn = con["transactions"][txnName];
        let schema = txn[attr];
        if (!schema) {
            return null;
        }
        let ref = schema["$ref"];
        if (ref) {
            let shared = md["components"][ref.substring(13)];
            return JSON.stringify(shared, null, 2);
        }
        return JSON.stringify(schema, null, 2);;
    }

    setMetadata = (name: string, connector: any) => {
        if (this.getMetadata(name)) {
            console.log("metadata is already set");
            return;
        }
        for (let setting of connector.settings) {
            if (setting.name === "contract" && setting.value) {
                let content = this.extractFileContent(setting.value);
                console.log(content);
                this.metadataMap.set(name, JSON.parse(content));
            }
        }
    }

    getMetadata = (name: string): object => {
        return this.metadataMap.get(name)
    }

    extractFileContent = (selector: object): string => {
        let content = selector["content"]
        let data = content.substring(content.indexOf("base64,")+7)
//        console.log(data)
        return atob(data);
    }
}
