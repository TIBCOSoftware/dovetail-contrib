/// <amd-dependency path="./common"/>
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

const commonjs = require("./common");

@WiContrib({})
@Injectable()
export class SubflowActivityContributionHandler extends WiServiceHandlerContribution {
    constructor(private injector: Injector, private http: Http,) {
        super(injector, http);
    }
   
    value = (fieldName: string, context: IActivityContribution): any | Observable<any> => {
        switch(fieldName) {
            case "input":
                let schemaSelection = context.getField("schemaSelection").value;
                if (schemaSelection == "user"){
                    if(Boolean(context.getField("inputParams").value))
                        return commonjs.createFlowInputSchema(context.getField("inputParams").value.value)  
                } else {
                    return Observable.create(observer => {
                        this.getSchemas(schemaSelection).subscribe( schema => {
                            observer.next(schema);
                        });
                    }); 
                }
                break;
            case "schemaSelection":
                let connectionRefs = [];
                connectionRefs.push({
                    "unique_id": "user",
                    "name": "User Defined..."
                });
                return Observable.create(observer => {
                    WiContributionUtils.getConnections(this.http, "CorDApp").subscribe((data: IConnectorContribution[]) => {
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
        }
        return null;
    }
 
    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        let schemaSelection = context.getField("schemaSelection").value;
        switch (fieldName) {
            case "inputParams":
                return Observable.create(observer => {
                    let vresult: IValidationResult = ValidationResult.newValidationResult();
                    vresult.setVisible(schemaSelection == "user");
                    observer.next(vresult);
                });
        }
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
                                    if(setting.name === "inputParams") {
                                        observer.next(commonjs.createFlowInputSchema(setting.value.value));
                                        break;
                                    }
                                }
                            });
                        });
    }
}