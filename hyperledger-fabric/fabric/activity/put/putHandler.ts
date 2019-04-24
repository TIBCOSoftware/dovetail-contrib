
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
    IActivityContribution    
} from "wi-studio/app/contrib/wi-contrib";
import { IValidationResult, ValidationResult } from "wi-studio/common/models/validation";

@WiContrib({})
@Injectable()
export class putHandler extends WiServiceHandlerContribution {
    constructor( @Inject(Injector) injector) {
        super(injector);
    }

    value = (fieldName: string, context: IActivityContribution): Observable<any> | any => {
        if (fieldName === "result") {
            let valueTypeField: IFieldDefinition = context.getField("valueType");
            if (valueTypeField.value && valueTypeField.value === "object") {
                let dataField: IFieldDefinition = context.getField("data");
                if (dataField && dataField.value && dataField.value.value) {
                    return dataField.value.value;
                }
            }
        }
        return null;
    }

    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        if (fieldName === "value") {
            let vresult: IValidationResult = ValidationResult.newValidationResult();
            let valueTypeField: IFieldDefinition = context.getField("valueType");
            let valueField: IFieldDefinition = context.getField("value");
            if (valueTypeField.value && valueTypeField.value === "object") {
                vresult.setVisible(false);
            } else {
                vresult.setVisible(true);
            }
            return vresult;
        } else if (fieldName === "collection") {
            let vresult: IValidationResult = ValidationResult.newValidationResult();
            let isPrivateField: IFieldDefinition = context.getField("isPrivate");
            let collectionField: IFieldDefinition = context.getField("collection");
            if (isPrivateField.value && isPrivateField.value === true) {
                if (collectionField.display && collectionField.display.visible == false) {
                    vresult.setVisible(true);
                }
            } else {
                vresult.setVisible(false);
            }
            return vresult;
        } else if (fieldName === "data") {
            let vresult: IValidationResult = ValidationResult.newValidationResult();
            let valueTypeField: IFieldDefinition = context.getField("valueType");
            let dataField: IFieldDefinition = context.getField("data");
            if (valueTypeField.value && valueTypeField.value === "object") {
                if (dataField.display && dataField.display.visible == false) {
                    vresult.setVisible(true);
                }
                if (dataField.value === null || dataField.value.value === null || dataField.value.value === "") {
                    vresult.setError("FABTIC-PUT-1010", "Data definition must not be empty");
                } else {
                    try {
                        let valRes;
                        valRes = JSON.parse(dataField.value.value);
                        valRes = JSON.stringify(valRes);
                    } catch (e) {
                        vresult.setError("FABTIC-PUT-1020", "Invalid JSON: " + e.toString());
                    }
                }
            } else {
                vresult.setVisible(false);
            }
            return vresult;
        } else if (fieldName === "result") {
            let vresult: IValidationResult = ValidationResult.newValidationResult();
            let valueTypeField: IFieldDefinition = context.getField("valueType");
            if (valueTypeField.value && valueTypeField.value === "object") {
                vresult.setVisible(true);
            } else {
                vresult.setVisible(false);
            }
            return vresult;
        }
        return null;
    }
}
