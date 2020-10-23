import { Injectable, Injector, Inject } from "@angular/core";
import { Http } from "@angular/http";
import { Observable } from "rxjs/Observable";
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
  IConnectorContribution,
} from "wi-studio/app/contrib/wi-contrib";
import {
  IValidationResult,
  ValidationResult,
} from "wi-studio/common/models/validation";

@WiContrib({})
@Injectable()
export class getbycompositekeyHandler extends WiServiceHandlerContribution {
  constructor(private injector: Injector, private http: Http) {
    super(injector, http);
  }

  value = (
    fieldName: string,
    context: IActivityContribution
  ): Observable<any> | any => {
    if (fieldName === "result") {
      let resultField: IFieldDefinition = context.getField("result");
      if (resultField.value) {
        let sch = JSON.parse(resultField.value);
        // check schema object type, and convert to key-value type
        if (sch["type"] == "object") {
          let data = {};
          data["type"] = "object";
          data["properties"] = sch["properties"];
          let keyval = {};
          keyval["key"] = { type: "string" };
          keyval["value"] = data;
          let item = {};
          item["type"] = "object";
          item["properties"] = keyval;
          let resultschema = {};
          resultschema["$schema"] = sch["$schema"];
          resultschema["type"] = "array";
          resultschema["items"] = item;
          return JSON.stringify(resultschema, null, 2);
        }
      }
    }
    return null;
  };

  validate = (
    fieldName: string,
    context: IActivityContribution
  ): Observable<IValidationResult> | IValidationResult => {
    if (
      fieldName === "pageSize" ||
      fieldName === "start" ||
      fieldName === "bookmark"
    ) {
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
        attributesParsed = JSON.parse(attributesField.value);
      } catch (e) {}

      for (let attr of attributesParsed) {
        if (!attr.parameterName) {
          errMessage = "Parameter Name should not be empty";
          vresult.setError("FABRIC-GETCOMPOSITE-1000", errMessage);
          vresult.setValid(false);
          break;
        } else {
          for (let paramName of arrParamNamesTmp) {
            if (paramName === attr.parameterName) {
              errMessage =
                "Attribute Name '" + attr.parameterName + "' already exists";
              vresult.setError("FABRIC-GETCOMPOSITE-1000", errMessage);
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
      if (resultField.value) {
        try {
          let valRes;
          valRes = JSON.parse(resultField.value);
          valRes = JSON.stringify(valRes);
        } catch (e) {
          vresult.setError(
            "FABRIC-GETCOMPOSITE-1010",
            "Invalid JSON: " + e.toString()
          );
        }
      }
      return vresult;
    }
    return null;
  };
}
