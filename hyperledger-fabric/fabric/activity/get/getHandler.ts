
import {Injectable, Injector, Inject} from "@angular/core";
import {Http} from "@angular/http";
import {Observable} from "rxjs/Observable";
import {
    WiContrib,
    WiServiceHandlerContribution,
    IContributionTypes,
    ActionResult,
    IActionResult,
    WiContribModelService,
    IFieldDefinition,
    WiContributionUtils,
    IActivityContribution,
    IConnectorContribution
} from "wi-studio/app/contrib/wi-contrib";
import { IValidationResult, ValidationResult } from "wi-studio/common/models/validation";

@WiContrib({})
@Injectable()
export class getHandler extends WiServiceHandlerContribution {
    constructor( private injector: Injector, private http: Http) {
        super(injector, http);
    }

    value = (fieldName: string, context: IActivityContribution): Observable<any> | any => {
        if (fieldName === "commonData") {
            return Observable.create(observer => {
                WiContributionUtils.getConnections(this.http, "fabric").subscribe((data: IConnectorContribution[]) => {
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
        } else if (fieldName === "dataType" ) {
            let connectorId = context.getField("commonData").value;
            if (connectorId) {
                return Observable.create(observer => {
                    WiContributionUtils.getConnection(this.http, connectorId).map(data => data)
                    .subscribe(data => {
                        for (let setting of data.settings) {
                            if (setting.name === "schema" && setting.value.value) {
//                                console.log("connector data: " + setting.name + " = " + setting.value.value);
                                let schema = [""];
                                let sobj = JSON.parse(setting.value.value);
                                Object.keys(sobj).forEach( (pkg) => {
                                    Object.keys(sobj[pkg]).forEach( (n) => {
                                        schema.push(pkg + "." + n);
                                    });
                                });
                                observer.next(schema);
                                break;
                            }
                        }
                    });
                });
            }
        } else if (fieldName === "result") {
            let dataTypeName = context.getField("dataType").value;
            let connectorId = context.getField("commonData").value;
            if (dataTypeName && connectorId) {
                // set pre-defined schema from shared data defs
                return Observable.create(observer => {
                    WiContributionUtils.getConnection(this.http, connectorId).map(data => data)
                    .subscribe(data => {
                        for (let setting of data.settings) {
                            if (setting.name === "schema" && setting.value && setting.value.value) {
//                                console.log("schema data: " + setting.name + " = " + setting.value.value);
                                let sch = JSON.parse(setting.value.value);
                                let idx = dataTypeName.lastIndexOf('.');
                                let pkg = sch[dataTypeName.substring(0, idx)];
                                observer.next(JSON.stringify(pkg[dataTypeName.substring(idx+1)], null, 2));
                                break;
                            }
                        }
                    });
                });
            }
        }
        return null;
    }

    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        if (fieldName === "dataType") {
            let vresult: IValidationResult = ValidationResult.newValidationResult();
            let commonDataField: IFieldDefinition = context.getField("commonData");
            if (commonDataField.value) {
                vresult.setVisible(true);
            } else {
                vresult.setVisible(false);
            }
            return vresult;
        } else if (fieldName === "result") {
            let vresult: IValidationResult = ValidationResult.newValidationResult();
            let dataTypeField: IFieldDefinition = context.getField("dataType");
            if (dataTypeField.value && dataTypeField.display.visible) {
//                console.log("data type is specified, set compositeKey readonly");
                vresult.setReadOnly(true);
            } else {
                vresult.setReadOnly(false);
            }
            let resultField: IFieldDefinition = context.getField("result");
            if (resultField.value) {
                try {
                    let valRes;
                    valRes = JSON.parse(resultField.value);
                    valRes = JSON.stringify(valRes);
                } catch (e) {
                    vresult.setError("FABRIC-GET-1000", "Invalid JSON: " + e.toString());
                }
            }
            return vresult;
        }
        return null;
    }
}
