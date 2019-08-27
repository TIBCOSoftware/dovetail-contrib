import {Injectable, Injector} from "@angular/core";
import {Http} from "@angular/http";
import {Observable} from "rxjs/Observable";
import {
    WiContrib,
    WiServiceHandlerContribution,
    IValidationResult,
    ValidationResult,
    ActionResult,
    IActionResult,
    ICreateFlowActionContext,
    CreateFlowActionResult,
    WiContribModelService,
    WiContributionUtils,
    IConnectorContribution
} from "wi-studio/app/contrib/wi-contrib";
import { ITriggerContribution, IFieldDefinition, IConnectionAllowedValue, MODE } from "wi-studio/common/models/contrib";
import * as lodash from "lodash";

@WiContrib({})
@Injectable()
export class SchedulerTriggerHandler extends WiServiceHandlerContribution {
    
    constructor(private injector: Injector, private http: Http, private contribModelService: WiContribModelService) {
        super(injector, http, contribModelService);
    }
    
    value = (fieldName: string, context: ITriggerContribution): Observable<any> | any => {
        let conId = context.getField("model").value;
        
        switch(fieldName) {
            case "model":
                let connectionRefs = [];
                
                return Observable.create(observer => {
                    WiContributionUtils.getConnections(this.http, "SmartContract").subscribe((data: IConnectorContribution[]) => {
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
                
            case "asset":
                if(Boolean(conId) == false)
                    return null;

                return Observable.create(observer => {
                    WiContributionUtils.getConnection(this.http, conId)
                                        .map(data => data)
                                        .subscribe(data => {
                                            for (let setting of data.settings) {
                                                if (setting.name === "assets") {
                                                    observer.next(setting.value);
                                                    break;
                                                }
                                            }
                                        });
                });
            case "transactionInput":
                let asset = context.getField("asset").value;
                if(Boolean(conId) == false || Boolean(asset) == false)
                    return null;

                return Observable.create(observer => {
                    this.getSchemas(conId).subscribe( schemas => {
                        observer.next(schemas[asset]);
                       /* let schema = JSON.parse(schemas[asset]);
                        
                        let newSchema = {};
                        newSchema["$schema"] = schema["$schema"];
                        newSchema["title"] = schema["title"];
                        newSchema["type"] = "object"
                        let properties = {};
                        properties["data"] = {type: "object", properties: {}};
                        properties["data"].properties = schema.properties;
                        properties["ref"] = {type: "string"}
                        newSchema["properties"] = properties;
                        newSchema["description"] = schema.description;
                        observer.next(JSON.stringify(newSchema));*/
                    });
                });
            default: 
                return null;
        }
            
    }

    validate = (fieldName: string, context: ITriggerContribution): Observable<IValidationResult> | IValidationResult => {
       return null;
    }

    action = (actionId: string, context: ICreateFlowActionContext): Observable<IActionResult> | IActionResult => {
       
        let result = CreateFlowActionResult.newActionResult();
        let conId = context.getField("model").value;

        let asset = context.getField("asset").value
         
        return Observable.create(observer => {
            this.createFlow(context, conId, asset, result);   
            let actionResult = ActionResult.newActionResult().setSuccess(true).setResult(result);
            observer.next(actionResult);
        });
    }

    createFlow(context, conId, asset, result) : string{
        let modelService = this.getModelService();
        let trigger = modelService.createTriggerElement("SmartContract-Corda/SchedulerTrigger");
        if (trigger) {
            for(let t = 0; t < trigger.settings.length; t++) {
                if (trigger.settings[t].name === "model" ) {
                    trigger.settings[t].value = conId;
                    break;
                }
            }
            for (let s = 0; s < trigger.handler.settings.length; s++) {
                if (trigger.handler.settings[s].name === "asset") {
                    trigger.handler.settings[s].value = asset;
                    break;
                } 
            }
        }

        let flowName = context.getFlowName();
        let flowModel = modelService.createFlow(flowName, context.getFlowDescription());
        let scheduledactivity = modelService.createFlowElement("SmartContract-Corda/scheduledtask");
        flowModel.addFlowElement(scheduledactivity);
        result = result.addTriggerFlowMapping(lodash.cloneDeep(trigger), lodash.cloneDeep(flowModel));
        return flowName;
    }

    getSchemas(conId):  Observable<any> {
        let schemas = new Map();
        return Observable.create(observer => {
            WiContributionUtils.getConnection(this.http, conId)
                            .map(data => data)
                            .subscribe(data => {
                                let schemas = new Map();
                                for (let setting of data.settings) {
                                    if(setting.name === "schemas") {
                                        setting.value.map(item => schemas[item[0]] = item[1]);
                                        observer.next(schemas);
                                        break;
                                    }
                                }
                            });
                        });
    }
}