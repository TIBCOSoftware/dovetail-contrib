
import {Injectable, Inject, Injector} from "@angular/core";
import {Http} from "@angular/http";
import {Observable} from "rxjs/Observable";
import {
    WiContrib,
    WiServiceHandlerContribution,
    ActionResult,
    IActionResult,
    IFieldDefinition,
    WiContribModelService,
    WiContributionUtils,
    IConnectorContribution,
    AUTHENTICATION_TYPE
} from "wi-studio/app/contrib/wi-contrib";
import { IValidationResult, ValidationResult, ValidationError } from "wi-studio/common/models/validation";

@WiContrib({})
@Injectable()
export class datadefHandler extends WiServiceHandlerContribution {

    constructor(private injector: Injector) {
        super(injector);
    }

    value = (fieldName: string, context: IConnectorContribution): Observable<any> | any => {
        return null;
    }

    validate = (fieldName: string, context: IConnectorContribution): Observable<IValidationResult> | IValidationResult => {
        if (fieldName === "schema" || fieldName === "data") {
            let vresult: IValidationResult = ValidationResult.newValidationResult();
            vresult.setReadOnly(true);
            return vresult;
        } else if (fieldName === "schemaList" || fieldName === "dataList") {
            let vresult: IValidationResult = ValidationResult.newValidationResult();
            let field: IFieldDefinition = context.getField(fieldName);
            let rows: any[] = [];
            try {
                rows = JSON.parse(field.value);
            } catch (e) {}
            if (rows) {
                let rowids: string[] = [];
                for (let row of rows) {
                    let myid = row["package"] + "." + row["name"];
                    let attr: string = "definition";
                    if (fieldName === "dataList") {
                        attr = "value";
                    }
                    if (!row[attr]) {
                        vresult.setError("FABRIC-DATA-1010", "Data definition cannot be empty for: " + myid);
                        vresult.setValid(false);
                        return vresult;
                    } else {
                        try {
                            let val = JSON.parse(row[attr]);
                            JSON.stringify(val);
                        } catch (e) {
                            vresult.setError("FABTIC-DATA-1020", "Invalid JSON for " + myid + ": " + e.toString());
                            vresult.setValid(false);
                            return vresult;    
                        }
                    }
                    for (let rid of rowids) {
                        if (myid === rid) {
                            vresult.setError("FABRIC-DATA-1000", "Package and name already exists: " + myid);
                            vresult.setValid(false);
                            return vresult;
                        }
                    }
                    rowids.push(myid);
                }
            }
        }
        return null;
    }

    action = (actionId: string, context: IConnectorContribution): Observable<IActionResult> | IActionResult => {
        let aresult: IActionResult = ActionResult.newActionResult();
        try {
            let schema = {};
            let schemaListDef: IFieldDefinition = context.getField("schemaList");
            if (schemaListDef.value) {
                let schemaList = JSON.parse(schemaListDef.value);
                console.log("schemaDef type: " + schemaList.constructor['name']);
                console.log("schemaDef length: " + schemaList.length + " value: " + schemaList);
                schemaList.forEach( (item) => {
                    let pName = item["package"];
                    if (!pName) {
                        pName = "default";
                    }
                    let pObj = schema[pName];
                    if (!pObj) {
                        schema[pName] = {};
                        pObj = schema[pName];
                    }
                    pObj[item["name"]] = JSON.parse(item["definition"]);
                });
            }
            let schemaDef: IFieldDefinition = context.getField("schema");
            schemaDef.value.value = JSON.stringify(schema, null, 2);

            let data = {};
            let dataListDef: IFieldDefinition = context.getField("dataList");
            if (dataListDef.value) {
                let dataList = JSON.parse(dataListDef.value);
                dataList.forEach( (item) => {
                    let pName = item["package"];
                    if (!pName) {
                        pName = "default";
                    }
                    let pObj = data[pName];
                    if (!pObj) {
                        data[pName] = {};
                        pObj = data[pName];
                    }
                    pObj[item["name"]] = JSON.parse(item["value"]);
                });
            }
            let dataDef: IFieldDefinition = context.getField("data");
            dataDef.value.value = JSON.stringify(data, null, 2);
        } catch(err) {
            console.log(actionId + ": " + err.message);
            aresult.setSuccess(false).setResult(new ValidationError("DATADEF_1000", "Failed to save: " + err.message));
            return aresult;
        }
        let actionResult = {
            context: context,
            authType: AUTHENTICATION_TYPE.BASIC,
            authData: {}
        };
        aresult.setResult(actionResult);
        return aresult;
    }
}
