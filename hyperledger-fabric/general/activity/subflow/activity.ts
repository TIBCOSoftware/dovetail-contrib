/*
 * Copyright © 2018. TIBCO Software Inc.
 * This file is subject to the license terms contained
 * in the license file that is distributed with this file.
 */
import {Observable} from "rxjs/Observable";
import {Http} from "@angular/http";
import {Inject, Injectable, Injector} from "@angular/core";
import {
    IActivityContribution,
    WiContrib,
    IValidationResult,
    WiServiceHandlerContribution,
    WiContribModelService,
    IFlow,
    IFlowElement,
    ValidationResult
} from "wi-studio/app/contrib/wi-contrib";
import * as lodash from "lodash";

@WiContrib({})
@Injectable()
export class SubflowActivityContribution extends WiServiceHandlerContribution {
    constructor(@Inject(Injector) injector, private http: Http, private contribModelService: WiContribModelService) {
        super(injector, http, contribModelService);
    }

    value = (fieldName: string, context: IActivityContribution): Observable<any> | any => {
        let modelService = this.getModelService();
        let applicationModel = modelService.getApplication();
        let flow = context.getField("flowURI").value;
        if (fieldName === "flowURI") {
            let list: string[] = [];
            let flows: any = [];
            if (applicationModel) {
                let triggerMappings = applicationModel.getTriggerFlowModelMaps();
                triggerMappings.map(triggerMapping => {
                    if ((context.getCurrentFlowName() !== triggerMapping.getFlowModel().getName()) && !triggerMapping.getFlowModel().isTriggerFlow()) {
                        if (list.indexOf(triggerMapping.getFlowModel().getName()) === -1) {
                            list.push(triggerMapping.getFlowModel().getName());
                            flows.push({
                                "unique_id": "res://flow:" + triggerMapping.getFlowModel().getName().replace(/ /g, "_"),
                                "name": triggerMapping.getFlowModel().getName()
                            });
                        }
                    }
                });
            }
            return flows;
        } else if (fieldName === "input") {
            if (applicationModel && flow) {
                let triggerMappings = applicationModel.getTriggerFlowModelMaps(), schema;
                for (let i = 0; i < triggerMappings.length; i++ ) {
                    let triggerMapping = triggerMappings[i];
                    if (flow === "res://flow:" + triggerMapping.getFlowModel().getName().replace(/ /g, "_")) {
                        schema = triggerMapping.getFlowModel().getFlowInputSchema().json;
                        break;
                    }
                }
                return schema;
            }
        } else if (fieldName === "output") {
            if (applicationModel && flow) {
                let triggerMappings = applicationModel.getTriggerFlowModelMaps(), schema;
                for (let i = 0; i < triggerMappings.length; i++ ) {
                    let triggerMapping = triggerMappings[i];
                    if (flow === "res://flow:" + triggerMapping.getFlowModel().getName().replace(/ /g, "_")) {
                        schema = triggerMapping.getFlowModel().getFlowOutputSchema().json;
                        break;
                    }
                };
                return schema;
            }
        }
        return null;
    }

    validate = (fieldName: string, context: IActivityContribution): Observable<IValidationResult> | IValidationResult => {
        let modelService = this.getModelService();
        let applicationModel = modelService.getApplication();
        let flow = context.getField("flowURI").value;
        if (fieldName === "flowURI" && flow) {
            let list: string[] = [];
            let flows: any = {};
            if (applicationModel) {
                let triggerMappings = applicationModel.getTriggerFlowModelMaps();
                triggerMappings.map(triggerMapping => {
                    if (!triggerMapping.getFlowModel().isTriggerFlow()) {
                        const name = "res://flow:" + triggerMapping.getFlowModel().getName().replace(/ /g, "_");
                        if (list.indexOf(name) === -1) {
                            list.push(name);
                            flows[name] = triggerMapping.getFlowModel();
                        }
                    }
                });
                if (this.isLoopDetected(flow, flow, flows)) {
                    return ValidationResult.newValidationResult().setError("SUBFLOW-LOOP", "Cyclic dependency detected in the subflow");
                } else {
                    return ValidationResult.newValidationResult();
                }
            }
        }
        return null;
    }

    isLoopDetected(destination: string, source: string, flows: any): boolean {
        if (flows && flows[source]) {
            const errorFlow = (<IFlow>flows[source]).getErrorFlow();
            if (this.checkLoopInFlow(destination, source, flows[source], flows)) {
                return true;
            } else if (errorFlow && this.checkLoopInFlow(destination, source, errorFlow, flows)) {
                return true;
            }
        }
        return false;
    }

    checkLoopInFlow(destination: string, source: string, flow: any, flows: any): boolean {
        const element = flow.getFlowElement();
        if (this.checkCurrentElement(element, destination, flows) || this.checkChildren(element, destination, flows)) {
            return true;
        }
        return false;
    }

    checkChildren(element: IFlowElement, destination: string, flows: any): boolean {
        const children = element.getChildren();
        let isError = false;
        for (let child of children) {
            if (this.checkCurrentElement(child, destination, flows)) {
                isError = true;
                break;
            } else {
                isError = this.checkChildren(child, destination, flows);
            }
        }
        return isError;

    }

    checkCurrentElement(element: IFlowElement, destination: string, flows: any) {
        if (element.name === "flogo-subflow") {
            let flow = element.getField("flowURI").value;
            if (flow === destination) {
                return true;
            } else {
                return this.isLoopDetected(destination, flow, flows);
            }
        }
        return false;
    }
}
