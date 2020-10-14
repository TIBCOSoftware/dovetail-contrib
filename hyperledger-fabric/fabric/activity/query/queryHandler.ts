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
export class queryHandler extends WiServiceHandlerContribution {
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
    if (fieldName === "queryParams") {
      let vresult = ValidationResult.newValidationResult();
      let queryParamsField: IFieldDefinition = context.getField(fieldName);
      let arrParamNamesTmp: any[] = [];
      let errMessage: string = "";
      let queryParamsParsed: any = {};

      try {
        queryParamsParsed = JSON.parse(queryParamsField.value);
      } catch (e) {}

      for (let queryParam of queryParamsParsed) {
        if (!queryParam.parameterName) {
          errMessage = "Parameter Name should not be empty";
          vresult.setError("FABTIC-QUERY-1010", errMessage);
          vresult.setValid(false);
          break;
        } else {
          for (let paramName of arrParamNamesTmp) {
            if (paramName === queryParam.parameterName) {
              errMessage =
                "Parameter Name '" +
                queryParam.parameterName +
                "' already exists";
              vresult.setError("FABRIC-QUERY-1010", errMessage);
              vresult.setValid(false);
              break;
            }
          }
          arrParamNamesTmp.push(queryParam.parameterName);
        }
      }
      return vresult;
    } else if (
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
            "FABRIC-QUERY-1000",
            "Invalid JSON: " + e.toString()
          );
        }
      }
      return vresult;
    }
    return null;
  };
}
