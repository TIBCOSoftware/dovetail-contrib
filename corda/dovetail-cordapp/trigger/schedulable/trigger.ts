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
const zstring = require("./lz-string")

@WiContrib({})
@Injectable()
export class R3FlowReceiverTriggerHandler extends WiServiceHandlerContribution {
    
    constructor(private injector: Injector, private http: Http, private contribModelService: WiContribModelService) {
        super(injector, http, contribModelService);
    }
    
    value = (fieldName: string, context: ITriggerContribution): Observable<any> | any => {
        let conId = context.getField("contract").value;
        switch(fieldName) {
            case "transactionInput":
                var event = context.getField("event").value;
                if(Boolean(conId) == false || Boolean(event) == false)
                    return null;

                return Observable.create(observer => {
                    this.getSchemas(conId).subscribe( schemas => {
                        var txnschema = schemas.get(event)
                        
                        observer.next(txnschema);
                        
                    });
                });  
            case "contract":
                let connectionRefs = [];
                return Observable.create(observer => {
                    WiContributionUtils.getConnections(this.http, "Dovetail-Contract").subscribe((data: IConnectorContribution[]) => {
                        data.forEach(connection => {
                            if ((<any>connection).isValid) {
                                for(let i=0; i < connection.settings.length; i++) {
                                    if(connection.settings[i].name === "display"){
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
            case "event":
                if(Boolean(conId) == false)
                    return null;

                return Observable.create(observer => {
                    WiContributionUtils.getConnection(this.http, conId)
                                        .map(data => data)
                                        .subscribe(data => {
                                            for (let setting of data.settings) {
                                                if (setting.name === "events") {
                                                    observer.next(setting.value);
                                                    break;
                                                }
                                            }
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
        let flows = []
        return Observable.create(observer => {
            let flowName = context.getFlowName();
            this.createFlow(context, flowName, result);                                     
            let actionResult = ActionResult.newActionResult().setSuccess(true).setResult(result);
            observer.next(actionResult);
        });
    }

    createFlow(context, flowName, result) : string{
        let modelService = this.getModelService();
        let trigger = modelService.createTriggerElement("Dovetail-CorDApp/R3SchedulableFlow");
        if (trigger) {
            for (let s = 0; s < trigger.handler.settings.length; s++) {
                if (trigger.handler.settings[s].name === "contract") {
                    trigger.handler.settings[s].value = context.getField("contract").value;
                } else  if (trigger.handler.settings[s].name === "event") {
                    trigger.handler.settings[s].value = context.getField("event").value;
                } 
            }

        }
        let flowModel = modelService.createFlow(flowName, context.getFlowDescription());
        let subflow = modelService.createFlowElement("Default/flogo-subflow");
        flowModel.addFlowElement(subflow);
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
                                        setting.value.map(item => {
                                            var s = item[1]
                                            var d = zstring.decompressFromUTF16(s) 
                                            schemas.set(item[0],d)
                                        });
                                        observer.next(schemas);
                                        break;
                                    }
                                }
                            });
                        });
    }
}