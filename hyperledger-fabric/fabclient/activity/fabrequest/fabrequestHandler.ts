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

    constructor(private injector: Injector, private http: Http) {
        super(injector, http);
    }

    value = (fieldName: string, context: IContributionTypes): Observable<any> | any => {
        if (fieldName === "connectionName") {
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
        } else {
            return null;
        }
    }

    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        if (fieldName === "parameters" || fieldName === "transient" || fieldName === "result") {
            return Observable.create(observer => {
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
                observer.next(vresult);
            });
        } else {
            return null;
        }
    }

    action = (actionId: string, context: IContributionTypes): Observable<IActionResult> | IActionResult => {
        return Observable.create(observer => {
            let aresult: IActionResult = ActionResult.newActionResult();
            observer.next(aresult);
        });
    }
}
