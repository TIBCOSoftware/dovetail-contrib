"use strict";
var __extends = (this && this.__extends) || (function () {
    var extendStatics = Object.setPrototypeOf ||
        ({ __proto__: [] } instanceof Array && function (d, b) { d.__proto__ = b; }) ||
        function (d, b) { for (var p in b) if (b.hasOwnProperty(p)) d[p] = b[p]; };
    return function (d, b) {
        extendStatics(d, b);
        function __() { this.constructor = d; }
        d.prototype = b === null ? Object.create(b) : (__.prototype = b.prototype, new __());
    };
})();
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
var __param = (this && this.__param) || function (paramIndex, decorator) {
    return function (target, key) { decorator(target, key, paramIndex); }
};
Object.defineProperty(exports, "__esModule", { value: true });
var http_1 = require("@angular/http");
var core_1 = require("@angular/core");
var wi_contrib_1 = require("wi-studio/app/contrib/wi-contrib");
var SubflowActivityContribution = (function (_super) {
    __extends(SubflowActivityContribution, _super);
    function SubflowActivityContribution(injector, http, contribModelService) {
        var _this = _super.call(this, injector, http, contribModelService) || this;
        _this.http = http;
        _this.contribModelService = contribModelService;
        _this.value = function (fieldName, context) {
            var modelService = _this.getModelService();
            var applicationModel = modelService.getApplication();
            var flow = context.getField("flowURI").value;
            if (fieldName === "flowURI") {
                var list_1 = [];
                var flows_1 = [];
                if (applicationModel) {
                    var triggerMappings = applicationModel.getTriggerFlowModelMaps();
                    triggerMappings.map(function (triggerMapping) {
                        if ((context.getCurrentFlowName() !== triggerMapping.getFlowModel().getName()) && !triggerMapping.getFlowModel().isTriggerFlow()) {
                            if (list_1.indexOf(triggerMapping.getFlowModel().getName()) === -1) {
                                list_1.push(triggerMapping.getFlowModel().getName());
                                flows_1.push({
                                    "unique_id": "res://flow:" + triggerMapping.getFlowModel().getName().replace(/ /g, "_"),
                                    "name": triggerMapping.getFlowModel().getName()
                                });
                            }
                        }
                    });
                }
                return flows_1;
            }
            else if (fieldName === "input") {
                if (applicationModel && flow) {
                    var triggerMappings = applicationModel.getTriggerFlowModelMaps(), schema = void 0;
                    for (var i = 0; i < triggerMappings.length; i++) {
                        var triggerMapping = triggerMappings[i];
                        if (flow === "res://flow:" + triggerMapping.getFlowModel().getName().replace(/ /g, "_")) {
                            schema = triggerMapping.getFlowModel().getFlowInputSchema().json;
                            break;
                        }
                    }
                    return schema;
                }
            }
            else if (fieldName === "output") {
                if (applicationModel && flow) {
                    var triggerMappings = applicationModel.getTriggerFlowModelMaps(), schema = void 0;
                    for (var i = 0; i < triggerMappings.length; i++) {
                        var triggerMapping = triggerMappings[i];
                        if (flow === "res://flow:" + triggerMapping.getFlowModel().getName().replace(/ /g, "_")) {
                            schema = triggerMapping.getFlowModel().getFlowOutputSchema().json;
                            break;
                        }
                    }
                    ;
                    return schema;
                }
            }
            return null;
        };
        _this.validate = function (fieldName, context) {
            var modelService = _this.getModelService();
            var applicationModel = modelService.getApplication();
            var flow = context.getField("flowURI").value;
            if (fieldName === "flowURI" && flow) {
                var list_2 = [];
                var flows_2 = {};
                if (applicationModel) {
                    var triggerMappings = applicationModel.getTriggerFlowModelMaps();
                    triggerMappings.map(function (triggerMapping) {
                        if (!triggerMapping.getFlowModel().isTriggerFlow()) {
                            var name_1 = "res://flow:" + triggerMapping.getFlowModel().getName().replace(/ /g, "_");
                            if (list_2.indexOf(name_1) === -1) {
                                list_2.push(name_1);
                                flows_2[name_1] = triggerMapping.getFlowModel();
                            }
                        }
                    });
                    if (_this.isLoopDetected(flow, flow, flows_2)) {
                        return wi_contrib_1.ValidationResult.newValidationResult().setError("SUBFLOW-LOOP", "Cyclic dependency detected in the subflow");
                    }
                    else {
                        return wi_contrib_1.ValidationResult.newValidationResult();
                    }
                }
            }
            return null;
        };
        return _this;
    }
    SubflowActivityContribution.prototype.isLoopDetected = function (destination, source, flows) {
        if (flows && flows[source]) {
            var errorFlow = flows[source].getErrorFlow();
            if (this.checkLoopInFlow(destination, source, flows[source], flows)) {
                return true;
            }
            else if (errorFlow && this.checkLoopInFlow(destination, source, errorFlow, flows)) {
                return true;
            }
        }
        return false;
    };
    SubflowActivityContribution.prototype.checkLoopInFlow = function (destination, source, flow, flows) {
        var element = flow.getFlowElement();
        if (this.checkCurrentElement(element, destination, flows) || this.checkChildren(element, destination, flows)) {
            return true;
        }
        return false;
    };
    SubflowActivityContribution.prototype.checkChildren = function (element, destination, flows) {
        var children = element.getChildren();
        var isError = false;
        for (var _i = 0, children_1 = children; _i < children_1.length; _i++) {
            var child = children_1[_i];
            if (this.checkCurrentElement(child, destination, flows)) {
                isError = true;
                break;
            }
            else {
                isError = this.checkChildren(child, destination, flows);
            }
        }
        return isError;
    };
    SubflowActivityContribution.prototype.checkCurrentElement = function (element, destination, flows) {
        if (element.name === "flogo-subflow") {
            var flow = element.getField("flowURI").value;
            if (flow === destination) {
                return true;
            }
            else {
                return this.isLoopDetected(destination, flow, flows);
            }
        }
        return false;
    };
    return SubflowActivityContribution;
}(wi_contrib_1.WiServiceHandlerContribution));
SubflowActivityContribution = __decorate([
    wi_contrib_1.WiContrib({}),
    core_1.Injectable(),
    __param(0, core_1.Inject(core_1.Injector)),
    __metadata("design:paramtypes", [Object, http_1.Http, wi_contrib_1.WiContribModelService])
], SubflowActivityContribution);
exports.SubflowActivityContribution = SubflowActivityContribution;
//# sourceMappingURL=activity.js.map