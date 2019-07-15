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

    metadata: object

    constructor(private injector: Injector, private http: Http) {
        super(injector, http);
    }

    value = (fieldName: string, context: IActivityContribution): Observable<any> | any => {
        if (fieldName === "connectionName") {
            // return list of connector refs
            return Observable.create(observer => {
                let connectionRefs = [];
                WiContributionUtils.getConnections(this.http, "fabclient").subscribe((data: IConnectorContribution[]) => {
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
                        this.setMetadata(data)
                        if (this.metadata) {
                            let txn = [""];
                            let con = this.metadata["contract"]
                            Object.keys(con["transactions"]).forEach( (t) => {
                                console.log(t);
                                txn.push(t);
                            });
                            observer.next(txn);
                        }
                    });
                });
            }
        } else if (fieldName === "requestType" ) {
            let txnName = context.getField("transactionName").value;
            if (txnName && this.metadata) {
                let con = this.metadata["contract"];
                let txn = con["transactions"][txnName];
                console.log("txn: " + txn + " operation: " + txn["operation"]);
                return txn["operation"];
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
            return vresult;
        }
        return null;
    }

    setMetadata = (connector: any) => {
        if (this.metadata) {
            console.log("metadata already set");
            return;
        }
        for (let setting of connector.settings) {
            if (setting.name === "contract" && setting.value) {
                let content = this.extractFileContent(setting.value);
                console.log(content);
                this.metadata = JSON.parse(content);
            }
        }
    }

    extractFileContent = (selector: object): string => {
        let content = selector["content"]
        let data = content.substring(content.indexOf("base64,")+7)
//        console.log(data)
        return atob(data);
    }
}
