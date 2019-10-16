"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
var http_1 = require("@angular/http");
var core_1 = require("@angular/core");
var wi_contrib_1 = require("wi-studio/app/contrib/wi-contrib");
var activity_1 = require("./activity");
var SubflowContribModule = (function () {
    function SubflowContribModule() {
    }
    return SubflowContribModule;
}());
SubflowContribModule = __decorate([
    core_1.NgModule({
        imports: [
            http_1.HttpModule,
        ],
        exports: [],
        declarations: [],
        entryComponents: [],
        providers: [
            {
                provide: wi_contrib_1.WiServiceContribution,
                useClass: activity_1.SubflowActivityContribution
            }
        ]
    })
], SubflowContribModule);
exports.default = SubflowContribModule;
//# sourceMappingURL=activity.module.js.map