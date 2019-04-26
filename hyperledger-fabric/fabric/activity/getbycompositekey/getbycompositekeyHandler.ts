
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
export class getbycompositekeyHandler extends WiServiceHandlerContribution {
    constructor( @Inject(Injector) injector) {
        super(injector);
    }

    value = (fieldName: string, context: IActivityContribution): Observable<any> | any => {
        return null;
    }

    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        if (fieldName === "collection") {
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
        } else if (fieldName === "pageSize" || fieldName === "start" || fieldName === "bookmark") {
            let vresult: IValidationResult = ValidationResult.newValidationResult();
            let paginationField: IFieldDefinition = context.getField("usePagination");
            let valueField: IFieldDefinition = context.getField(fieldName);
            if (paginationField.value && paginationField.value === true) {
                if (valueField.display && valueField.display.visible == false) {
                    vresult.setVisible(true);
                }
            } else {
                vresult.setVisible(false);
            }
            return vresult;
        } else if (fieldName === "attributes") {
            let vresult: IValidationResult = ValidationResult.newValidationResult();
            let attributesField: IFieldDefinition = context.getField(fieldName);
            let arrParamNamesTmp: any[] = [];
            let errMessage: string = "";
            let attributesParsed: any = {};

            try {
                attributesParsed = JSON.parse(attributesField.value.value);
            } catch (e) { }

            for (let attr of attributesParsed) {
                if (!attr.parameterName) {
                    errMessage = "Parameter Name should not be empty";
                    vresult.setError("FABTIC-GETCOMPOSITE-1000", errMessage);
                    vresult.setValid(false);
                    break;
                } else {
                    for (let paramName of arrParamNamesTmp) {
                        if (paramName === attr.parameterName) {
                            errMessage = "Attribute Name \'" + attr.parameterName + "\' already exists";
                            vresult.setError("FABTIC-GETCOMPOSITE-1000", errMessage);
                            vresult.setValid(false);
                            break;
                        }
                    }
                    arrParamNamesTmp.push(attr.parameterName);
                }
            }
            return vresult;
        } else if (fieldName === "result") {
            let vresult: IValidationResult = ValidationResult.newValidationResult();
            let resultField: IFieldDefinition = context.getField("result");
            if (resultField.value && resultField.value.value) {
                try {
                    let valRes;
                    valRes = JSON.parse(resultField.value.value);
                    valRes = JSON.stringify(valRes);
                } catch (e) {
                    vresult.setError("FABTIC-GETCOMPOSITE-1010", "Invalid JSON: " + e.toString());
                }
            }
            return vresult;
        }
        return null;
    }
}