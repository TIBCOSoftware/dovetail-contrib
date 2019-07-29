
import {Injectable, Inject, Injector} from "@angular/core";
import {Http} from "@angular/http";
import {Observable} from "rxjs/Observable";
import {
    WiContrib,
    WiServiceHandlerContribution,
    IActionResult,
    ActionResult,
    WiContribModelService,
    ICreateFlowActionContext,
    CreateFlowActionResult,
    WiContributionUtils,
    IActivityContribution,
    IConnectorContribution
} from "wi-studio/app/contrib/wi-contrib";
import { ITriggerContribution, IFieldDefinition, MODE } from "wi-studio/common/models/contrib";
import { IValidationResult, ValidationResult } from "wi-studio/common/models/validation";
import * as lodash from "lodash";

const ccEvent = {
    "type": "object",
    "properties": {
        "block": {
            "type": "integer"
        },
        "source": {
            "type": "string"
        },
        "txId": {
            "type": "string"
        },
        "chaincode": {
            "type": "string"
        },
        "name": {
            "type": "string"
        },
        "payload": {
            "type": "string"
        }
    }
}

const blockEvent = {
    "type": "object",
    "properties": {
        "block": {
            "type": "integer"
        },
        "source": {
            "type": "string"
        },
        "transactions": {
            "type": "array",
            "items": {
                "type": "object",
                "properties": {
                    "type": {
                        "type": "string"
                    },
                    "txId": {
                        "type": "string"
                    },
                    "txTime": {
                        "type": "string"
                    },
                    "channel": {
                        "type": "string"
                    },
                    "creator": {
                        "type": "object",
                        "properties": {
                            "wspid": {
                                "type": "string"
                            },
                            "subject": {
                                "type": "string"
                            },
                            "issuer": {
                                "type": "string"
                            },
                            "cert": {
                                "type": "string"
                            }        
                        }
                    },
                    "actions": {
                        "type": "array",
                        "items": {
                            "type": "object",
                            "properties": {
                                "chaincode": {
                                    "type": "object",
                                    "properties": {
                                        "type": {
                                            "type": "string"
                                        },
                                        "name": {
                                            "type": "string"
                                        },
                                        "version": {
                                            "type": "string"
                                        }            
                                    }
                                },
                                "input": {
                                    "type": "object",
                                    "properties": {
                                        "function": {
                                            "type": "string"
                                        },
                                        "args": {
                                            "type": "array",
                                            "items": {
                                                "type": "string"
                                            }
                                        },
                                        "transient": {
                                            "type": "string"
                                        }
                                    }
                                },
                                "result": {
                                    "type": "object",
                                    "properties": {
                                        "rwset": {
                                            "type": "integer"
                                        },
                                        "response": {
                                            "type": "object",
                                            "properties": {
                                                "status": {
                                                    "type": "integer"
                                                },
                                                "message": {
                                                    "type": "string"
                                                },
                                                "payload": {
                                                    "type": "string"
                                                }        
                                            }
                                        },
                                        "event": {
                                            "type": "object",
                                            "properties": {
                                                "name": {
                                                    "type": "string"
                                                },
                                                "payload": {
                                                    "type": "string"
                                                }        
                                            }
                                        }
                                    }
                                },
                                "endorsers": {
                                    "type": "integer"
                                }    
                            }
                        }
                    }
                }
            }
        }
    }
}

@WiContrib({})
@Injectable()
export class eventlistenerHandler extends WiServiceHandlerContribution {

    constructor(private injector: Injector, private http: Http, private contribModelService: WiContribModelService) {
        super(injector, http, contribModelService);
    }

    value = (fieldName: string, context: ITriggerContribution): Observable<any> | any => {
        if (fieldName === "connectionName") {
            return Observable.create(observer => {
                WiContributionUtils.getConnections(this.http, "fabclient").subscribe((data: IConnectorContribution[]) => {
                    let connectionRefs = [];
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
        } else if (fieldName === "data") {
            let eventTypeField: IFieldDefinition = context.getField("eventType");
            if (eventTypeField && eventTypeField.value == "Chaincode") {
                return JSON.stringify(ccEvent, null, 2);
            } else {
                return JSON.stringify(blockEvent, null, 2);
            }
        }
        return null;
    }

    // verify user entries are valid JSON string
    validate = (fieldName: string, context: ITriggerContribution): Observable<IValidationResult> | IValidationResult => {
        let vresult: IValidationResult = ValidationResult.newValidationResult();
        if (fieldName === "eventFilter" || fieldName === "chaincodeID") {
            let eventTypeField: IFieldDefinition = context.getField("eventType");
            if (eventTypeField.value == "Chaincode") {
                vresult.setVisible(true);
            } else {
                vresult.setVisible(false);
            }
            return vresult;
        }
        return null;
    }

    // used to configure trigger with data from "Add a trigger" wizard. it also give you option to copy schema to flow on creation
    action = (actionId: string, context: ICreateFlowActionContext): Observable<IActionResult> | IActionResult => {
        let modelService = this.getModelService();
        let result = CreateFlowActionResult.newActionResult();
        if (context.handler && context.handler.settings && context.handler.settings.length > 0) {
            let connectionField = <IFieldDefinition>context.getField("connectionName");
            let typeField = <IFieldDefinition>context.getField("eventType");

            if (typeField && typeField.value) {
                let trigger = modelService.createTriggerElement("fabclient/fabclient-eventlistener");
                if (trigger && trigger.handler && trigger.handler.settings && trigger.handler.settings.length > 0) {
                    for (let j = 0; j < trigger.handler.settings.length; j++) {
                        if (trigger.handler.settings[j].name === "connectionName") {
                            trigger.handler.settings[j].value = connectionField.value;
                        } else if (trigger.handler.settings[j].name === "eventType") {
                            trigger.handler.settings[j].value = typeField.value;
                        }
                    }
                }
                if (trigger && trigger.outputs && trigger.outputs.length > 0) {
                    for (let j = 0; j < trigger.outputs.length; j++) {
                        if (trigger.outputs[j].name === "data") {
                            let dataSchema = JSON.stringify(ccEvent, null, 2);
                            if (typeField.value != "Chaincode") {
                                dataSchema = JSON.stringify(blockEvent, null, 2)
                            }
                            trigger.outputs[j].value = {
                                "value": "",
                                "metadata": dataSchema
                            };
                            break;
                        }
                    }
                }
                let flowModel = modelService.createFlow(context.getFlowName(), context.getFlowDescription());
                result = result.addTriggerFlowMapping(lodash.cloneDeep(trigger), lodash.cloneDeep(flowModel));
            }
        }
        let actionResult = ActionResult.newActionResult().setSuccess(true).setResult(result);
        return actionResult;
    }
}
