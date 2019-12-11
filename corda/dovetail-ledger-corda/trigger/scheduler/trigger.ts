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
        let asset = context.getField("asset").value;
        switch(fieldName) {
            case "asset":
                let connectionRefs = [];
                
                return Observable.create(observer => {
                    WiContributionUtils.getConnections(this.http, "Dovetail-Ledger").subscribe((data: IConnectorContribution[]) => {
                        data.forEach(connection => {
                            if ((<any>connection).isValid) {
                                for(let i=0; i < connection.settings.length; i++) {
                                    if(connection.settings[i].name === "displayname"){
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
            case "assetName":
                if(Boolean(asset) == false)
                    return null;

                return Observable.create(observer => {
                    this.getSchema(asset).subscribe( result => {
                        observer.next(result.get("module") + "." + result.get("name"))
                    })
                });
            case "transactionInput":
                if(Boolean(asset) == false)
                    return null;

                return Observable.create(observer => {
                    this.getSchema(asset).subscribe( result => {
                        var json = JSON.parse(result.get("schema"))
                        var metadata = JSON.parse(json.description)
                        metadata.metadata.type = "ScheduledEvent"
                        metadata.metadata.asset = result.get("module") + "." + result.get("name")
                        metadata.metadata["actors"] = []
                        var authorized = context.getField("actors").value.value
                        if(authorized){
                            var authjson = JSON.parse(authorized)
                            for(var i=0; i<authjson.length; i++){
                                metadata.metadata.actors.push(authjson[i].party + "|" + authjson[i].certAttributes)
                            }
                        }
                        json.description = JSON.stringify(metadata)
                        observer.next(JSON.stringify(json));
                    });
                });
            case "data":
                   var schema = {
                        "$schema": "http://json-schema.org/draft-07/schema#", 
                        "type": "object", 
                        "properties": {"scheduledAt":{"type":"string", "format":"date-time"}} ,
                        "required":["scheduledAt"]
                    }
                    return JSON.stringify(schema)
            case "actors":
                    var actors = context.getField("actors").value
                    if(Boolean(asset) == false || Boolean(actors) == false || Boolean(actors.value))
                        return null;
    
                    return Observable.create(observer => {
                        WiContributionUtils.getConnection(this.http, asset)
                                            .map(data => data)
                                            .subscribe(data => {
                                                var party = []
                                                var json1 = JSON.parse(this.getSettingValue(data, "schema"))
                                                var json2 = JSON.parse(json1.description)
                                                var available = Array.from(new Set(json2.metadata.issueSigners.concat(json2.metadata.exitSigners).concat(json2.metadata.participants)))
                                                
                                                for (var p of available) {
                                                    party.push({party:p, certAttributes:""})
                                                }
                                                
                                                if(party.length > 0)
                                                    observer.next(JSON.stringify(party));
                                                else
                                                    observer.next(null)
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
        let asset = context.getField("asset").value
         
        return Observable.create(observer => {
            this.createFlow(context, asset, result);   
            let actionResult = ActionResult.newActionResult().setSuccess(true).setResult(result);
            observer.next(actionResult);
        });
    }

    createFlow(context, conId, result) : string{
        let modelService = this.getModelService();
        let trigger = modelService.createTriggerElement("Dovetail-Ledger-Corda/SchedulerTrigger");
        if (trigger) {
            for (let s = 0; s < trigger.handler.settings.length; s++) {
                if (trigger.handler.settings[s].name === "asset") {
                    trigger.handler.settings[s].value = conId;
                    break;
                } 
            }
        }

        let flowName = context.getFlowName();
        let flowModel = modelService.createFlow(flowName, context.getFlowDescription());
        let scheduledactivity = modelService.createFlowElement("Dovetail-Ledger-Corda/scheduledtask");
        flowModel.addFlowElement(scheduledactivity);
        result = result.addTriggerFlowMapping(lodash.cloneDeep(trigger), lodash.cloneDeep(flowModel));
        return flowName;
    }

    getSchema(conId):  Observable<any> {
        return Observable.create(observer => {
            WiContributionUtils.getConnection(this.http, conId)
                            .map(data => data)
                            .subscribe(data => {
                                var result = new Map()
                                for (let setting of data.settings) {
                                    if(setting.name === "schema") {
                                        result.set("schema", setting.value)
                                    } else if(setting.name === "name") {
                                        result.set("name", setting.value)
                                    } else if(setting.name === "module") {
                                        result.set("module", setting.value)
                                    }
                                }
                                observer.next(result);
                            });
                        });
    }

    getSettingValue(connection, setting):string {
        for(let i=0; i < connection.settings.length; i++) {
            if(connection.settings[i].name === setting){
                return connection.settings[i].value
            }
        }
    }
}