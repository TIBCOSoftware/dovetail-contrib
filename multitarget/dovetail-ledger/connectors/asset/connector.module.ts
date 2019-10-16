import {NgModule} from "@angular/core"
import {HttpModule} from "@angular/http";
import {WiServiceContribution} from "wi-studio/app/contrib/wi-contrib"
import {AssetConnectorService} from "./connector"

@NgModule({
    imports: [
        HttpModule
    ],
    providers: [
        {
            provide: WiServiceContribution,
            useClass: AssetConnectorService
        }
    ]
})
export default class AssetConnectorModule {

}
