import { HttpModule, Http } from "@angular/http";
import { NgModule } from "@angular/core";
import { WiServiceContribution } from "wi-studio/app/contrib/wi-contrib";
import { SubflowActivityContribution } from "./activity";


@NgModule({
    imports: [
        HttpModule,
    ],
    exports: [

    ],
    declarations: [

    ],
    entryComponents: [

    ],
    providers: [
        {
            provide: WiServiceContribution,
            useClass: SubflowActivityContribution
        }
    ]
})
export default class SubflowContribModule {

}
